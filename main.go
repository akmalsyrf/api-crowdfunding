package main

import (
	"api-crowdfunding/auth"
	"api-crowdfunding/campaign"
	"api-crowdfunding/handler"
	"api-crowdfunding/middleware"
	"api-crowdfunding/payment"
	"api-crowdfunding/transaction"
	"api-crowdfunding/user"
	"fmt"
	"log"
	"os"

	cors "github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getIndex(c *gin.Context) {
	c.Data(200, "text/plain; charset=utf-8", []byte("Hello world"))
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}
	// dsn := "root:@tcp(127.0.0.1:5432)/bwa_startup?charset=utf8mb4&parseTime=True&loc=Local"
	DB_HOST := os.Getenv("DB_HOST")
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	DB_NAME := os.Getenv("DB_NAME")
	DB_PORT := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", DB_HOST, DB_USER, DB_PASS, DB_NAME, DB_PORT)
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
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	router := gin.Default()
	router.Static("/images", "./images")
	router.Use(cors.Default())
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
	api.GET("/transactions", middleware.AuthMiddleware(authService, userService), transactionHandler.GetUserTransaction)
	api.POST("/transaction", middleware.AuthMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", middleware.AuthMiddleware(authService, userService), transactionHandler.GetNotification)

	router.Run(":8080")
}
