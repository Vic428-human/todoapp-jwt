package repository

import (
	"context"
	"fmt"
	"time"
	"todo_api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// 預期是給交易所建立掛單商品用的
// 建立物件 → 寫入資料庫 → 回傳完整物件
func CreateProduct(pool *pgxpool.Pool, ownerId string, title string, game string, platform string, username string, views int, monthlyViews int, price int, description string, verified bool, country string, featured bool) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 設定建立與更新時間
	createdAt := time.Now()
	updatedAt := time.Now()

	// 建立新一筆的產品
	product := &models.Product{
		OwnerID:      ownerId,
		Title:        title,
		Game:         game,
		Platform:     platform,
		Username:     username,
		Views:        views,
		MonthlyViews: monthlyViews,
		Price:        price,
		Description:  description,
		Verified:     verified,
		Country:      country,
		Featured:     featured,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
	// 執行 SQL 插入語句，RETURNING id 以取得主鍵
	query := `
		INSERT INTO products (owner_id, title, game, platform, username, views, monthly_views, price, description, verified, country, featured, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`
	// 使用 QueryRow(...).Scan(...) 來取得 RETURNING id，確保 todo.ID 被正確填入
	err := pool.QueryRow(ctx, query,
		product.OwnerID,
		product.Title,
		product.Game,
		product.Platform,
		product.Username,
		product.Views,
		product.MonthlyViews,
		product.Price,
		product.Description,
		product.Verified,
		product.Country,
		product.Featured,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert product: %w", err)
	}

	return product, nil
}
