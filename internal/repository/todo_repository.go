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
	// 保護查詢不會無限等待，只等5秒
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Always cancel to release resources

	// 建立待辦事項
	todo := &models.Todo{Title: title, Completed: completed, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	// 執行 SQL 插入語句，RETURNING id 以取得主鍵
	query := `INSERT INTO todos (title, completed, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id `
	// 使用 QueryRow(...).Scan(...) 來取得 RETURNING id，確保 todo.ID 被正確填入
	err := pool.QueryRow(ctx, query, todo.Title, todo.Completed, todo.CreatedAt, todo.UpdatedAt).Scan(&todo.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert todo: %w", err)
	}

	return todo, nil
}
