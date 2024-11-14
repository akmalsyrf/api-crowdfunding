package handler

import (
	"api-crowdfunding/middleware"
	"api-crowdfunding/service/auth"
	"api-crowdfunding/service/user"
	"api-crowdfunding/utils/helper"
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

func (h *userHandler) Router(router *gin.RouterGroup) {
	router.POST("/", h.RegisterUser)
	router.POST("/session", h.Login)
	router.POST("/email-check", h.CheckEmailAvailability)
	router.POST("/avatar", middleware.AuthMiddleware(h.authService, h.userService), h.UploadAvatar)
	router.GET("/fetch", middleware.AuthMiddleware(h.authService, h.userService), h.FetchUser)
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

	isNotExist, err := h.userService.IsEmailAvailable(user.CheckEmailInput{Email: input.Email})
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login failed", 500, "failed", errorMessage)
		c.JSON(500, response)
		return
	}
	if !isNotExist {
		response := helper.APIResponse("Email has been registered", 409, "failed", nil)
		c.JSON(409, response)
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
		response := helper.APIResponse("Login failed", 400, "failed", errorMessage)
		c.JSON(400, response)
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

func (h *userHandler) FetchUser(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)

	formatter := user.FormatUser(currentUser, "")
	response := helper.APIResponse("Success fetch current user", 200, "success", formatter)
	c.JSON(200, response)
}
