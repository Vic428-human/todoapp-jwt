package handlers

import (
	"net/http"
	"strings"
	"time"
	"todo_api/internal/config"
	"todo_api/internal/models"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var registerRequest RegisterRequest

		// BindJSON 的「標準」用法
		if err := c.BindJSON(&registerRequest); err != nil {
			// 只需要 return
			// 因為 BindJSON 內部已經調用了 AbortWithError(400)
			// 客戶端會收到 400 狀態碼，但 Response Body 通常是空的或由全局錯誤處理器決定
			return
		}

		if len(registerRequest.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 6 characters long"})
			return
		}
		// 把加鹽密碼存在db，但回傳的時候不要把密碼回傳給USER
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		user := &models.User{Email: registerRequest.Email, Password: string(hashedPassword)}

		createdUser, err := repository.CreateUser(pool, user)

		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server ==> " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, createdUser)
	}
}

func LoginHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest LoginRequest

		if err := c.BindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		user, err := repository.GetUserByEmail(pool, loginRequest.Email)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		// 把存在db的加鹽密碼跟前端傳來的密碼比對
		// 例如，db中的密碼是 $2a$10$uyx7Xo1MTx1OYiBC4Gs.qO5NdWt2Wt55bGVDz.oOSHPKij23vf.Ni 這是加鹽後的密碼
		// uyx7Xo1MTx1OYiBC4Gs 就是加鹽本身，當使用者登入的時候，會拿 uyx7Xo1MTx1OYiBC4Gs 這段+使用者輸入的密碼進行加鹽，若加鹽後跟db的加鹽後一致，則等於密碼相同
		// Authorization: 是賦予權限， Authentication是進行權限驗證
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))

		if err != nil {
			// 密碼錯誤代表沒成功被賦予權限，所以失敗
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		// map[string]interface{}{}
		// map[string]any{}
		claims := jwt.MapClaims{}
		claims["user_id"] = user.ID
		claims["email"] = user.Email
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	}
}
