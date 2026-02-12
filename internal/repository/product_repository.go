package repository

import (
	"context"
	"fmt"
	"strings"
	"time"
	"todo_api/internal/models" // 或 "github.com/gin-gonic/gin" 的 logger

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

/* 獲得所有刊登商品，但熱推優先顯示 */
func GetAllProducts(pool *pgxpool.Pool) ([]models.Product, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var products []models.Product
	// featured 熱推商品
	// true (1) 排在 false (0) 前面
	// 相同 featured 的產品，按建立時間由新至舊排序
	// 效果 : featured 為ture 會排在最前面，即使刊登時間不是最新，也會排在最前面，然後再從熱推當中時間最新的在最前面
	// 第二順位才是 featured是 false (非熱推) 進行排序，但一樣是 非熱推中最新的擺最前面
	query := `
		SELECT id, owner_id, title, game, platform, username, views, monthly_views, price, description, verified, country, featured, created_at, updated_at
		FROM products
		ORDER BY featured DESC, created_at DESC
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		if err := rows.Scan(
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
		); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

// 需要知道特定id才查得到商品
func GetProductById(pool *pgxpool.Pool, id int) (*models.Product, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, owner_id, title, game, platform, username, views, monthly_views, price, description, verified, country, featured, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var product models.Product

	err := pool.QueryRow(ctx, query, id).Scan(
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

// 假設你要在 title 裡找含有「太陽神」三個字的商品（不分大小寫）：
func SearchProducts(pool *pgxpool.Pool, keyword string) ([]models.Product, error) {

	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT id, owner_id, title, game, platform, username, views, monthly_views,
               price, description, verified, country, featured, created_at, updated_at
        FROM products
        WHERE title ILIKE $1
    `

	pattern := "%" + keyword + "%"
	rows, err := pool.Query(ctx, query, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.ID,
			&p.OwnerID,
			&p.Title,
			&p.Game,
			&p.Platform,
			&p.Username,
			&p.Views,
			&p.MonthlyViews,
			&p.Price,
			&p.Description,
			&p.Verified,
			&p.Country,
			&p.Featured,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return products, nil
}

// 可編輯欄位白名單（排除 id, owner_id, verified, featured, timestamps）
var editableFields = []string{
	"title", "game", "platform", "username",
	"views", "monthly_views", "price", "description", "country",
}

// 更新商品內容
// 解決痛點:
// 1. 不限制前端更新那些欄位，之前的寫法是所有欄位的值都要提供給後端，現在是只提供需要更新的欄位。
// 2. 那些不該被更改的內容，像是 ownerId 這種與帳號綁定的，如前端誤傳錯誤的值，後端這邊會擋掉。
// 3. 在此處 ToUpdates 方法中，定義了那些欄位才能被更新，沒出現在這上面的都不能修改。
func UpdateProduct(pool *pgxpool.Pool, id int, updates map[string]any) (*models.Product, error) {
	var ctx context.Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 建構動態 SET 子句
	var setClauses []string
	var args []any
	argIndex := 1 // $1 給 WHERE id

	for _, field := range editableFields {
		if val, ok := updates[field]; ok && val != nil {
			setClauses = append(setClauses, fmt.Sprintf("%s=$%d", field, argIndex+1))
			args = append(args, val)
			argIndex++
		}
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no valid fields to update")
	}

	// 動態查詢（安全，參數化）
	query := fmt.Sprintf(`
        UPDATE products 
        SET %s, updated_at = NOW()
        WHERE id = $1
        RETURNING id, owner_id, title, game, platform, username, views, monthly_views, 
                  price, description, verified, country, featured, created_at, updated_at
    `, strings.Join(setClauses, ", "))

	args = append([]interface{}{id}, args...) // id 第一個

	var updatedProduct models.Product
	err := pool.QueryRow(ctx, query, args...).Scan(
		&updatedProduct.ID, &updatedProduct.OwnerID, &updatedProduct.Title,
		&updatedProduct.Game, &updatedProduct.Platform, &updatedProduct.Username,
		&updatedProduct.Views, &updatedProduct.MonthlyViews, &updatedProduct.Price,
		&updatedProduct.Description, &updatedProduct.Verified, &updatedProduct.Country,
		&updatedProduct.Featured, &updatedProduct.CreatedAt, &updatedProduct.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return &updatedProduct, nil
}
