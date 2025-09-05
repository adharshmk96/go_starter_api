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
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Email     string         `json:"email" gorm:"unique"`
	Password  string         `json:"password"`
}

type AccountActivity struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	AccountID uint   `json:"account_id"`
	Activity  string `json:"activity"`
}

type AccountService interface {
	GenerateAuthToken(account *Account) (string, error)
	ValidateAuthToken(token string) (uint, error)
	HashPassword(password string) (string, error)
	ComparePassword(password, hash string) (bool, error)

	GeneratePasswordResetToken(account *Account) (string, error)
	ValidatePasswordResetToken(token string) (uint, error)
	SendPasswordResetEmail(email string, token string) error
}

var (
	ErrPasswordEmpty     = errors.New("password cannot be empty")
	ErrInvalidHashFormat = errors.New("invalid hash format")
	ErrServerURLNotSet   = errors.New("server url is not set")
)

var (
	ActivityLogin          = "login"
	ActivityLogout         = "logout"
	ActivityRegister       = "register"
	ActivityUpdate         = "update"
	ActivityDelete         = "delete"
	ActivityResetPassword  = "reset_password"
	ActivityForgotPassword = "forgot_password"
	ActivityChangePassword = "change_password"
)

type AccountRepository interface {
	CreateAccount(account *Account) (*Account, error)
	GetAccountByEmail(email string) (*Account, error)
	GetAccountByID(id uint) (*Account, error)
	UpdateAccount(account *Account) (*Account, error)
	DeleteAccount(id uint) error

	LogAccountActivity(accountID uint, activity string) error
}
