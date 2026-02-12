package handlers

import (
	"log"
	"net/http"
	"strconv"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

type UpdateProductRequest struct {
	Title *string `json:"title"`
	// Game  *string `json:"game"` 編輯刊登商品不會同時修改到遊戲，不然你賣太陽神頭盔明明只會出現在ro，卻被你改去天堂，就不正常
	// Platform *string `json:"platform"` // 編輯刊登商品不會同時修改到平台，因為這個是跟登入帳號時就連動的
	Username *string `json:"username"` // 綁定帳號的預期只有 ownerId ，例如一個帳號就是一個 ownerId，但一個 ownerId 可以有很多角色名稱，若特定帳號有問題，直接從 ownerId 去調資料就好
	// Views        *int    `json:"views"` 編輯刊登商品不會同時修改到觀看次數
	// MonthlyViews *int    `json:"monthly_views"` 編輯刊登商品不會同時修改到月觀看次數
	Price       *int    `json:"price"`
	Description *string `json:"description"`
	// Country     *string `json:"country"` 編輯刊登商品不會同時修改到國家，因為這個是跟登入帳號時就連動的
}

// 轉 map 的 helper
func (r *UpdateProductRequest) ToUpdates() map[string]interface{} {
	updates := make(map[string]interface{})
	if r.Title != nil {
		updates["title"] = *r.Title
	}
	// if r.Platform != nil {
	// 	updates["platform"] = *r.Platform
	// }
	if r.Username != nil {
		updates["username"] = *r.Username
	}
	if r.Price != nil {
		updates["price"] = *r.Price
	}
	if r.Description != nil {
		updates["description"] = *r.Description
	}

	return updates
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

func GetProductByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID PRODUCT ID"})
			return
		}
		product, err := repository.GetProductById(pool, id)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "product not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, product)
	}
}

// TODO: 這個是搜尋關鍵字的名稱之後再改，有點容易混淆
func ListProductsHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyword := c.Query("keyword")
		log.Printf("DEBUG: received keyword='%s'", keyword) // ← 加這行
		if keyword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "keyword required"})
			return
		}
		products, err := repository.SearchProducts(pool, keyword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, products)
	}
}

func UpdateProductHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID PRODUCT ID"})
			return
		}

		// 如果傳入已經被排除的欄位，代碼會自動忽略 ex:　http://localhost:3000/products/3
		// 你預期修改的是 id=3，但卻同時想要修改 id成4，這不被允許，所以會排除這個寫入
		var input UpdateProductRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ✅ 轉成 map，只包含有值的欄位（真正的 PATCH）
		updates := input.ToUpdates()
		if len(updates) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
			return
		}

		// ✅ 呼叫新版 repository
		updated, err := repository.UpdateProduct(pool, id, updates)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": updated})
	}
}
