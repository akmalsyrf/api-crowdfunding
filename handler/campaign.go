package handler

import (
	"api-crowdfunding/campaign"
	"api-crowdfunding/helper"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
}

func NewCampaignHandler(campaignService campaign.Service) *campaignHandler {
	return &campaignHandler{campaignService}
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
		response := helper.APIResponse("Failed to get detail campaign", 400, "error", nil)
		c.JSON(400, response)
		return
	}

	campaignDetail, err := h.campaignService.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail campaign", 400, "error", nil)
		c.JSON(400, response)
		return
	}

	response := helper.APIResponse("Success get detail campaign", 200, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(200, response)
}
