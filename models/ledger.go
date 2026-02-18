package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type LedgerEntry struct {
	ID        int             `gorm:"primaryKey;autoIncrement"`
	AccountID string          `gorm:"not null;index"`
	Account   Account         `gorm:"foreignKey:AccountID;references:AccountID"`
	PaymentID string          `gorm:"not null;index"`
	Payment   Payment         `gorm:"foreignKey:PaymentID;references:PaymentID"`
	Amount    decimal.Decimal `gorm:"type:numeric(18,2);not null"`
	CreatedAt time.Time
}
