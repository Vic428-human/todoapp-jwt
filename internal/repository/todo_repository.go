package repository

import (
	"context"
	"fmt"
	"time"
	"todo_api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// 建立物件 → 寫入資料庫 → 回傳完整物件
func CreateTodo(pool *pgxpool.Pool, title string, completed bool) (*models.Todo, error) {
	// 建立帶有背景上下文的連線池
	var ctx context.Context
	var cancel context.CancelFunc
	// 帶有 5 秒 timeout 的 context，避免查詢卡住。
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Always cancel to release resources

	query := `INSERT INTO todos (title, completed) VALUES ($1, $2) RETURNING id, title, completed, created_at, updated_at`

	var todo models.Todo
	// 後面傳入的參數 (todo.Title, todo.Completed, todo.CreatedAt, todo.UpdatedAt) 會對應到 SQL 裡的 $1, $2, ... 佔位符。
	// Scan(...) 會把查詢結果的欄位值依序填入 todo.ID, todo.Title, todo.Completed, todo.CreatedAt, todo.UpdatedAt。
	err := pool.QueryRow(ctx, query, todo.Title, todo.Completed).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to insert todo: %w", err)
	}

	return &todo, nil
}
