package domain

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Email     string         `json:"email" gorm:"unique"`
	Password  string         `json:"password"`
}

type AccountService interface {
	GenerateToken(account *Account) (string, error)
	ValidateToken(token string) (uint, error)
	HashPassword(password string) (string, error)
	ComparePassword(password, hash string) (bool, error)
}

var (
	ErrPasswordEmpty     = errors.New("password cannot be empty")
	ErrInvalidHashFormat = errors.New("invalid hash format")
)

type AccountRepository interface {
	CreateAccount(account *Account) (*Account, error)
	GetAccountByEmail(email string) (*Account, error)
	GetAccountByID(id string) (*Account, error)
	UpdateAccount(account *Account) (*Account, error)
	DeleteAccount(id string) error
}
