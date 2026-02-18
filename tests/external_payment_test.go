package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grey/structs"
	"github.com/stretchr/testify/assert"
)

func TestExternalPayment(t *testing.T) {
	// Setup test environment
	db := SetupTestEnvironment(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Create test data
	user := CreateTestUser(t, db)
	fromAccount := CreateTestAccount(t, db, user.ID, 1000.0)

	// Create test router
	router := SetupTestRouterWithDB(db)
	token := CreateTestJWT(t, user.Email)

	// Test case 1: Successful bank transfer
	t.Run("Successful Bank Transfer", func(t *testing.T) {
		payload := structs.ExternalPaymentRequest{
			Account:         fromAccount.AccountID,
			Amount:          200.0,
			Currency:        "USD",
			TransactionType: "BANK_TRANSFER",
			Recipient: structs.RecipientDetails{
				RecipientNumber: "1234567890",
				RecipientName:   "John Doe",
			},
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/payment/api/external_payment", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Payment created successfully", response["message"])
	})

	// Test case 2: Successful mobile money transfer
	t.Run("Successful Mobile Money Transfer", func(t *testing.T) {
		payload := structs.ExternalPaymentRequest{
			Account:         fromAccount.AccountID,
			Amount:          150.0,
			Currency:        "USD",
			TransactionType: "MOBILE_MONEY",
			Recipient: structs.RecipientDetails{
				RecipientNumber: "0771234567",
				RecipientName:   "Jane Smith",
			},
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/payment/api/external_payment", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Payment created successfully", response["message"])
	})

	// Test case 3: Insufficient balance
	t.Run("Insufficient Balance", func(t *testing.T) {
		payload := structs.ExternalPaymentRequest{
			Account:         fromAccount.AccountID,
			Amount:          2000.0, // More than available balance
			Currency:        "USD",
			TransactionType: "BANK_TRANSFER",
			Recipient: structs.RecipientDetails{
				RecipientNumber: "1234567890",
				RecipientName:   "John Doe",
			},
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/payment/api/external_payment", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Insufficient balance", response["message"])
	})

	// Test case 4: Invalid transaction type
	t.Run("Invalid Transaction Type", func(t *testing.T) {
		payload := structs.ExternalPaymentRequest{
			Account:         fromAccount.AccountID,
			Amount:          100.0,
			Currency:        "USD",
			TransactionType: "INVALID_TYPE",
			Recipient: structs.RecipientDetails{
				RecipientNumber: "1234567890",
				RecipientName:   "John Doe",
			},
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/payment/api/external_payment", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "Invalid transaction type")
	})

	// Test case 5: Account not found
	t.Run("Account Not Found", func(t *testing.T) {
		payload := structs.ExternalPaymentRequest{
			Account:         "non-existent-account",
			Amount:          100.0,
			Currency:        "USD",
			TransactionType: "BANK_TRANSFER",
			Recipient: structs.RecipientDetails{
				RecipientNumber: "1234567890",
				RecipientName:   "John Doe",
			},
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/payment/api/external_payment", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Account not found", response["message"])
	})

	// Test case 6: Unauthorized access
	t.Run("Unauthorized Access", func(t *testing.T) {
		payload := structs.ExternalPaymentRequest{
			Account:         fromAccount.AccountID,
			Amount:          100.0,
			Currency:        "USD",
			TransactionType: "BANK_TRANSFER",
			Recipient: structs.RecipientDetails{
				RecipientNumber: "1234567890",
				RecipientName:   "John Doe",
			},
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/payment/api/external_payment", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		// No authorization header

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
