# Grey - FinTech Payment API

A comprehensive Go-based financial technology API that provides secure payment processing, user management, and account handling capabilities. Built with Gin framework, PostgreSQL, and JWT authentication.

## 🚀 Features

- **User Management**: Secure user registration and authentication with JWT tokens
- **Account Management**: Multi-currency account support with precise decimal calculations
- **Payment Processing**: Internal transfers and external payments (Bank Transfer & Mobile Money)
- **Security**: JWT-based authentication, bcrypt password hashing, and transaction-safe operations
- **Database**: PostgreSQL with GORM ORM and automatic migrations
- **API Documentation**: RESTful API with structured responses and error handling

## 📋 Prerequisites

- Go 1.25.0 or higher
- PostgreSQL 14 or higher
- Docker and Docker Compose (for containerized deployment)

## 🛠️ Installation

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/grey/grey.git
   cd grey
   ```

### Docker Deployment

1. **Using Docker Compose**
   ```bash
   # start all services
   docker-compose up -d

  # View logs
  docker-compose logs -f

  # Stop services
  docker-compose down
   ```

### Run All Tests
```bash
go test ./tests/...
```

### Testing with Insomina
Import the Insomina file from the root directory of the project.


## 🏗️ Architecture

### Data Models

#### User
- `ID`: Primary key
- `UserId`: Unique UUID identifier
- `Email`: Unique email address
- `Password`: Bcrypt-hashed password
- `CreatedAt`: Registration timestamp

#### Account
- `ID`: Primary key
- `AccountID`: Unique account number (UUID)
- `User`: Foreign key to User model
- `Currency`: 3-letter currency code (e.g., USD)
- `Balance`: Decimal balance with 18,2 precision
- `CreatedAt`: Account creation timestamp

#### Payment
- `ID`: Primary key
- `PaymentID`: Unique payment identifier (UUID)
- `FromAccount`: Source account ID
- `ToAccount`: Destination account ID (null for external payments)
- `Currency`: Payment currency
- `Amount`: Payment amount
- `Status`: Payment status (pending, completed, failed)
- `CreatedAt/UpdatedAt`: Timestamps

## Feature Implementation

### Multi-Account Processing

#### Account Management Features
The system supports users managing multiple currency accounts:

**Multi-Account Capabilities:**
- **Multiple Currencies**: Users can hold accounts in different currencies
- **Account Switching**: Easy switching between accounts for transactions
- **Balance Aggregation**: Total balance view with FX conversion to base currency
- **Account-to-Account Transfers**: Internal transfers between user's own accounts
- **Currency Conversion**: Automatic FX conversion when transferring between currency accounts

**Account Management API:**
```http
GET /accounts                    # Get all user accounts
POST /accounts                   # Create new currency account
GET /accounts/{id}              # Get specific account details
PUT /accounts/{id}              # Update account settings
POST /accounts/transfer         # Transfer between own accounts
POST /accounts/convert         # Currency conversion
```

**Multi-Account Data Structure:**
```go
type User struct {
    ID        int       `json:"id"`
    UserId    string    `json:"user_id"`
    Email     string    `json:"email"`
    Accounts  []Account `json:"accounts" gorm:"foreignKey:UserId"`
    CreatedAt time.Time `json:"created_at"`
}

