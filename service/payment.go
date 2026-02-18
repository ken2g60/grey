package service

import (
	"context"
	"errors"

	"github.com/grey/models"
	"github.com/grey/structs"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func ProcessInternalPayment(ctx context.Context, DB *gorm.DB, fromAccount, toAccount string, amount decimal.Decimal, currency string) (Payment models.Payment, err error) {

	if amount.LessThanOrEqual(decimal.Zero) {
		return Payment, errors.New("invalid amount")
	}

	tx := DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return Payment, tx.Error
	}
	// Defer rollback in case of error
	defer tx.Rollback()

	// 2. Deduct from source (Lock row)
	result := tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE account_id = $2", amount, fromAccount)
	if result.Error != nil {
		return Payment, result.Error // Likely insufficient funds or database error
	}
	if result.RowsAffected == 0 {
		return Payment, errors.New("insufficient balance")
	}

	// 3. Add to target
	tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE account_id = $2", amount, toAccount)
	if tx.Error != nil {
		return Payment, tx.Error
	}

	response := models.Payment{
		FromAccount: fromAccount,
		ToAccount:   toAccount,
		Currency:    currency,
		Amount:      amount,
		Status:      models.Pending,
		Description: "Internal Payment",
	}
	err = tx.Create(&response).Error
	if err != nil {
		return Payment, err
	}

	tx.Commit()
	return response, nil
}

func ProcessExternalPayment(ctx context.Context, DB *gorm.DB, recipient structs.RecipientDetails, fromAccount string, amount decimal.Decimal, currency string) (response structs.ExternalPaymentResponse, err error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return structs.ExternalPaymentResponse{}, errors.New("invalid amount")
	}

	tx := DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return structs.ExternalPaymentResponse{}, tx.Error
	}
	// Defer rollback in case of error
	defer tx.Rollback()

	// 2. Deduct from source (Lock row)
	result := tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE account_id = $2 AND balance >= $1", amount, fromAccount)
	if result.Error != nil {
		return structs.ExternalPaymentResponse{}, result.Error // Likely insufficient funds or database error
	}
	if result.RowsAffected == 0 {
		return structs.ExternalPaymentResponse{}, errors.New("insufficient balance")
	}

	payment := models.Payment{
		FromAccount: fromAccount,
		ToAccount:   fromAccount,
		Currency:    currency,
		Amount:      amount,
		Status:      models.Completed,
		Description: "External Payment",
	}

	err = tx.Create(&payment).Error
	if err != nil {
		return structs.ExternalPaymentResponse{}, err
	}

	tx.Commit()
	return structs.ExternalPaymentResponse{
		PaymentID:      payment.PaymentID,
		Recipient:      recipient,
		Status:         "success",
		ProviderStatus: "success",
	}, nil
}

func TopUpProcess(ctx context.Context, DB *gorm.DB, fromAccount string, amount decimal.Decimal, currency string) (response structs.TopUpResponse, err error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return structs.TopUpResponse{}, errors.New("invalid amount")
	}

	tx := DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return structs.TopUpResponse{}, tx.Error
	}
	// Defer rollback in case of error
	defer tx.Rollback()

	// 2. Add funds to account
	tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE account_id = $2", amount, fromAccount)
	if tx.Error != nil {
		return structs.TopUpResponse{}, tx.Error // Likely insufficient funds or database error
	}

	payment := models.Payment{
		FromAccount: fromAccount,
		ToAccount:   fromAccount,
		Amount:      amount,
		Currency:    currency,
		Status:      models.Completed,
		Description: "Top up",
	}

	err = tx.Create(&payment).Error
	if err != nil {
		return structs.TopUpResponse{}, err
	}

	// Create ledger entry
	ledgerEntry := models.LedgerEntry{
		AccountID: fromAccount,
		PaymentID: payment.PaymentID,
		Amount:    amount,
	}
	err = tx.Create(&ledgerEntry).Error
	if err != nil {
		return structs.TopUpResponse{}, err
	}

	tx.Commit()
	return structs.TopUpResponse{
		PaymentID: payment.PaymentID,
		Status:    "success",
	}, nil
}
