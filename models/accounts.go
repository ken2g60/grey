package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Account struct {
	ID        int             `json:"id" gorm:"type:integer;primaryKey"`
	UserID    int             `json:"user_id" gorm:"not null;index"`
	AccountID string          `json:"account_id" gorm:"type:uuid;not null;index"` // account number
	Currency  string          `json:"currency" gorm:"type:varchar(3);not null"`
	Balance   decimal.Decimal `json:"balance" gorm:"type:numeric(18,2);not null;default:0"`
	CreatedAt time.Time
}

func (account *Account) BeforeCreate(tx *gorm.DB) error {
	account.AccountID = uuid.NewString()
	return nil
}

func CreateAccount(ctx context.Context, db *gorm.DB, account *Account) (err error) {
	err = db.WithContext(ctx).Create(&account).Error
	if err != nil {
		return err
	}
	return nil
}

func IsAccountExists(ctx context.Context, db *gorm.DB, uuid string) (*Account, error) {
	var account Account
	err := db.Debug().WithContext(ctx).Model(&account).Where("account_id = ?", uuid).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func UpdateAccountBalance(ctx context.Context, db *gorm.DB, accountID string, newBalance decimal.Decimal) error {
	return db.WithContext(ctx).Model(&Account{}).Where("account_id = ?", accountID).Update("balance", newBalance).Error
}
