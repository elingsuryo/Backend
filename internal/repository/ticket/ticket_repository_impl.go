package ticket

import (
	"fmt"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TicketRepositoryImpl struct {
	repository.RepositoryImpl[entity.Ticket]
	Log *logrus.Logger
}

func NewTicketRepository(db *gorm.DB, log *logrus.Logger) *TicketRepositoryImpl {
	return &TicketRepositoryImpl{
		RepositoryImpl: repository.RepositoryImpl[entity.Ticket]{DB: db},
		Log:            log,
	}
}

func (r *TicketRepositoryImpl) CreateBatch(db *gorm.DB, tickets []*entity.Ticket) error {
	return db.Create(tickets).Error
}

func (r *TicketRepositoryImpl) Find(db *gorm.DB, opts *model.TicketQueryOptions) ([]*entity.Ticket, error) {
	// First, get total count
	var totalCount int64
	countQuery := db.Model(&entity.Ticket{})

	// Apply filters to count query
	if opts.ID != nil {
		countQuery = countQuery.Where("UPPER(id) = UPPER(?)", *opts.ID)
	}
	if opts.EventID != nil {
		countQuery = countQuery.Where("event_id = ?", *opts.EventID)
	}
	if opts.OrderID != nil {
		countQuery = countQuery.Where("order_id = ?", *opts.OrderID)
	}
	if opts.Price != nil {
		countQuery = countQuery.Where("price = ?", *opts.Price)
	}
	if opts.SeatNumbers != nil && len(*opts.SeatNumbers) > 0 {
		countQuery = countQuery.Where("UPPER(seat_number) IN (?)", *opts.SeatNumbers)
	}

	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Main query with preloads
	query := db.Model(&entity.Ticket{}).
		Preload("Event").
		Preload("Order")

	// Apply same filters to main query
	if opts.ID != nil {
		query = query.Where("UPPER(id) = UPPER(?)", *opts.ID)
	}
	if opts.EventID != nil {
		query = query.Where("event_id = ?", *opts.EventID)
	}
	if opts.OrderID != nil {
		query = query.Where("order_id = ?", *opts.OrderID)
	}
	if opts.Price != nil {
		query = query.Where("price = ?", *opts.Price)
	}
	if opts.SeatNumbers != nil && len(*opts.SeatNumbers) > 0 {
		query = query.Where("UPPER(seat_number) IN (?)", *opts.SeatNumbers)
	}

	// Apply pagination
	if opts.Page > 0 && opts.Size > 0 {
		offset := (opts.Page - 1) * opts.Size
		query = query.Offset(offset).Limit(opts.Size)
	}

	// Apply sorting
	if opts.Sort != "" && opts.Order != "" {
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, opts.Order))
	}

	var tickets []*entity.Ticket
	if err := query.Find(&tickets).Error; err != nil {
		return nil, err
	}

	// Store total count in the first ticket's metadata (if any tickets exist)
	if len(tickets) > 0 {
		tickets[0].Metadata = map[string]interface{}{
			"total_count": totalCount,
		}
	}

	return tickets, nil
}

func (r *TicketRepositoryImpl) GetLastTicketNumber(db *gorm.DB, eventID uint, ticketType string) (*entity.Ticket, error) {
	var ticket entity.Ticket

	err := db.Transaction(func(tx *gorm.DB) error {
		// Lock the tickets table
		if err := tx.Exec("LOCK TABLE tickets IN SHARE ROW EXCLUSIVE MODE").Error; err != nil {
			return err
		}

		// Get the last ticket number for the specific event and type
		result := tx.Where("event_id = ? AND type = ?", eventID, ticketType).
			Order("CAST(SPLIT_PART(seat_number, '-', 2) AS INTEGER) DESC").
			First(&ticket)

		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &ticket, nil
}
