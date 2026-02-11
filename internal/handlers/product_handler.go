package handlers

import (
	"net/http"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateProductRequest struct {
	OwnerID      string `json:"ownerId" binding:"required"`
	Title        string `json:"title" binding:"required"`
	Game         string `json:"game" binding:"required"`
	Platform     string `json:"platform" binding:"required"`
	Username     string `json:"username" binding:"required"`
	Views        int    `json:"views" binding:"required"`
	MonthlyViews int    `json:"monthly_views" binding:"required"`
	Price        int    `json:"price" binding:"required"`
	Description  string `json:"description" binding:"required"`
	// TODO: 以正常管道註冊的是驗證用戶，如果使開發人員發的帳號就屬於非驗證用戶
	Verified bool   `json:"verified" binding:"required"`
	Country  string `json:"country" binding:"required"`
	// 熱播推薦
	Featured bool `json:"featured"`
}

func CreatteProductHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateProductRequest

		// 先驗證 client 傳來的資料
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// 沒問題後，把資料傳給 repository 層，透過 sql 方式把資料寫入到DB
		proudct, err := repository.CreateProduct(pool, input.OwnerID, input.Title, input.Game, input.Platform, input.Username, input.Views, input.MonthlyViews, input.Price, input.Description, input.Verified, input.Country, input.Featured)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		// 資料庫寫入正確後，回傳訊息到 client 端
		c.JSON(http.StatusCreated, proudct)
	}
}

func GetAllProductsHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := repository.GetAllProducts(pool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, products)
	}
}
