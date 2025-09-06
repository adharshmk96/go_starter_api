package account

import (
	"servicehub_api/pkg/domain"

	"gorm.io/gorm"
)

type AccountRepo struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) domain.AccountRepository {
	return &AccountRepo{db: db}
}

func (r *AccountRepo) CreateAccount(account *domain.Account) (*domain.Account, error) {
	err := r.db.Create(account).Error
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (r *AccountRepo) GetAccountByEmail(email string) (*domain.Account, error) {
	var account domain.Account
	err := r.db.Where("email = ?", email).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepo) GetAccountByID(id uint) (*domain.Account, error) {
	var account domain.Account
	err := r.db.Where("id = ?", id).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepo) UpdateAccount(account *domain.Account) (*domain.Account, error) {
	err := r.db.Save(account).Error
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (r *AccountRepo) DeleteAccount(id uint) error {
	return r.db.Delete(&domain.Account{}, id).Error
}

func (r *AccountRepo) LogAccountActivity(accountID uint, activity string) error {
	return r.db.Create(&domain.AccountActivity{AccountID: accountID, Activity: activity}).Error
}
