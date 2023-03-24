package handler

import (
	"api-crowdfunding/helper"
	"api-crowdfunding/transaction"
	"api-crowdfunding/user"
	"fmt"

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
	fmt.Println("PARAMS ", input.ID)

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
