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

	errInput := c.ShouldBindJSON(&input)
	if errInput != nil {
		errors := helper.FormatValidationError(errInput)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed create data", 422, "failed", errorMessage)
		c.JSON(422, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Failed create data", 500, "failed", errorMessage)
		c.JSON(500, response)
		return
	}

	formatter := user.FormatUser(newUser, "")

	response := helper.APIResponse("Account has been registered", 200, "success", formatter)
	c.JSON(200, response)

}

func (h *userHandler) Login(c *gin.Context) {
	var input user.LoginInput

	errInput := c.ShouldBindJSON(&input)

	if errInput != nil {
		errors := helper.FormatValidationError(errInput)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login failed", 422, "failed", errorMessage)
		c.JSON(422, response)
		return
	}

	loggedInUser, err := h.userService.Login(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login failed", 500, "failed", errorMessage)
		c.JSON(500, response)
		return
	}

	formatter := user.FormatUser(loggedInUser, "")

	response := helper.APIResponse("Successfully logged in", 200, "success", formatter)
	c.JSON(200, response)
}
