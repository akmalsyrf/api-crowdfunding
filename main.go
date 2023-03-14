package main

import (
	"api-crowdfunding/auth"
	"api-crowdfunding/campaign"
	"api-crowdfunding/handler"
	"api-crowdfunding/helper"
	"api-crowdfunding/user"
	"fmt"
	"log"
	"strings"

	"github.com/dgrijalva/jwt-go"
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

	if err := db.AutoMigrate(&user.User{}, &campaign.Campaign{}); err != nil {
		log.Fatalln(err)
	}

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)

	userService := user.NewService(userRepository)
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	router := gin.Default()
	router.Static("/images", "./images")
	api := router.Group("/api/v1")

	router.GET("/", getIndex)
	api.POST("/user", userHandler.RegisterUser)
	api.POST("/session", userHandler.Login)
	api.POST("/email-check", userHandler.CheckEmailAvailability)
	api.POST("/avatar", authMiddleware(authService, userService), userHandler.UploadAvatar)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaign/:id", campaignHandler.GetCampaign)

	router.Run(":8080")
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", 401, "error", nil)
			c.AbortWithStatusJSON(401, response)
			return
		}

		// "Bearer (token)"
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")

		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)

		if err != nil {
			response := helper.APIResponse("Unauthorized", 401, "error", nil)
			c.AbortWithStatusJSON(401, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", 401, "error", nil)
			c.AbortWithStatusJSON(401, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserById(userID)

		if err != nil {
			response := helper.APIResponse("Unauthorized", 401, "error", nil)
			c.AbortWithStatusJSON(401, response)
			return
		}

		c.Set("currentUser", user)
	}
}
