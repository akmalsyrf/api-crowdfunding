package handler

import (
	"api-crowdfunding/middleware"
	"api-crowdfunding/service/auth"
	"api-crowdfunding/service/transaction"
	"api-crowdfunding/service/user"
	"api-crowdfunding/utils/helper"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service     transaction.Service
	userService user.Service
	authService auth.Service
}

func NewTransactionHandler(transactionService transaction.Service, userService user.Service, authService auth.Service) *transactionHandler {
	return &transactionHandler{transactionService, userService, authService}
}

func (h *transactionHandler) Router(router *gin.RouterGroup) {
	router.GET("/:id/campaign", middleware.AuthMiddleware(h.authService, h.userService), h.GetCampaignTransactions)
	router.GET("/", middleware.AuthMiddleware(h.authService, h.userService), h.GetUserTransaction)
	router.POST("/", middleware.AuthMiddleware(h.authService, h.userService), h.CreateTransaction)
	router.POST("/notification", middleware.AuthMiddleware(h.authService, h.userService), h.GetNotification)
}

func (h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input transaction.GetTransactionByCampaignInput

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	errInput := c.ShouldBindUri(&input)
	if errInput != nil {
		errors := helper.FormatValidationError(errInput)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to get campaign's transaction", 400, "error", errorMessage)
		c.JSON(400, response)
		return
	}

	transactions, err := h.service.GetTransactionByCampaignID(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Failed to get campaign's transaction", 400, "error", errorMessage)
		c.JSON(400, response)
		return
	}

	response := helper.APIResponse("Success get campaign transaction", 200, "success", transaction.FormatCampaignTransactions(transactions))
	c.JSON(200, response)
}

func (h *transactionHandler) GetUserTransaction(c *gin.Context) {

	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	transactions, err := h.service.GetTransactionByUserID(userID)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Failed to get campaign's transaction", 400, "error", errorMessage)
		c.JSON(400, response)
		return
	}

	formatter := transaction.FormatUserTransactions(transactions)
	response := helper.APIResponse("Success get campaign transaction", 200, "success", formatter)
	c.JSON(200, response)
}

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	var input transaction.CreateTransactionInput

	errInput := c.ShouldBindJSON(&input)
	if errInput != nil {
		errors := helper.FormatValidationError(errInput)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to create transaction", 422, "error", errorMessage)
		c.JSON(422, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newTransaction, err := h.service.CreateTransaction(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Failed create transaction", 400, "error", errorMessage)
		c.JSON(400, response)
		return
	}

	formatter := transaction.FormatTransaction(newTransaction)
	response := helper.APIResponse("Success create transaction", 200, "success", formatter)
	c.JSON(200, response)
}

func (h *transactionHandler) GetNotification(c *gin.Context) {
	var input transaction.TransactionNotificationInput

	errInput := c.ShouldBindJSON(&input)
	if errInput != nil {
		errors := helper.FormatValidationError(errInput)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to get notification", 422, "error", errorMessage)
		c.JSON(422, response)
		return
	}

	err := h.service.ProcessPayment(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Failed to get notification", 400, "error", errorMessage)
		c.JSON(400, response)
		return
	}

	c.JSON(200, input)
}
