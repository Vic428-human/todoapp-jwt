// security check point
package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"todo_api/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// once we receive our request, the request is going to have something with it that's called a header
		// we need to make sure it is the same user who has logged in who can create todos
		bearer := c.GetHeader("Authorization")

		if bearer == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort() // https://pkg.go.dev/github.com/gin-gonic/gin#Context.Abort
			return
		}

		parts := strings.Split(bearer, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "empty token"})
			c.Abort()
			return
		}

		// Parse 讓你用 MapClaims 自己取值，所以這邊不是用 ParseWithClaims，所以不需要自訂 struct 直接拿欄位
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { //063951
				return nil, fmt.Errorf("不支援的簽章演算法：%v", token.Method.Alg())
			}

			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// parsing and validating a token using the HMAC signing method
		// 為什麼 JWT需要使用 Claims ?
		// 使用者基本資訊、做權限控制與授權，因為有些api操作是基於特定權限條件滿足後才能使用，例如交易所這邊是刊登商品，有權限的人才能刊登
		if claims, ok := token.Claims.(jwt.MapClaims); ok { // https://pkg.go.dev/github.com/golang-jwt/jwt/v5#section-readme
			userID := claims["user_id"].(string)

			// time.Now().Unix() => int64
			//  claims["exp"] => float64
			expFloat := claims["exp"].(float64)
			exp := int64(expFloat) // Go 不會幫你自動轉型

			if time.Now().Unix() > exp {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
				c.Abort()
				return
			}

			c.Set("user_id", userID)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		}

	}
}
