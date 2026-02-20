package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grey/database"
	"github.com/grey/middlewares"
	"github.com/grey/models"
	"github.com/grey/service"
	"github.com/grey/structs"
	"github.com/grey/utils"
	"github.com/shopspring/decimal"
)

type PaymentGroup struct{}

func (repository *PaymentGroup) InternalPayment(c *gin.Context) {

	claimPayload, exists := c.Get("x-claim-payload")
	_, _ = claimPayload.(middlewares.JwtSessionPayload)
	if !exists {
		c.AbortWithStatus(401)
		return
	}

	var form structs.InternalPaymentRequest
	err := c.ShouldBindJSON(&form)
	if err != nil {
		utils.ErrorResponse(c, "Please check your payment details and ensure all required fields are filled correctly.")
		return
	}

	fromID, err := models.IsAccountExists(c.Request.Context(), database.Db, form.FromAccount)
	if err != nil {
		utils.ErrorResponse(c, "Account does not exist")
		return
	}

	proceesedAmount := decimal.NewFromFloat(form.Amount)
	if fromID.Balance.Cmp(proceesedAmount) < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Insufficient balance",
			"status":  http.StatusBadRequest,
		})
		return
	}

	toID, err := models.IsAccountExists(c.Request.Context(), database.Db, form.ToAccount)
	if err != nil {
		utils.ErrorResponse(c, "Account does not exist")
		return
	}

	processedAmount := decimal.NewFromFloat(form.Amount)
	response, err := service.ProcessInternalPayment(c.Request.Context(), database.Db, fromID.AccountID, toID.AccountID, processedAmount, form.Currency)
	if err != nil {
		if err.Error() == "invalid amount" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid amount",
				"error":   err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to process payment",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(200, gin.H{
		"response": response,
		"message":  "Payment created successfully",
	})
}

func (repository *PaymentGroup) ExternalPayment(c *gin.Context) {
	claimPayload, exists := c.Get("x-claim-payload")
	_, _ = claimPayload.(middlewares.JwtSessionPayload)
	if !exists {
		c.AbortWithStatus(401)
		return
	}

	var form structs.ExternalPaymentRequest
	err := c.ShouldBindJSON(&form)
	if err != nil {
		utils.ErrorResponse(c, "Please check your payment details and ensure all required fields are filled correctly.")
		return
	}

	fromID, err := models.IsAccountExists(c.Request.Context(), database.Db, form.Account)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Account not found",
			"error":   err.Error(),
		})
		return
	}

	proceesedAmount := decimal.NewFromFloat(form.Amount)
	if fromID.Balance.Cmp(proceesedAmount) < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Insufficient balance",
			"status":  http.StatusBadRequest,
		})
		return
	}

	switch form.TransactionType {
	case "BANK_TRANSFER":
		response, err := service.ProcessExternalPayment(c.Request.Context(), database.Db, form.Recipient, fromID.AccountID, decimal.NewFromFloat(form.Amount), form.Currency)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to process payment",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"response": response,
			"message":  "Payment created successfully",
		})
	case "MOBILE_MONEY":
		response, err := service.ProcessExternalPayment(c.Request.Context(), database.Db, form.Recipient, fromID.AccountID, decimal.NewFromFloat(form.Amount), form.Currency)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to process payment",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"response": response,
			"message":  "Payment created successfully",
		})

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid transaction type support BANK_TRANSFER or MOBILE_MONEY",
		})
	}
}

func (repository *PaymentGroup) TopUp(c *gin.Context) {
	claimPayload, exists := c.Get("x-claim-payload")
	_, _ = claimPayload.(middlewares.JwtSessionPayload)
	if !exists {
		c.AbortWithStatus(401)
		return
	}

	var form structs.TopUp
	err := c.ShouldBindJSON(&form)
	if err != nil {
		utils.ErrorResponse(c, "Please check your payment details and ensure all required fields are filled correctly.")
		return
	}

	fromID, err := models.IsAccountExists(c.Request.Context(), database.Db, form.Account)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Account not found",
			"error":   err.Error(),
		})
		return
	}

	response, err := service.TopUpProcess(c.Request.Context(), database.Db, fromID.AccountID, decimal.NewFromFloat(form.Amount), form.Currency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to process payment",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"response": response,
		"message":  "topup successfully",
	})

}
