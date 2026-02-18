# Architecture Documentation

## System Overview

The Grey FinTech API is a microservice-based financial system designed for secure payment processing and user management. The architecture follows clean architecture principles with clear separation of concerns.

## Architecture Patterns

### Clean Architecture
The system implements clean architecture with the following layers:

1. **Presentation Layer** (`controllers/`, `routers/`)
   - HTTP request handling
   - Request validation
   - Response formatting

2. **Business Logic Layer** (`service/`)
   - Core business rules
   - Payment processing logic
   - Transaction management

3. **Data Access Layer** (`models/`, `database/`)
   - Database operations
   - Data models
   - ORM interactions

4. **Cross-cutting Concerns** (`middlewares/`, `utils/`)
   - Authentication
   - Error handling
   - Common utilities

### Repository Pattern
The system uses a repository-like pattern where controllers act as repositories, abstracting data access logic.

## Component Architecture

### Core Components

#### 1. Application Entry Point (`main.go`)
```go
func main() {
    // Database initialization
    db := database.InitDb()
    
    // Database migrations
    database.RunMigrations(migrations)
    
    // HTTP server setup
    srv := &http.Server{
        Addr:    ":" + port,
        Handler: routers.NewRouter(),
    }
    
    // Graceful shutdown handling
}
```

**Responsibilities**:
- Application bootstrap
- Database connection setup
- Server lifecycle management
- Graceful shutdown

#### 2. Router Configuration (`routers/routers.go`)
```go
func NewRouter() *gin.Engine {
    router := gin.New()
    router.Use(gin.Logger())
    router.Use(gin.Recovery())
    router.Use(cors.New(config))
    
    // Route groups and middleware setup
}
```

**Responsibilities**:
- HTTP routing configuration
- Middleware chaining
- CORS setup
- Request/response pipeline

#### 3. Controllers Layer (`controllers/`)

**User Controller** (`controllers/users.go`):
- User registration with transaction safety
- Authentication and token generation
- Account auto-creation on registration

**Payment Controller** (`controllers/payments.go`):
- Internal payment processing
- External payment processing
- Transaction type validation

#### 4. Service Layer (`service/payment.go`)
```go
func ProcessInternalPayment(ctx context.Context, DB *gorm.DB, 
    fromAccount, toAccount string, amount decimal.Decimal, 
    currency string) (Payment models.Payment, err error)
```

**Responsibilities**:
- Business logic implementation
- Transaction management
- Balance operations with row locking
- Payment record creation

#### 5. Models Layer (`models/`)

**User Model** (`models/users.go`):
- User entity definition
- Password hashing with bcrypt
- Email uniqueness validation
- CRUD operations

**Account Model** (`models/accounts.go`):
- Account entity with decimal precision
- Balance management
- User relationship
- UUID generation

**Payment Model** (`models/payments.go`):
- Payment transaction entity
- Status management (pending, completed, failed)
- Transaction relationships

## Data Flow Architecture

### User Registration Flow
```
Client Request → Router → User Controller → Database Transaction → 
User Creation → Account Creation → Transaction Commit → Response
```

### Payment Processing Flow
```
Client Request → Auth Middleware → Payment Controller → 
Service Layer → Database Transaction → Balance Updates → 
Payment Record → Transaction Commit → Response
```

## Database Architecture

### Schema Design
The database uses PostgreSQL with the following design principles:

#### Normalization
- Third Normal Form (3NF) compliance
- Proper foreign key relationships
- Minimal data redundancy

#### Data Types
- UUID for unique identifiers
- Decimal(18,2) for monetary values
- Timestamps for audit trails
- VARCHAR(3) for currency codes

#### Indexing Strategy
- Primary keys on all tables
- Unique indexes on emails
- Indexes on UUID fields for lookups

### Transaction Management
```go
tx := DB.WithContext(ctx).Begin()
defer tx.Rollback()

// Database operations
tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE account_id = $2 AND balance >= $1", amount, fromAccount)

tx.Commit()
```

