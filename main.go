package main

import (
	_config "api-crowdfunding/config"
	"api-crowdfunding/handler"
	"api-crowdfunding/service/auth"
	"api-crowdfunding/service/campaign"
	"api-crowdfunding/service/payment"
	"api-crowdfunding/service/transaction"
	"api-crowdfunding/service/user"
	"api-crowdfunding/utils/database"
	"api-crowdfunding/utils/logger"
	"fmt"
	"log"
	"time"

	cors "github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func getIndex(c *gin.Context) {
	c.Data(200, "text/plain; charset=utf-8", []byte("Hello world"))
}

func main() {
	logger.Info.Println("========= Server Setup =========")
	config, err := _config.LoadEnv()
	if err != nil {
		log.Fatal(err.Error())
	}
	db, err := database.Setup(config.Database)
	if err != nil {
		log.Fatal(err.Error())
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
	campaignHandler := handler.NewCampaignHandler(campaignService, userService, authService)
	transactionHandler := handler.NewTransactionHandler(transactionService, userService, authService)

	router := gin.Default()
	router.Static("/images", "./images")
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	api := router.Group("/api/v1")

	router.GET("/", getIndex)
	userHandler.Router(api.Group("/user"))
	campaignHandler.Router(api.Group(("/campaign")))
	transactionHandler.Router(api.Group(("/transaction")))

	logger.Info.Println("========= Server Started =========")
	port := fmt.Sprintf(":%v", config.Service.Port)
	err = router.Run(port)
	if err != nil {
		logger.Info.Println(err)
		logger.Info.Println("========= Server Ended =========")
	}
}
