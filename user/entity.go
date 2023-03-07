package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID             int `gorm:"primary_key"`
	Name           string
	Occupation     string
	Email          string
	PasswordHash   string
	AvatarFileName string
	Role           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
