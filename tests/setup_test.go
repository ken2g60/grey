package tests

import (
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grey/controllers"
	"github.com/grey/database"
	"github.com/grey/middlewares"
	"github.com/grey/models"
	"github.com/grey/routers"
	"github.com/grey/utils"
	"github.com/shopspring/decimal"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestEnvironment initializes the test environment
func SetupTestEnvironment(t *testing.T) *gorm.DB {
	// Use in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database connection: %v", err)
	}

	// Run migrations
	migrations := database.Migrations{
		DB: db,
		Models: []interface{}{
			&models.User{},
			&models.Account{},
			&models.Payment{},
			&models.LedgerEntry{},
		},
	}
	database.RunMigrations(migrations)

	return db
}

// CreateTestUser creates a test user for testing
func CreateTestUser(t *testing.T, db *gorm.DB) *models.User {
	user := &models.User{
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	err := db.Create(user).Error
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return user
}

// CreateTestAccount creates a test account for testing
func CreateTestAccount(t *testing.T, db *gorm.DB, userID int, balance float64) *models.Account {
	account := &models.Account{
		UserID:   userID,
		Currency: "USD",
		Balance:  decimal.NewFromFloat(balance),
	}

	err := db.Create(account).Error
	if err != nil {
		t.Fatalf("Failed to create test account: %v", err)
	}

	return account
}

// CreateTestJWT creates a mock JWT token for testing
func CreateTestJWT(t *testing.T, email string) string {
	token, err := utils.GenerateToken("1", email)
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}
	return token
}

// SetupTestRouter creates a test router with middleware using the test database
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return routers.NewRouter()
}

// SetupTestRouterWithDB creates a test router with middleware using a specific database
func SetupTestRouterWithDB(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// Initialize repositories with test database
	userRepo := &controllers.UserGroup{}
	paymentRepo := &controllers.PaymentGroup{}

	// Override the global database for tests
	database.Db = db

	// Stripe API endpoints
	paymentGroup := router.Group("/payment/api")
	{
		paymentGroup.POST("/internal_payment", middlewares.SessionMiddleware(), paymentRepo.InternalPayment)
		paymentGroup.POST("/external_payment", middlewares.SessionMiddleware(), paymentRepo.ExternalPayment)
		paymentGroup.POST("/topup", middlewares.SessionMiddleware(), paymentRepo.TopUp)
	}

	userGroup := router.Group("/user/api")
	{
		userGroup.POST("/register", userRepo.CreateUser)
		userGroup.POST("/login", userRepo.Login)
		userGroup.GET("/profile", middlewares.SessionMiddleware(), userRepo.UserProfile)
	}

	return router
}
