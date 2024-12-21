package ticket

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/domain/model/converter"
	"github.com/TrinityKnights/Backend/internal/repository/ticket"
	"github.com/TrinityKnights/Backend/pkg/cache"
	domainErrors "github.com/TrinityKnights/Backend/pkg/errors"
	"github.com/TrinityKnights/Backend/pkg/helper"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TicketServiceImpl struct {
	DB               *gorm.DB
	Cache            *cache.ImplCache
	Log              *logrus.Logger
	Validate         *validator.Validate
	TicketRepository ticket.TicketRepository
	helper           *helper.ContextHelper
}

func NewTicketServiceImpl(db *gorm.DB, cacheImpl *cache.ImplCache, log *logrus.Logger, validate *validator.Validate, ticketRepository ticket.TicketRepository) *TicketServiceImpl {
	return &TicketServiceImpl{
		DB:               db,
		Cache:            cacheImpl,
		Log:              log,
		Validate:         validate,
		TicketRepository: ticketRepository,
		helper:           helper.NewContextHelper(),
	}
}

const MaxTicketCount = 1000

func (s *TicketServiceImpl) CreateTicket(ctx context.Context, request *model.CreateTicketRequest) ([]*model.TicketResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	if request.Count < 0 || request.Count > MaxTicketCount {
		return nil, domainErrors.ErrBadRequest
	}

	ticketType := helper.TicketUpper(request.Type)
	if ticketType.Short == "" || ticketType.Long == "" {
		return nil, domainErrors.ErrBadRequest
	}

	// Get the last ticket number using the repository method
	lastTicket, err := s.TicketRepository.GetLastTicketNumber(tx, request.EventID, ticketType.Short)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Log.Errorf("failed to get last ticket: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	// Initialize starting seat number
	startingNumber := 1
	if lastTicket != nil {
		parts := strings.Split(lastTicket.SeatNumber, "-")
		if len(parts) == 2 {
			if num, err := strconv.Atoi(parts[1]); err == nil {
				startingNumber = num + 1
			}
		}
	}

	tickets := make([]*entity.Ticket, request.Count)
	for i := 0; i < request.Count; i++ {
		tickets[i] = &entity.Ticket{
			ID:         fmt.Sprintf("T-%s", uuid.NewString()[:6]),
			EventID:    request.EventID,
			Price:      request.Price,
			Type:       ticketType.Long,
			SeatNumber: fmt.Sprintf("%s-%d", ticketType.Short, startingNumber+i),
		}
	}

	if err := s.TicketRepository.CreateBatch(tx, tickets); err != nil {
		s.Log.Errorf("failed to create tickets: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	return converter.TicketsToResponses(tickets), nil
}

func (s *TicketServiceImpl) UpdateTicket(ctx context.Context, request *model.UpdateTicketRequest) (*model.TicketResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	if _, err := s.TicketRepository.Find(tx, &model.TicketQueryOptions{
		ID: &request.ID,
	}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrNotFound
		}
		return nil, domainErrors.ErrInternalServer
	}

	data := &entity.Ticket{
		ID:         request.ID,
		EventID:    request.EventID,
		OrderID:    request.OrderID,
		Price:      request.Price,
		Type:       request.Type,
		SeatNumber: request.SeatNumber,
	}

	if err := s.TicketRepository.Update(tx, data); err != nil {
		s.Log.Errorf("failed to update ticket: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		return nil, domainErrors.ErrInternalServer
	}

	return converter.TicketEntityToResponse(data), nil
}

func (s *TicketServiceImpl) GetTicketByID(ctx context.Context, request *model.GetTicketRequest) (*model.TicketResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	key := fmt.Sprintf("ticket:get:id:%s", request.ID)
	var cacheResponse *model.TicketResponse
	if err := s.Cache.Get(key, &cacheResponse); err == nil {
		return cacheResponse, nil
	}

	tx := s.DB.WithContext(ctx)

	t, err := s.TicketRepository.Find(tx.Preload("Order"), &model.TicketQueryOptions{
		ID: &request.ID,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrNotFound
		}
		s.Log.Errorf("failed to get ticket: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if len(t) == 0 {
		return nil, domainErrors.ErrNotFound
	}

	response := converter.TicketEntityToResponse(t[0])

	if err := s.Cache.Set(key, response, 5*time.Minute); err != nil {
		s.Log.Errorf("failed to set cache: %v", err)
	}

	return response, nil
}

func (s *TicketServiceImpl) GetTickets(ctx context.Context, request *model.TicketsRequest) (*model.Response[[]*model.TicketResponse], error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	opts := model.TicketQueryOptions{
		Page:  request.Page,
		Size:  request.Size,
		Sort:  request.Sort,
		Order: request.Order,
	}

	if opts.Size <= 0 {
		opts.Size = 10
	}
	if opts.Page <= 0 {
		opts.Page = 1
	}

	cacheKey := fmt.Sprintf("ticket:get:page:%d:size:%d:sort:%s:order:%s", opts.Page, opts.Size, opts.Sort, opts.Order)
	var cacheResponse model.Response[[]*model.TicketResponse]
	if err := s.Cache.Get(cacheKey, &cacheResponse); err == nil {
		return &cacheResponse, nil
	}

	db := s.DB.WithContext(ctx)

	tickets, err := s.TicketRepository.Find(db, &model.TicketQueryOptions{
		Page:  opts.Page,
		Size:  opts.Size,
		Sort:  opts.Sort,
		Order: opts.Order,
	})
	if err != nil {
		s.Log.Errorf("failed to get tickets: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrNotFound
		}
		return nil, domainErrors.ErrInternalServer
	}

	if len(tickets) == 0 {
		return nil, domainErrors.ErrNotFound
	}

	response := converter.TicketsToPaginatedResponse(tickets, int64(len(tickets)), opts.Page, opts.Size)

	if err := s.Cache.Set(cacheKey, response, 5*time.Minute); err != nil {
		s.Log.Errorf("failed to set cache: %v", err)
	}

	return response, nil
}

func (s *TicketServiceImpl) SearchTickets(ctx context.Context, request *model.TicketSearchRequest) (*model.Response[[]*model.TicketResponse], error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	// Set default values for pagination
	if request.Size <= 0 {
		request.Size = 10
	}
	if request.Page <= 0 {
		request.Page = 1
	}

	// Generate cache key based on search parameters
	ticketType := strings.ToUpper(request.Type)
	cacheKey := fmt.Sprintf("ticket:search:id:%s:event:%d:order:%d:price:%.2f:type:%s:seat:%s:page:%d:size:%d:sort:%s:order:%s",
		request.ID, request.EventID, request.OrderID, request.Price, ticketType, request.SeatNumber,
		request.Page, request.Size, request.Sort, request.Order)

	var cacheResponse model.Response[[]*model.TicketResponse]
	err := s.Cache.Get(cacheKey, &cacheResponse)
	if err == nil {
		return &cacheResponse, nil
	}

	opts := model.TicketQueryOptions{
		Page:  request.Page,
		Size:  request.Size,
		Sort:  request.Sort,
		Order: request.Order,
	}

	// Only set pointer fields if they have non-zero values
	if request.ID != "" {
		opts.ID = &request.ID
	}
	if request.EventID != 0 {
		opts.EventID = &request.EventID
	}
	if request.OrderID != 0 {
		opts.OrderID = &request.OrderID
	}
	if request.Price != 0 {
		opts.Price = &request.Price
	}
	if request.Type != "" {
		opts.Type = &ticketType
	}
	if request.SeatNumber != "" {
		opts.SeatNumbers = &[]string{request.SeatNumber}
	}

	db := s.DB.WithContext(ctx)

	tickets, err := s.TicketRepository.Find(db, &opts)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrNotFound
		}
		s.Log.Errorf("failed to get tickets: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if len(tickets) == 0 {
		return nil, domainErrors.ErrNotFound
	}

	response := converter.TicketsToPaginatedResponse(tickets, int64(len(tickets)), opts.Page, opts.Size)

	if err := s.Cache.Set(cacheKey, response, 5*time.Minute); err != nil {
		s.Log.Errorf("failed to set cache: %v", err)
	}

	return response, nil
}
