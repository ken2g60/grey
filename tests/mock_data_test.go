package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/grey/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// MockDataGenerator provides helper functions to generate mock data for testing
type MockDataGenerator struct {
	db *gorm.DB
}

func NewMockDataGenerator(db *gorm.DB) *MockDataGenerator {
	return &MockDataGenerator{db: db}
}

// GenerateTestUsers creates multiple test users with different roles
func (md *MockDataGenerator) GenerateTestUsers(t *testing.T, count int) []*models.User {
	users := make([]*models.User, count)

	for i := 0; i < count; i++ {
		user := &models.User{
			Email:    fmt.Sprintf("testuser%d-%d@example.com", i, time.Now().UnixNano()),
			Password: "hashedpassword123",
		}

		err := md.db.Create(user).Error
		if err != nil {
			t.Fatalf("Failed to create test user %d: %v", i, err)
		}

		users[i] = user
	}

	return users
}

// GenerateTestAccounts creates multiple test accounts with different balances
func (md *MockDataGenerator) GenerateTestAccounts(t *testing.T, userID int, count int, balances []float64) []*models.Account {
	accounts := make([]*models.Account, count)

	for i := 0; i < count; i++ {
		balance := 1000.0 // default balance
		if i < len(balances) {
			balance = balances[i]
		}

		account := &models.Account{
			UserID:   userID,
			Currency: "USD",
			Balance:  decimal.NewFromFloat(balance),
		}

		err := md.db.Create(account).Error
		if err != nil {
			t.Fatalf("Failed to create test account %d: %v", i, err)
		}

		accounts[i] = account
	}

	return accounts
}

// GenerateTestPayments creates test payment records
func (md *MockDataGenerator) GenerateTestPayments(t *testing.T, count int) []*models.Payment {
	payments := make([]*models.Payment, count)

	for i := 0; i < count; i++ {
		payment := &models.Payment{
			FromAccount: fmt.Sprintf("account-from-%d", i),
			ToAccount:   fmt.Sprintf("account-to-%d", i),
			Currency:    "USD",
			Amount:      decimal.NewFromFloat(float64((i + 1) * 100)),
			Status:      models.Pending,
		}

		err := md.db.Create(payment).Error
		if err != nil {
			t.Fatalf("Failed to create test payment %d: %v", i, err)
		}

		payments[i] = payment
	}

	return payments
}

// GenerateLedgerEntries creates test ledger entries
func (md *MockDataGenerator) GenerateLedgerEntries(t *testing.T, count int) []*models.LedgerEntry {
	entries := make([]*models.LedgerEntry, count)

	for i := 0; i < count; i++ {
		entry := &models.LedgerEntry{
			AccountID: fmt.Sprintf("account-%d", i),
			PaymentID: fmt.Sprintf("payment-%d", i),
			Amount:    decimal.NewFromFloat(float64((i + 1) * 50)),
		}

		err := md.db.Create(entry).Error
		if err != nil {
			t.Fatalf("Failed to create ledger entry %d: %v", i, err)
		}

		entries[i] = entry
	}

	return entries
}

// TestMockDataGeneration tests the mock data generation functions
func TestMockDataGeneration(t *testing.T) {
	db := SetupTestEnvironment(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	generator := NewMockDataGenerator(db)

	t.Run("Generate Test Users", func(t *testing.T) {
		users := generator.GenerateTestUsers(t, 3)
		assert.Len(t, users, 3)

		for i, user := range users {
			assert.NotZero(t, user.ID)
			assert.Contains(t, user.Email, fmt.Sprintf("testuser%d-", i))
			assert.Contains(t, user.Email, "@example.com")
		}
	})

	t.Run("Generate Test Accounts", func(t *testing.T) {
		user := generator.GenerateTestUsers(t, 1)[0]
		balances := []float64{500.0, 1000.0, 1500.0}
		accounts := generator.GenerateTestAccounts(t, user.ID, 3, balances)

		assert.Len(t, accounts, 3)

		for i, account := range accounts {
			assert.NotZero(t, account.ID)
			assert.Equal(t, user.ID, account.UserID)
			assert.Equal(t, "USD", account.Currency)
			assert.True(t, account.Balance.Equal(decimal.NewFromFloat(balances[i])))
		}
	})

	t.Run("Generate Test Payments", func(t *testing.T) {
		payments := generator.GenerateTestPayments(t, 3)
		assert.Len(t, payments, 3)

		for i, payment := range payments {
			assert.NotZero(t, payment.ID)
			assert.Equal(t, fmt.Sprintf("account-from-%d", i), payment.FromAccount)
			assert.Equal(t, fmt.Sprintf("account-to-%d", i), payment.ToAccount)
			assert.Equal(t, "USD", payment.Currency)
			assert.True(t, payment.Amount.Equal(decimal.NewFromFloat(float64((i+1)*100))))
			assert.Equal(t, models.Pending, payment.Status)
		}
	})

	t.Run("Generate Ledger Entries", func(t *testing.T) {
		entries := generator.GenerateLedgerEntries(t, 3)
		assert.Len(t, entries, 3)

		for i, entry := range entries {
			assert.NotZero(t, entry.ID)
			assert.Equal(t, fmt.Sprintf("account-%d", i), entry.AccountID)
			assert.Equal(t, fmt.Sprintf("payment-%d", i), entry.PaymentID)
			assert.True(t, entry.Amount.Equal(decimal.NewFromFloat(float64((i+1)*50))))
		}
	})
}