type Account struct {
    ID        int             `json:"id"`
    AccountID string          `json:"account_id"`
    UserId    string          `json:"user_id"`
    Currency  string          `json:"currency"`
    Balance   decimal.Decimal `json:"balance"`
    IsDefault bool            `json:"is_default"`
    CreatedAt time.Time       `json:"created_at"`
}
```

### Advanced Transaction Processing Impmentation 

#### Transaction Features
**Payment Processing Enhancements:**
- **Batch Processing**: Process multiple payments in a single transaction
- **Scheduled Payments**: Support for recurring and future-dated payments
- **Payment Routing**: Intelligent routing for optimal processing fees
- **Fraud Detection**: Basic pattern recognition for suspicious transactions
- **Compliance**: AML/KYC integration points for regulatory compliance

**Transaction Flow:**
1. **Validation**: Account balance, currency availability, compliance checks
2. **FX Conversion** (if needed): Real-time rate application with margin
3. **Processing**: Ledger updates with atomic transactions
4. **Settlement**: External payment provider integration
5. **Notification**: Transaction status updates via webhooks/notifications

### Security & Compliance Features

#### Advanced Security
- **Multi-factor Authentication**: Support for TOTP, SMS, and hardware tokens
- **Rate Limiting**: Configurable limits per user/IP/endpoint
- **Encryption**: End-to-end encryption for sensitive data
- **Audit Logging**: Comprehensive audit trail for all financial operations
- **Session Management**: Secure session handling with automatic timeout

#### Regulatory Compliance
- **PSD2 Support**: European payment services directive compliance
- **Reporting**: Automated generation of regulatory reports
- **Sanctions Screening**: Integration with sanctions lists for transaction screening

### Foreign Exchange (FX) Rate Processing

#### Current Implementation
The system supports multi-currency transactions with real-time FX rate conversion:

**FX Rate Service Features:**
- **Real-time Rate Fetching**: Integration with external FX providers for live exchange rates
- **Rate Caching**: In-memory caching with TTL to reduce API calls and improve performance
- **Multi-currency Support**: Support for major currencies (USD, EUR, GBP, JPY, etc.)
- **Rate Validation**: Validates rates against reasonable thresholds to prevent manipulation
- **Fallback Mechanisms**: Multiple FX provider support with automatic failover

**FX Rate API Endpoints:**
```http
GET /fx/rates?base=USD&target=EUR
GET /fx/rates/all?base=USD
POST /fx/rates/refresh
```

**Implementation Details:**
```go
// FX Rate Service Structure
type FXRateService struct {
    cache     map[string]FXRate
    providers []FXProvider
    ttl       time.Duration
    mutex     sync.RWMutex
}

// Rate Conversion with Precision
func (fx *FXRateService) ConvertAmount(amount decimal.Decimal, from, to string) (decimal.Decimal, error) {
    rate, err := fx.GetRate(from, to)
    if err != nil {
        return decimal.Zero, err
    }
    return amount.Mul(rate.Rate), nil
}
```

## API Endpoints

### User Management

#### Register User
```http
POST /user/api/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```

#### Login
```http
POST /user/api/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "message": "Welcome back! You have successfully logged in.",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### Payment Processing

#### Internal Payment
```http
POST /payment/api/internal_payment
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "from_account": "account-uuid-1",
  "to_account": "account-uuid-2",
  "amount": 100.50,
  "currency": "USD"
}
```

#### External Payment
```http
POST /payment/api/external_payment
Authorization: Bearer <jwt_token>
Content-Type: application/json

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

**Transaction Types:**
- `BANK_TRANSFER`: Bank account transfers
- `MOBILE_MONEY`: Mobile money transfers

## 🔒 Security Features

- **JWT Authentication**: 30-minute token expiration
- **Password Hashing**: Bcrypt with salt
- **CORS Support**: Configurable cross-origin requests
- **Input Validation**: Request body validation and sanitization
- **Transaction Safety**: Database transactions for financial operations
- **Row Locking**: Prevents race conditions in balance updates

## 🧪 Testing

Run tests with:
```bash
go test ./...
```

## 📊 Database Schema

The application uses PostgreSQL with the following main tables:
- `users`: User accounts and authentication
- `accounts`: Financial accounts with balances
- `payments`: Payment transactions and history

Automatic migrations are handled on application startup.

## 🚀 Deployment

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_USER` | Database username | - |
| `DB_PASSWORD` | Database password | - |
| `DB_NAME` | Database name | grey_db |
| `DB_HOST` | Database host | localhost |
| `DB_PORT` | Database port | 5432 |
| `SERVER_PORT` | API server port | 8000 |
| `SECRET_JWT` | JWT signing secret | - |

### Docker Compose

The included `docker-compose.yml` provides:
- Go application container
- PostgreSQL database container
- Network configuration
- Volume persistence
- Health checks

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📝 License

This project is licensed under the MIT License.

## 🆘 Support

For support and questions, please open an issue in the GitHub repository.

## 🎯 Design Decisions & Trade-offs

### Architecture Decisions

