package repository

import (
	"context"
	"fmt"
	"time"
	"todo_api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TODO: 交易所API
// 建立物件 → 寫入資料庫 → 回傳完整物件
func CreateProduct(pool *pgxpool.Pool, ownerId string, title string, game string, platform string, username string, views int, monthlyViews int, price int, description string, verified bool, country string, featured bool) (*models.Product, error) {
	// 建立帶有背景上下文的連接池
	var ctx context.Context
	var cancel context.CancelFunc
	// 帶有 5 秒 timeout 的 context，避免查詢卡住。
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // 釋放記憶體

	// 建立新一筆的產品
	var product models.Product

	// 在資料表名稱 products 中，對 表 的欄位新增一筆資料
	query := `
		INSERT INTO products (owner_id, title, game, platform, username, views, monthly_views, price, description, verified, country, featured)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, owner_id, title, game, platform, username, views, monthly_views, price, description, verified, country, featured, created_at, updated_at
	`

	// 其實是在做「執行 SQL（只拿一筆結果）→ 把回傳欄位塞進 todo 這個 struct」
	// ownerId, title, game ...等, 會依序對應到 SQL 裡的欄位，也就是 VALUES ($1, $2)
	// => 所以前端傳來的 ownerId, title, game ...等, 會依序對應到 VALUES ($1, $2)
	err := pool.QueryRow(ctx, query, ownerId, title, game, platform, username, views, monthlyViews, price, description, verified, country, featured).Scan(
		&product.ID,
		&product.OwnerID,
		&product.Title,
		&product.Game,
		&product.Platform,
		&product.Username,
		&product.Views,
		&product.MonthlyViews,
		&product.Price,
		&product.Description,
		&product.Verified,
		&product.Country,
		&product.Featured,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to insert product: %w", err)
	}

	return &product, nil
}
