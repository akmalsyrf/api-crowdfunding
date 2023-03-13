package handler

import (
	"api-crowdfunding/auth"
	"api-crowdfunding/helper"
	"api-crowdfunding/user"
	"fmt"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
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

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login failed", 500, "failed", errorMessage)
		c.JSON(500, response)
		return
	}

	formatter := user.FormatUser(newUser, token)

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

	token, err := h.authService.GenerateToken(loggedInUser.ID)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login failed", 500, "failed", errorMessage)
		c.JSON(500, response)
		return
	}

	formatter := user.FormatUser(loggedInUser, token)

	response := helper.APIResponse("Successfully logged in", 200, "success", formatter)
	c.JSON(200, response)
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	var input user.CheckEmailInput
	errInput := c.ShouldBindJSON(&input)

	if errInput != nil {
		errorMessage := gin.H{"errors": errInput.Error()}

		response := helper.APIResponse("Login failed", 422, "failed", errorMessage)
		c.JSON(422, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)

	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}
		response := helper.APIResponse("Email checking failed", 500, "failed", errorMessage)
		c.JSON(500, response)
		return
	}

	var metaMessage string
	if isEmailAvailable {
		metaMessage = "Email is available"
	} else {
		metaMessage = "Email has been registered"
	}

	data := gin.H{"isAvailable": isEmailAvailable}

	response := helper.APIResponse(metaMessage, 200, "success", data)
	c.JSON(200, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")

	if err != nil {
		errorMessage := gin.H{"is_uploaded": false, "error": err}
		response := helper.APIResponse("Failed to upload avatar image", 400, "failed", errorMessage)
		c.JSON(400, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	path := "images/"

	err = helper.ValidateFolderExist(path)
	if err != nil {
		errorMessage := gin.H{"is_uploaded": false, "error": err}
		response := helper.APIResponse("Failed to upload avatar image", 500, "failed", errorMessage)
		c.JSON(400, response)
		return
	}

	path = fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		errorMessage := gin.H{"is_uploaded": false, "error": err}
		response := helper.APIResponse("Failed to upload avatar image", 400, "failed", errorMessage)
		c.JSON(400, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		errorMessage := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", 400, "failed", errorMessage)
		c.JSON(400, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Success upload avatar image", 200, "success", data)
	c.JSON(200, response)
}