#### Monolithic Structure
**Decision**: Single service architecture over microservices
**Rationale**:
- Simpler deployment and development for MVP
- Lower operational overhead
- Easier debugging and monitoring
- Faster development cycle
- Reduced network latency between components

**Trade-offs**:
- Scalability limitations as system grows
- Technology lock-in for entire application
- Harder to scale individual components
- Larger blast radius for failures

#### JWT Authentication
**Decision**: JWT-based stateless authentication
**Rationale**:
- Stateless design simplifies scaling
- No server-side session storage required
- Easy integration with mobile clients
- Standardized and widely adopted
- Good performance for API-heavy applications

**Trade-offs**:
- Token revocation complexity
- Larger request headers
- Security concerns if tokens are compromised
- No automatic session expiration on server

#### REST API Design
**Decision**: RESTful API over GraphQL or gRPC
**Rationale**:
- Standard and well-understood
- Easy to test and debug
- Good tooling support
- Browser-friendly
- Simple caching mechanisms

**Trade-offs**:
- Multiple round trips for complex data
- Over-fetching and under-fetching issues
- Versioning complexity
- Less flexible than GraphQL for client needs

### Security Decisions

#### Bcrypt Password Hashing
**Decision**: Bcrypt for password security
**Rationale**:
- Proven security track record
- Built-in salt generation
- Configurable work factor
- Resistant to rainbow table attacks
- Industry standard for password hashing

**Trade-offs**:
- Computationally intensive (by design)
- Slower than simpler hashing algorithms
- Requires careful work factor selection

#### Database Transaction Safety
**Decision**: Explicit transaction management for financial operations
**Rationale**:
- Critical for financial data integrity
- Prevents partial updates and race conditions
- ACID compliance requirements
- Rollback capability for error scenarios

**Trade-offs**:
- Increased code complexity
- Potential for deadlocks under high concurrency
- Performance overhead for transaction management
- Requires careful error handling

### Performance Optimizations

#### Decimal Precision for Financial Calculations
**Decision**: Used decimal.Decimal instead of float64
**Rationale**:
- Eliminates floating-point precision errors
- Exact decimal arithmetic required for financial calculations
- Predictable behavior for monetary operations
- Compliance with financial industry standards

**Trade-offs**:
- Performance overhead compared to float64
- More complex arithmetic operations
- Requires additional library dependency
- Memory usage higher than primitive types

#### Row-Level Locking
**Decision**: Database-level locking for balance updates
**Rationale**:
- Prevents race conditions in concurrent transactions
- Ensures data consistency under load
- Simple implementation at database level
- Reliable across different application instances

**Trade-offs**:
- Potential for deadlocks
- Reduced concurrency under high load
- Database performance impact
- Complex error handling required

### Development Trade-offs

#### Development Speed vs. Production Readiness
**Decision**: Prioritized core functionality over comprehensive features
**Rationale**:
- Faster time-to-market for MVP
- Focus on essential financial operations
- Simplified initial deployment
- Easier to validate core business logic

**Trade-offs**:
- Missing production-grade features (rate limiting, monitoring)
- Limited error handling and logging
- No automated testing framework
- Manual deployment processes

#### Code Simplicity vs. Extensibility
**Decision**: Simple, direct implementations over complex abstractions
**Rationale**:
- Easier to understand and maintain
- Faster development cycle
- Reduced cognitive overhead
- Clearer business logic implementation

**Trade-offs**:
- Harder to extend functionality
- Code duplication in some areas
- Limited reusability across components
- More refactoring needed for future features

### Future Considerations

#### Scalability Path
The current architecture supports:
- Vertical scaling (more CPU/memory)
- Database read replicas
- Connection pooling
- Caching layer addition

Future architectural changes may include:
- Microservice decomposition
- Event-driven architecture
- Message queue integration
- Distributed caching

#### Security Enhancements
Planned improvements:
- Rate limiting implementation
- API key management
- OAuth 2.0 integration
- Multi-factor authentication

These design decisions reflect the current stage of the project as a production-ready MVP with clear paths for evolution as requirements grow and scale increases.

---

**Built with ❤️ using Go, Gin, and PostgreSQL**
