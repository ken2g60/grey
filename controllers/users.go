package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grey/database"
	"github.com/grey/middlewares"
	"github.com/grey/models"
	"github.com/grey/structs"
	"github.com/grey/utils"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
)

type UserGroup struct{}

func (repository *UserGroup) CreateUser(c *gin.Context) {

	var form structs.User
	err := c.ShouldBindJSON(&form)
	if err != nil {
		utils.ErrorResponse(c, "Provide all the required fields")
		return
	}

	// check if user exists
	user, err := models.IsEmailExists(c.Request.Context(), database.Db, form.Email)
	if err != nil {
		utils.ErrorResponse(c, "We couldn't check your email at this time. Please try again later.")
		return
	}

	if user.ID != 0 {
		utils.ErrorResponse(c, "email already exists")
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), 10)
	if err != nil {
		utils.ErrorResponse(c, "We encountered an issue while securing your password. Please try again.")
		return
	}

	user = &models.User{
		Email:    form.Email,
		Password: string(hashPassword),
	}

	// transaction to create user and account
	tx := database.Db.WithContext(c.Request.Context()).Begin()

	// create user
	err = tx.Create(user).Error
	if err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, "error creating user")
		return
	}

	// create account for user
	account := &models.Account{
		UserID:   user.ID,
		Currency: "USD",
		Balance:  decimal.NewFromFloat(0.00),
	}
	err = tx.Create(account).Error
	if err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, "error creating account")
		return
	}

	err = tx.Commit().Error
	if err != nil {
		utils.ErrorResponse(c, "error committing transaction")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success",
		"data":    "user registered successfully",
		"status":  http.StatusCreated,
	})
}

func (repository *UserGroup) Login(c *gin.Context) {
	var form structs.User
	err := c.ShouldBindJSON(&form)
	if err != nil {
		utils.ErrorResponse(c, "Please check your login details and ensure all required fields are filled correctly.")
		return
	}

	user, err := models.IsEmailExists(c, database.Db, form.Email)
	if err != nil {
		utils.ErrorResponse(c, "We couldn't check your email at this time. Please try again later.")
		return
	}

	if user.Email == "" {
		utils.ErrorResponse(c, "We couldn't find an account with that email. Please check your email address or sign up for a new account.")
		return
	}

	err = models.PasswordCompare(form.Password, user.Password)
	if err != nil {
		utils.ErrorResponse(c, "Invalid email or password")
		return
	}

	token, err := utils.GenerateToken(user.UserId, user.Email)
	if err != nil {
		utils.ErrorResponse(c, "We couldn't log you in at this time. Please try again later.")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome back! You have successfully logged in.",
		"data": gin.H{
			"token": token,
		},
	})

}

func (repository *UserGroup) UserProfile(c *gin.Context) {
	claimPayload, exists := c.Get("x-claim-payload")
	JwtSessionPayload, _ := claimPayload.(middlewares.JwtSessionPayload)
	if !exists {
		c.AbortWithStatus(401)
		return
	}

	user, err := models.UserProfile(c.Request.Context(), database.Db, JwtSessionPayload.UserID)
	if err != nil {
		utils.ErrorResponse(c, "We couldn't retrieve your profile at this time. Please try again later.")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    user,
	})

}
