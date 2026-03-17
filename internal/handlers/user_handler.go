package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"todo_api/internal/config"
	"todo_api/internal/models"
	"todo_api/internal/repository"
	"todo_api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

		user := &models.User{
			Email:    registerRequest.Email,
			Password: string(hashedPassword),
		}

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
			return
		}

		user, err := repository.GetUserByEmail(pool, loginRequest.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		// 把存在db的加鹽密碼跟前端傳來的密碼比對
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
		if err != nil {
			// 密碼錯誤代表沒成功被賦予權限，所以失敗
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		// creating, signing, and encoding a JWT token using the HMAC signing method
		t := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"user_id": user.ID,
				"email":   user.Email,
				"exp":     time.Now().Add(24 * time.Hour).Unix(), // Unix() 代表 UTC 秒數時間戳
			})

		tokenString, err := t.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
	}
}

func SetProfileImageHandler(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 從 URL 取得 user id
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "user id is required",
			})
			return
		}

		// 2. 從 multipart/form-data 取得圖片檔案
		// Postman / 前端欄位名稱要用 image
		imageFileHeader, err := c.FormFile("image")
		if err != nil {
			log.Printf("failed to get image file from request: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "failed to get image file",
			})
			return
		}

		log.Printf(
			"received upload file: userID=%s filename=%s size=%d\n",
			id,
			imageFileHeader.Filename,
			imageFileHeader.Size,
		)

		// 3. 呼叫 service 處理完整流程
		updatedUser, err := userService.SetProfileImage(
			c.Request.Context(),
			id,
			imageFileHeader,
		)
		if err != nil {
			log.Printf("failed to set profile image: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to update profile image",
			})
			return
		}

		// 4. 回傳更新後的 user
		c.JSON(http.StatusOK, gin.H{
			"message": "profile image updated successfully",
			"user":    updatedUser,
		})
	}
}

func TestProtectedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// "Comma Ok" 慣例，從 map 取值
		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "protected route accessed successfully",
			"user_id": userID,
		})
	}
}
