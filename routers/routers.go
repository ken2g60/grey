package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grey/controllers"
	"github.com/grey/middlewares"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// Initialize repositories
	userRepo := &controllers.UserGroup{}
	paymentRepo := &controllers.PaymentGroup{}

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
