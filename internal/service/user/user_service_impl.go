package user

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"time"

	"embed"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/domain/model/converter"
	"github.com/TrinityKnights/Backend/internal/repository/user"
	domainErrors "github.com/TrinityKnights/Backend/pkg/errors"
	"github.com/TrinityKnights/Backend/pkg/gomail"
	"github.com/TrinityKnights/Backend/pkg/helper"
	"github.com/TrinityKnights/Backend/pkg/jwt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//go:embed template/*.html
var templateFS embed.FS

type UserServiceImpl struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *user.UserRepositoryImpl
	JWTService     jwt.JWTService
	Gomail         *gomail.ImplGomail
	helper         *helper.ContextHelper
}

func NewUserServiceImpl(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *user.UserRepositoryImpl, jwtService jwt.JWTService, mail *gomail.ImplGomail) *UserServiceImpl {
	return &UserServiceImpl{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
		JWTService:     jwtService,
		Gomail:         mail,
		helper:         helper.NewContextHelper(),
	}
}

func (s *UserServiceImpl) Register(ctx context.Context, request *model.RegisterRequest) (*model.UserResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	existingUser := &entity.User{}
	err := s.UserRepository.GetByEmail(tx, existingUser, request.Email)
	if err == nil {
		s.Log.Errorf("email already exists: %v", request.Email)
		return nil, domainErrors.ErrEmailAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Log.Errorf("failed to check existing user: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	first := &entity.User{}
	if err := s.UserRepository.GetFirst(tx, first); err != nil {
		s.Log.Errorf("failed to get first user: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrNotFound
		}
		return nil, domainErrors.ErrInternalServer
	}

	var role string
	if first.ID == "" {
		role = "admin"
	} else {
		role = "buyer"
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Log.Errorf("failed to generate password: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	token := uuid.NewString()
	var replaceEmail = struct {
		Token string
	}{
		Token: token,
	}

	data := &entity.User{
		ID:               uuid.NewString(),
		Email:            request.Email,
		Password:         string(password),
		Name:             request.Name,
		Role:             role,
		VerifyEmailToken: token,
		IsVerified:       false,
		Status:           true,
	}

	if err := s.UserRepository.Create(tx, data); err != nil {
		s.Log.Errorf("failed to create user: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	tmpl, err := template.ParseFS(templateFS, "template/verify-email.html")
	if err != nil {
		s.Log.Errorf("failed to parse template: %v", err)
		return nil, domainErrors.ErrInternalServer
	}
	var body bytes.Buffer
	if err := tmpl.Execute(&body, &replaceEmail); err != nil {
		s.Log.Errorf("failed to execute template: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	emailRequest := &gomail.SendEmail{
		EmailTo:   data.Email,
		EmailFrom: s.Gomail.GetFromEmail(),
		Subject:   "[TrinityKnights] Verify Email",
		Body:      body,
	}
	if err := s.Gomail.SendEmail(emailRequest); err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("failed to commit transaction: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	return converter.UserToResponse(data), nil
}

func (s *UserServiceImpl) Login(ctx context.Context, request *model.LoginRequest) (*model.TokenResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	data := &entity.User{}
	if err := s.UserRepository.GetByEmail(tx, data, request.Email); err != nil {
		s.Log.Errorf("failed to get user by email: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrUnauthorized
		}
		return nil, domainErrors.ErrInternalServer
	}

	if err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(request.Password)); err != nil {
		s.Log.Errorf("failed to compare password: %v", err)
		return nil, domainErrors.ErrUnauthorized
	}

	accessToken, err := s.JWTService.GenerateAccessToken(data.ID, data.Email, data.Role)
	if err != nil {
		s.Log.Errorf("failed to generate access token: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	refreshToken, err := s.JWTService.GenerateRefreshToken(data.ID, data.Email, data.Role)
	if err != nil {
		s.Log.Errorf("failed to generate refresh token: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	t := time.Now()
	data.LastLogin = &t
	if err := s.UserRepository.Update(tx, data); err != nil {
		s.Log.Errorf("failed to update user: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("failed to commit transaction: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	return converter.LoginToTokenResponse(accessToken, refreshToken), nil
}

func (s *UserServiceImpl) Profile(ctx context.Context) (*model.UserResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	claims, err := s.helper.GetJWTClaims(ctx)
	if err != nil {
		s.Log.Errorf("failed to get jwt claims: %v", err)
		return nil, domainErrors.ErrUnauthorized
	}

	data := &entity.User{}
	if err := s.UserRepository.GetByID(tx, data, claims.UserID); err != nil {
		s.Log.Errorf("failed to get user by id: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrNotFound
		}
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("failed to commit transaction: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	return converter.UserToResponse(data), nil
}

func (s *UserServiceImpl) Update(ctx context.Context, request *model.UpdateRequest) (*model.UserResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	claims, err := s.helper.GetJWTClaims(ctx)
	if err != nil {
		s.Log.Errorf("failed to get jwt claims: %v", err)
		return nil, domainErrors.ErrUnauthorized
	}

	data := &entity.User{}
	if err := s.UserRepository.GetByID(tx, data, claims.UserID); err != nil {
		s.Log.Errorf("failed to get user by id: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrNotFound
		}
		return nil, domainErrors.ErrInternalServer
	}

	if request.Email != "" {
		data.Email = request.Email
	}

	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			s.Log.Errorf("failed to generate password: %v", err)
			return nil, domainErrors.ErrInternalServer
		}
		data.Password = string(password)
	}

	if request.Name != "" {
		data.Name = request.Name
	}

	if err := s.UserRepository.Update(tx, data); err != nil {
		s.Log.Errorf("failed to update user: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("failed to commit transaction: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	return converter.UserToResponse(data), nil
}

func (s *UserServiceImpl) RefreshToken(ctx context.Context, request *model.RefreshTokenRequest) (*model.TokenResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	claims, err := s.JWTService.ValidateToken(request.RefreshToken)
	if err != nil {
		s.Log.Errorf("failed to validate token: %v", err)
		return nil, domainErrors.ErrUnauthorized
	}

	accessToken, err := s.JWTService.GenerateAccessToken(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		s.Log.Errorf("failed to generate access token: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	refreshToken, err := s.JWTService.GenerateRefreshToken(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		s.Log.Errorf("failed to generate refresh token: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	return converter.LoginToTokenResponse(accessToken, refreshToken), nil
}

func (s *UserServiceImpl) RequestReset(ctx context.Context, request *model.ReqResetPasswordRequest) (*model.VerifyResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	existingUser := &entity.User{}
	if err := s.UserRepository.GetByEmail(tx, existingUser, request.Email); err != nil {
		s.Log.Errorf("email already exists: %v", request.Email)
		return nil, domainErrors.ErrNotFound
	}

	token := uuid.NewString()

	update := &entity.User{
		ID:                 existingUser.ID,
		Email:              existingUser.Email,
		Password:           existingUser.Password,
		Name:               existingUser.Name,
		Role:               existingUser.Role,
		Status:             existingUser.Status,
		ResetPasswordToken: token,
	}
	if err := s.UserRepository.Update(tx, update); err != nil {
		s.Log.Errorf("failed to update user: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("failed to commit transaction: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	var replaceEmail = struct {
		Token string
	}{
		Token: token,
	}

	tmpl, err := template.ParseFS(templateFS, "template/reset-password.html")
	if err != nil {
		s.Log.Errorf("failed to parse template: %v", err)
		return nil, domainErrors.ErrInternalServer
	}
	var body bytes.Buffer
	if err := tmpl.Execute(&body, &replaceEmail); err != nil {
		s.Log.Errorf("failed to execute template: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	emailRequest := &gomail.SendEmail{
		EmailTo:   existingUser.Email,
		EmailFrom: s.Gomail.GetFromEmail(),
		Subject:   "[TrinityKnights] Reset Password",
		Body:      body,
	}
	if err := s.Gomail.SendEmail(emailRequest); err != nil {
		return nil, err
	}

	return &model.VerifyResponse{Status: "success"}, nil
}

func (s *UserServiceImpl) ResetPassword(ctx context.Context, request *model.ResetPasswordRequest) (*model.VerifyResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	u := &entity.User{}
	if err := s.UserRepository.GetByResetPasswordToken(tx, u, request.Token); err != nil {
		s.Log.Errorf("failed to get user by reset password token: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrNotFound
		}
		return nil, domainErrors.ErrInternalServer
	}

	if request.Token != u.ResetPasswordToken {
		return nil, domainErrors.ErrBadRequest
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		s.Log.Errorf("failed to generate password: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	u.Password = string(hashedPassword)
	err = s.UserRepository.Update(tx, u)
	if err != nil {
		s.Log.Errorf("failed to update user: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("failed to commit transaction: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	return &model.VerifyResponse{Status: "success"}, nil
}

func (s *UserServiceImpl) VerifyEmail(ctx context.Context, request *model.VerifyRequest) (*model.VerifyResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrBadRequest
	}

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	u := &entity.User{}
	if err := s.UserRepository.GetByVerifyEmailToken(tx, u, request.Token); err != nil {
		s.Log.Errorf("failed to get user by reset password token: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrNotFound
		}
		return nil, domainErrors.ErrInternalServer
	}

	if request.Token != u.VerifyEmailToken {
		return nil, domainErrors.ErrBadRequest
	}

	u.IsVerified = true
	if err := s.UserRepository.Update(tx, u); err != nil {
		s.Log.Errorf("failed to update user: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("failed to commit transaction: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	return &model.VerifyResponse{Status: "success"}, nil
}
