package database

import (
	"api-crowdfunding/config"
	"api-crowdfunding/service/campaign"
	"api-crowdfunding/service/transaction"
	"api-crowdfunding/service/user"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Setup(config config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", config.DbHost, config.DbUser, config.DbPass, config.DbName, config.DbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to postgres")
	if err := db.AutoMigrate(&user.User{}, &campaign.Campaign{},
		&campaign.CampaignImage{}, &transaction.Transaction{}); err != nil {
		log.Fatalln(err)
	}
	return db, nil
}
