package main

import (
	"api-crowdfunding/handler"
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

	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Fatalln(err)
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	api := router.Group("/api/v1")

	router.GET("/", getIndex)
	api.POST("/user", userHandler.RegisterUser)
	api.POST("/session", userHandler.Login)
	api.POST("/email-check", userHandler.CheckEmailAvailability)

	router.Run(":8080")
}
