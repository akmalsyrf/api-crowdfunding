package transaction

import (
	"api-crowdfunding/campaign"
	"api-crowdfunding/user"
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ID         int `gorm:"primary_key"`
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	PaymentUrl string
	Campaign   campaign.Campaign
	CreatedAt  time.Time
	UpdatedAt  time.Time
	User       user.User
}
