package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	Pending   PaymentStatus = "pending"
	Completed PaymentStatus = "completed"
	Failed    PaymentStatus = "failed"
)

type Payment struct {
	ID          int             `json:"id" gorm:"type:integer;primaryKey"`
	PaymentID   string          `json:"payment_id" gorm:"type:uuid;not null;index"`
	FromAccount string          `json:"from_account" gorm:"type:uuid;not null"`
	ToAccount   string          `json:"to_account" gorm:"type:uuid"`
	Currency    string          `json:"currency" gorm:"type:varchar(3);not null"`
	Amount      decimal.Decimal `json:"amount" gorm:"type:numeric(18,2);not null"`
	Status      PaymentStatus   `json:"status" gorm:"type:varchar(20);default:'pending'"`
	Description string          `json:"description" gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	p.PaymentID = uuid.NewString()
	return nil
}

func CreatePayment(ctx context.Context, db *gorm.DB, payment *Payment) error {
	return db.WithContext(ctx).Create(payment).Error
}
