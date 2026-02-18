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

func TestInternalPayment(t *testing.T) {
	// Setup test environment
	db := SetupTestEnvironment(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Create test data
	user := CreateTestUser(t, db)
	fromAccount := CreateTestAccount(t, db, user.ID, 1000.0)
	toAccount := CreateTestAccount(t, db, user.ID, 500.0)

	// Create test router
	router := SetupTestRouterWithDB(db)
	token := CreateTestJWT(t, user.Email)

	// Test case 1: Successful internal payment
	t.Run("Successful Internal Payment", func(t *testing.T) {
		payload := structs.InternalPaymentRequest{
			FromAccount: fromAccount.AccountID,
			ToAccount:   toAccount.AccountID,
			Amount:      100.0,
			Currency:    "USD",
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/payment/api/internal_payment", bytes.NewBuffer(jsonPayload))
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

	// Test case 2: Insufficient balance
	t.Run("Insufficient Balance", func(t *testing.T) {
		payload := structs.InternalPaymentRequest{
			FromAccount: fromAccount.AccountID,
			ToAccount:   toAccount.AccountID,
			Amount:      2000.0, // More than available balance
			Currency:    "USD",
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/payment/api/internal_payment", bytes.NewBuffer(jsonPayload))
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

	// Test case 3: Invalid account
	t.Run("Invalid Account", func(t *testing.T) {
		payload := structs.InternalPaymentRequest{
			FromAccount: "invalid-account-id",
			ToAccount:   toAccount.AccountID,
			Amount:      100.0,
			Currency:    "USD",
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/payment/api/internal_payment", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "Account does not exist")
	})

	// Test case 4: Unauthorized access
	t.Run("Unauthorized Access", func(t *testing.T) {
		payload := structs.InternalPaymentRequest{
			FromAccount: fromAccount.AccountID,
			ToAccount:   toAccount.AccountID,
			Amount:      100.0,
			Currency:    "USD",
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/payment/api/internal_payment", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		// No authorization header

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
