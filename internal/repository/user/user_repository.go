package user

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/repository"
	"gorm.io/gorm"
)

type UserRepository interface {
	repository.Repository[entity.User]
	GetFirst(db *gorm.DB, user *entity.User) error
	GetByID(db *gorm.DB, user *entity.User, id string) error
	GetByEmail(db *gorm.DB, user *entity.User, email string) error
	CountByRole(db *gorm.DB, role string) (int64, error)
	GetByResetPasswordToken(db *gorm.DB, user *entity.User, token string) error
	GetByVerifyEmailToken(db *gorm.DB, user *entity.User, token string) error
}
