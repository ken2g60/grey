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

func TestTopUp(t *testing.T) {
	// Setup test environment
	db := SetupTestEnvironment(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Create test data
	user := CreateTestUser(t, db)
	account := CreateTestAccount(t, db, user.ID, 500.0)

	// Create test router
	router := SetupTestRouterWithDB(db)
	token := CreateTestJWT(t, user.Email)

	// Test case 1: Successful top-up
	t.Run("Successful Top Up", func(t *testing.T) {
		payload := structs.TopUp{
			Account:  account.AccountID,
			Amount:   200.0,
			Currency: "USD",
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/payment/api/topup", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "topup successfully", response["message"])

		// Verify the response contains payment ID
		responseData := response["response"].(map[string]interface{})
		assert.NotEmpty(t, responseData["payment_id"])
		assert.Equal(t, "success", responseData["status"])
	})

	// Test case 2: Invalid amount (zero)
	t.Run("Invalid Amount - Zero", func(t *testing.T) {
		payload := structs.TopUp{
			Account:  account.AccountID,
			Amount:   0.0,
			Currency: "USD",
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/payment/api/topup", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to process payment", response["message"])
	})

	// Test case 3: Invalid amount (negative)
	t.Run("Invalid Amount - Negative", func(t *testing.T) {
		payload := structs.TopUp{
			Account:  account.AccountID,
			Amount:   -100.0,
			Currency: "USD",
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/payment/api/topup", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to process payment", response["message"])
	})

	// Test case 4: Account not found
	t.Run("Account Not Found", func(t *testing.T) {
		payload := structs.TopUp{
			Account:  "non-existent-account",
			Amount:   100.0,
			Currency: "USD",
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/payment/api/topup", bytes.NewBuffer(jsonPayload))
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

	// Test case 5: Unauthorized access
	t.Run("Unauthorized Access", func(t *testing.T) {
		payload := structs.TopUp{
			Account:  account.AccountID,
			Amount:   100.0,
			Currency: "USD",
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/payment/api/topup", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		// No authorization header

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Test case 6: Large amount top-up
	t.Run("Large Amount Top Up", func(t *testing.T) {
		payload := structs.TopUp{
			Account:  account.AccountID,
			Amount:   10000.0,
			Currency: "USD",
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/payment/api/topup", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "topup successfully", response["message"])
	})
}
