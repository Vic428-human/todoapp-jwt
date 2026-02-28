// security check point
package middleware

import (
	"fmt"
	"net/http"
	"strings"
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

		jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 演算法白名單檢查，我只接受 HMAC 系列的簽名方法（也就是 HS256、HS384、HS512）
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("不支援的簽章演算法：%v", token.Method.Alg())
			}

		})
	}
}
