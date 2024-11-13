package middleware

import (
	"api-crowdfunding/service/auth"
	"api-crowdfunding/service/user"
	"api-crowdfunding/utils/helper"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
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
