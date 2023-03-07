package handler

import (
	"api-crowdfunding/helper"
	"api-crowdfunding/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	user, err := h.userService.RegisterUser(input)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		response := helper.APIResponse("Account has been registered", 200, "success", user)
		c.JSON(200, response)
	}
}
