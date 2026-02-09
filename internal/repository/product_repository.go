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
	item := &models.Product{
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

	query := `
		INSERT INTO products (owner_id, title, game, platform, username, views, monthly_views, price, description, verified, country, featured, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`

	err := pool.QueryRow(ctx, query,
		item.OwnerID,
		item.Title,
		item.Game,
		item.Platform,
		item.Username,
		item.Views,
		item.MonthlyViews,
		item.Price,
		item.Description,
		item.Verified,
		item.Country,
		item.Featured,
		item.CreatedAt,
		item.UpdatedAt,
	).Scan(&item.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to insert product: %w", err)
	}

	return item, nil
}
