# API Documentation

## Overview

The Grey FinTech API provides a comprehensive set of endpoints for user management, account operations, and payment processing. All endpoints use JSON for request/response formats and follow RESTful conventions.

## Base URL

```
http://localhost:8000
```

## Authentication

The API uses JWT (JSON Web Token) authentication for protected endpoints. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

**Token Expiration**: 30 minutes

## Response Format

All API responses follow a consistent structure:

### Success Response
```json
{
  "message": "Success message",
  "data": { ... },
  "status": 200
}
```

### Error Response
```json
{
  "message": "Error description",
  "status": 400
}
```

## Endpoints

### User Management

#### Register User
Creates a new user account and automatically creates a USD account.

**Endpoint**: `POST /user/api/register`

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response**:
```json
{
  "message": "user registered successfully",
  "data": "user registered successfully",
  "status": 201
}
```

**Validation Rules**:
- Email must be unique and valid format
- Password minimum length: 6 characters
- Both fields are required

#### Login
Authenticates a user and returns a JWT token.

**Endpoint**: `POST /user/api/login`

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response**:
```json
{
  "message": "Welcome back! You have successfully logged in.",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Error Responses**:
- `400`: Invalid credentials or missing fields
- `401`: User not found or incorrect password

### Payment Processing

All payment endpoints require JWT authentication.

#### Internal Payment
Transfers funds between two accounts within the system.

**Endpoint**: `POST /payment/api/internal_payment`

**Headers**:
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body**:
```json
{
  "from_account": "account-uuid-1",
  "to_account": "account-uuid-2", 
  "amount": 100.50,
  "currency": "USD"
}
```

**Response**:
```json
{
  "response": {
    "id": 1,
    "payment_id": "payment-uuid",
    "from_account": "account-uuid-1",
    "to_account": "account-uuid-2",
    "currency": "USD",
    "amount": "100.50",
    "status": "pending",
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  },
  "message": "Payment created successfully"
}
```

**Validation Rules**:
- Both accounts must exist
- Sufficient balance in source account
- Amount must be positive
- Currency must be valid (3-letter code)

#### External Payment
Processes payments to external recipients (bank transfers or mobile money).

**Endpoint**: `POST /payment/api/external_payment`

**Headers**:
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body**:
```json
{
  "from_account": "account-uuid",
  "amount": 250.00,
  "currency": "USD",
  "transaction_type": "BANK_TRANSFER",
  "recipient": {
    "recipientNumber": "1234567890",
    "recipientName": "John Doe"
  }
}
```

**Response**:
```json
{
  "response": {
    "payment_id": "payment-uuid",
    "recipient": {
      "recipientNumber": "1234567890",
      "recipientName": "John Doe"
    },
    "status": "success",
    "provider_status": "processing"
  },
  "message": "Payment created successfully"
}
```

**Transaction Types**:
- `BANK_TRANSFER`: Transfer to bank account
- `MOBILE_MONEY`: Transfer to mobile money service

**Validation Rules**:
- Source account must exist
- Sufficient balance in source account
- Amount must be positive
- Valid transaction type
- Recipient details required

## Error Handling

### Common Error Codes

| Status Code | Description | Example |
|-------------|-------------|---------|
| 400 | Bad Request | Invalid input data |
| 401 | Unauthorized | Missing or invalid token |
| 404 | Not Found | Account or user not found |
| 500 | Internal Server Error | Database or system error |

### Error Response Format

```json
{
  "message": "Error description",
  "status": 400
}
```

### Common Error Messages

- `"Provide all the required fields"` - Missing required input fields
- `"email already exists"` - Email already registered
- `"Password is incorrect"` - Invalid login credentials
- `"Account does not exist"` - Account not found
- `"Insufficient balance"` - Not enough funds for transfer
- `"Invalid transaction type support BANK_TRANSFER or MOBILE_MONEY"` - Unsupported transaction type

## Data Types

### UUID Format
All account IDs, user IDs, and payment IDs use UUID format:
```
"550e8400-e29b-41d4-a716-446655440000"
```

### Decimal Format
All monetary values use decimal format with 2 decimal places:
```json
"amount": "100.50"
```

### Currency Codes
3-letter ISO 4217 currency codes:
```json
"currency": "USD"
```

## Security Considerations

1. **Token Management**: Tokens expire after 30 minutes
2. **HTTPS**: Use HTTPS in production environments
3. **Input Validation**: All inputs are validated and sanitized
4. **Rate Limiting**: Consider implementing rate limiting for production
5. **Transaction Safety**: All financial operations use database transactions

## Testing

### Example cURL Commands

#### Register User
```bash
curl -X POST http://localhost:8000/user/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

#### Login
```bash
curl -X POST http://localhost:8000/user/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

#### Internal Payment
```bash
curl -X POST http://localhost:8000/payment/api/internal_payment \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "from_account": "account-uuid-1",
    "to_account": "account-uuid-2",
    "amount": 100.50,
    "currency": "USD"
  }'
```

## Rate Limiting

Currently not implemented, but recommended for production:
- 100 requests per minute per IP
- 10 payment requests per minute per user

## Webhooks

Webhooks are not currently implemented but planned for:
- Payment status updates
- Account balance notifications
- Transaction confirmations

## SDKs

SDKs are not currently available but planned for:
- JavaScript/Node.js
- Python
- Go

## Support

For API support and questions:
- Create an issue in the GitHub repository
- Check the documentation for common solutions
- Review error messages for specific issues