**Features**:
- ACID compliance
- Row-level locking for balance updates
- Automatic rollback on errors
- Context-aware transactions

## Security Architecture

### Authentication Flow
```
Client → JWT Token → Middleware → Token Validation → 
Claim Extraction → Context Setting → Controller
```

### Security Layers

#### 1. JWT Authentication (`middlewares/auth.go`)
- Token-based authentication
- 30-minute token expiration
- Claim-based authorization
- Bearer token format

#### 2. Password Security (`models/users.go`)
- Bcrypt hashing with salt
- Secure password comparison
- No plain text storage

#### 3. Input Validation
- Request body validation
- SQL injection prevention via ORM
- Type safety with Go's type system

#### 4. CORS Configuration
- Configurable origin policies
- Method and header restrictions
- Development-friendly defaults

## Performance Architecture

### Database Optimization
- Connection pooling via GORM
- Prepared statements
- Efficient indexing
- Row locking for concurrent operations

### Memory Management
- Context-aware operations
- Proper resource cleanup
- Minimal memory footprint

### Concurrency Handling
- Database-level locking
- Transaction isolation
- Race condition prevention

## Scalability Architecture

### Horizontal Scaling
- Stateless application design
- External database dependency
- Container-ready deployment

### Vertical Scaling
- Efficient resource usage
- Minimal memory allocation
- Optimized database queries

### Caching Strategy
Currently not implemented but planned:
- Redis for session storage
- Application-level caching
- Database query caching

## Deployment Architecture

### Containerization (`Dockerfile`)
```dockerfile
FROM golang:1.25-alpine AS builder
# Multi-stage build for optimized images
FROM alpine:latest
# Minimal runtime image
```

### Orchestration (`docker-compose.yml`)
- Multi-container setup
- Network isolation
- Volume persistence
- Health checks

### Environment Configuration
- Environment variable-based configuration
- Secret management
- Development vs production settings

## Monitoring & Observability

### Current Implementation
- Basic logging with Gin
- Error tracking
- Database connection monitoring

### Planned Enhancements
- Structured logging with Zap
- Metrics collection
- Health check endpoints
- Distributed tracing

## Error Handling Architecture

### Error Types
1. **Validation Errors**: Input validation failures
2. **Business Logic Errors**: Insufficient funds, account not found
3. **System Errors**: Database failures, network issues
4. **Authentication Errors**: Invalid tokens, expired sessions

### Error Response Strategy
```go
func ErrorResponse(c *gin.Context, message string) {
    c.JSON(http.StatusBadRequest, gin.H{
        "message": message,
        "status":  http.StatusBadRequest,
    })
}
```

### Error Propagation
- Context-aware error handling
- Transaction rollback on errors
- Graceful degradation
- User-friendly error messages

## Future Architecture Enhancements

### Microservices Decomposition
- Separate user service
- Payment service isolation
- Account service extraction
- API gateway implementation

### Event-Driven Architecture
- Message queue integration
- Event sourcing for payments
- Async processing capabilities
- Event replay functionality

### API Evolution
- GraphQL support
- API versioning strategy
- Backward compatibility
- Deprecation policies

### Security Enhancements
- OAuth 2.0 implementation
- Multi-factor authentication
- Rate limiting
- API key management

## Technology Stack Rationale

### Go Language Choice
- Performance and concurrency
- Strong typing
- Standard library richness
- Deployment simplicity

### Gin Framework
- HTTP routing efficiency
- Middleware ecosystem
- JSON handling
- Development productivity

### PostgreSQL Database
- ACID compliance
- JSON support
- Advanced indexing
- Transaction reliability

### GORM ORM
- Developer productivity
- Type safety
- Migration support
- Relationship handling

## Development Workflow

### Code Organization
- Package-based structure
- Clear separation of concerns
- Testable architecture
- Documentation standards

### Development Practices
- Clean code principles
- SOLID design principles
- Test-driven development
- Code review process

This architecture provides a solid foundation for a secure, scalable, and maintainable financial technology platform.
