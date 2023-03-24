package main

import (
	"api-crowdfunding/auth"
	"api-crowdfunding/campaign"
	"api-crowdfunding/handler"
	"api-crowdfunding/middleware"
	"api-crowdfunding/transaction"
	"api-crowdfunding/user"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getIndex(c *gin.Context) {
	c.Data(200, "text/plain; charset=utf-8", []byte("Hello world"))
}

func main() {
	// dsn := "root:@tcp(127.0.0.1:5432)/bwa_startup?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "host=localhost user=postgres password=password dbname=crowdfunding port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Connected to postgres")

	if err := db.AutoMigrate(&user.User{}, &campaign.Campaign{},
		&campaign.CampaignImage{}, &transaction.Transaction{}); err != nil {
		log.Fatalln(err)
	}

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	userService := user.NewService(userRepository)
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository)
	transactionService := transaction.NewService(transactionRepository, campaignRepository)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	router := gin.Default()
	router.Static("/images", "./images")
	api := router.Group("/api/v1")

	router.GET("/", getIndex)
	api.POST("/user", userHandler.RegisterUser)
	api.POST("/session", userHandler.Login)
	api.POST("/email-check", userHandler.CheckEmailAvailability)
	api.POST("/avatar", middleware.AuthMiddleware(authService, userService), userHandler.UploadAvatar)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaign/:id", campaignHandler.GetCampaign)
	api.POST("/campaign", middleware.AuthMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaign/:id", middleware.AuthMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-image", middleware.AuthMiddleware(authService, userService), campaignHandler.UploadImage)

	api.GET("/campaign/:id/transactions", middleware.AuthMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)

	router.Run(":8080")
}
