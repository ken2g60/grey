# Payment API Tests

This directory contains comprehensive test suites for the payment API endpoints including internal payments, external payments, and top-up functionality.

## Test Structure

### Test Files

1. **setup_test.go** - Common test setup utilities and helper functions
2. **internal_payment_test.go** - Tests for internal payment endpoint
3. **external_payment_test.go** - Tests for external payment endpoint  
4. **topup_test.go** - Tests for top-up endpoint
5. **service_test.go** - Tests for payment service layer functions
6. **mock_data_test.go** - Mock data generation utilities and tests
7. **all_tests.go** - Test runner that executes all payment tests

### Test Coverage

#### Internal Payment Tests (`/payment/api/internal_payment`)
- ✅ Successful internal payment
- ✅ Insufficient balance scenarios
- ✅ Invalid account handling
- ✅ Unauthorized access

#### External Payment Tests (`/payment/api/external_payment`)
- ✅ Successful bank transfer
- ✅ Successful mobile money transfer
- ✅ Insufficient balance scenarios
- ✅ Invalid transaction type handling
- ✅ Account not found scenarios
- ✅ Unauthorized access

#### Top-Up Tests (`/payment/api/topup`)
- ✅ Successful top-up
- ✅ Invalid amount (zero/negative)
- ✅ Account not found scenarios
- ✅ Unauthorized access
- ✅ Large amount top-ups

#### Service Layer Tests
- ✅ `ProcessInternalPayment` function
- ✅ `ProcessExternalPayment` function
- ✅ `TopUpProcess` function
- ✅ Error handling for invalid amounts
- ✅ Insufficient balance handling

#### Mock Data Generation
- ✅ Test user generation
- ✅ Test account generation
- ✅ Test payment generation
- ✅ Ledger entry generation

## Running Tests

### Run All Tests
```bash
go test ./tests/...
```

### Run Tests with Verbose Output
```bash
go test ./tests/... -v
```