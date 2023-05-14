package handler

import (
	"api-crowdfunding/helper"
	"api-crowdfunding/transaction"
	"api-crowdfunding/user"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(transactionService transaction.Service) *transactionHandler {
	return &transactionHandler{transactionService}
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
