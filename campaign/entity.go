package campaign

import (
	"api-crowdfunding/user"
	"time"

	"gorm.io/gorm"
)

type Campaign struct {
	gorm.Model
	ID               int `gorm:"primary_key"`
	UserID           int
	Name             string
	ShortDescription string
	Description      string
	Perks            string
	BackerCount      int
	GoalAmount       int
	CurrentAmount    int
	Slug             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CampaignImages   []CampaignImage
	User             user.User
}

type CampaignImage struct {
	gorm.Model
	ID         int `gorm:"primary_key"`
	CampaignID int
	FileName   string
	IsPrimary  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
