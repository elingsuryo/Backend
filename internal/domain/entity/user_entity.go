package entity

import "gorm.io/gorm"

type User struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password" gorm:"not null"`
	Name     string `json:"name" gorm:"not null"`
	Role     string `json:"role" gorm:"not null"`
	Status   bool   `json:"status" gorm:"not null"`
	VerifyEmailToken   string `json:"verify_email_token" gorm:"null"`
	IsVerify           bool   `json:"is_verify" gorm:"not null"`
	ResetPasswordToken string `json:"reset_password_token" gorm:"null"`
	gorm.Model
}

func (u *User) TableName() string {
	return "users"
}
