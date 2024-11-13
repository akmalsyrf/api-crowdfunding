package handler

import (
	"api-crowdfunding/middleware"
	"api-crowdfunding/service/auth"
	"api-crowdfunding/service/campaign"
	"api-crowdfunding/service/user"
	"api-crowdfunding/utils/helper"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
	userService     user.Service
	authService     auth.Service
}

func NewCampaignHandler(campaignService campaign.Service, userService user.Service, authService auth.Service) *campaignHandler {
	return &campaignHandler{campaignService, userService, authService}
}

func (h *campaignHandler) Router(router *gin.RouterGroup) {
	router.GET("/", h.GetCampaigns)
	router.GET("/:id", h.GetCampaign)
	router.POST("/", middleware.AuthMiddleware(h.authService, h.userService), h.CreateCampaign)
	router.PUT("/:id", middleware.AuthMiddleware(h.authService, h.userService), h.UpdateCampaign)
	router.POST("/image", middleware.AuthMiddleware(h.authService, h.userService), h.UploadImage)
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.campaignService.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns", 400, "error", nil)
		c.JSON(400, response)
		return
	}

	formatter := campaign.FormatCampaigns(campaigns)
	response := helper.APIResponse("List of campaigns", 200, "success", formatter)
	c.JSON(200, response)
}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	errInput := c.ShouldBindUri(&input)
	if errInput != nil {
		errors := helper.FormatValidationError(errInput)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to get detail campaign", 400, "error", errorMessage)
		c.JSON(400, response)
		return
	}

	campaignDetail, err := h.campaignService.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail campaign", 400, "error", err.Error())
		c.JSON(400, response)
		return
	}

	response := helper.APIResponse("Success get detail campaign", 200, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(200, response)
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	var input campaign.CreateCampaignInput
	input.User = currentUser

	errInput := c.ShouldBindJSON(&input)
	if errInput != nil {
		errors := helper.FormatValidationError(errInput)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to create campaign", 422, "error", errorMessage)
		c.JSON(422, response)
		return
	}

	newCampaign, err := h.campaignService.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Failed to create campaign", 400, "error", nil)
		c.JSON(400, response)
		return
	}

	response := helper.APIResponse("Success create campaign", 200, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(200, response)
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput
	errInputID := c.ShouldBindUri(&inputID)
	if errInputID != nil {
		response := helper.APIResponse("Failed to update campaign", 400, "error", nil)
		c.JSON(400, response)
		return
	}

	var inputData campaign.CreateCampaignInput
	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	errInput := c.ShouldBindJSON(&inputData)
	if errInput != nil {
		errors := helper.FormatValidationError(errInput)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to update campaign", 422, "error", errorMessage)
		c.JSON(422, response)
		return
	}

	updatedCampaign, err := h.campaignService.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", 400, "error", err.Error())
		c.JSON(400, response)
		return
	}

	response := helper.APIResponse("Success update campaign", 200, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(200, response)
}

func (h *campaignHandler) UploadImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput

	errInput := c.ShouldBind(&input)
	if errInput != nil {
		errors := helper.FormatValidationError(errInput)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to create campaign", 422, "error", errorMessage)
		c.JSON(422, response)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		errorMessage := gin.H{"is_uploaded": false, "error": err.Error()}
		response := helper.APIResponse("Failed to upload campaign image", 400, "failed", errorMessage)
		c.JSON(400, response)
		return
	}
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	path := "images/"

	err = helper.ValidateFolderExist(path)
	if err != nil {
		errorMessage := gin.H{"is_uploaded": false, "error": err.Error()}
		response := helper.APIResponse("Failed to upload campaign image", 500, "failed", errorMessage)
		c.JSON(400, response)
		return
	}

	path = fmt.Sprintf("images/%d-%s", input.User.ID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		errorMessage := gin.H{"is_uploaded": false, "error": err.Error()}
		response := helper.APIResponse("Failed to upload campaign image", 400, "failed", errorMessage)
		c.JSON(400, response)
		return
	}

	_, err = h.campaignService.SaveCampaignImage(input, path)
	if err != nil {
		errFile := helper.DeleteFile(path)
		if errFile != nil {
			errorMessage := gin.H{"is_uploaded": false, "error": errFile.Error()}
			response := helper.APIResponse("Failed to upload campaign image", 400, "failed", errorMessage)
			c.JSON(400, response)
			return
		}

		errorMessage := gin.H{"is_uploaded": false, "error": err.Error()}
		response := helper.APIResponse("Failed to upload campaign image", 400, "failed", errorMessage)
		c.JSON(400, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Success to upload campaign image", 200, "success", data)
	c.JSON(200, response)
}
