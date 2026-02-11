package repository

import (
	"context"
	"fmt"
	"time"
	"todo_api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// 建立物件 → 寫入資料庫 → 回傳完整物件
// 傳入的是 todo 結構體對應的json的key名稱 	(上層)todo, err := repository.CreateTodo(pool, input.Title, input.Completed)
func CreateTodo(pool *pgxpool.Pool, title string, completed bool) (*models.Todo, error) {
	// 建立帶有背景上下文的連線池
	var ctx context.Context
	var cancel context.CancelFunc
	// 帶有 5 秒 timeout 的 context，避免查詢卡住。
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // 釋放記憶體

	// 在資料表名稱 todos 中，對 表 的欄位新增一筆資料
	query := `INSERT INTO todos (title, completed) VALUES ($1, $2) RETURNING id, title, completed, created_at, updated_at`

	var todo models.Todo
	// 其實是在做「執行 SQL（只拿一筆結果）→ 把回傳欄位塞進 todo 這個 struct」
	// title, completed：會依序對應到 SQL 裡的 $1, $2，也就是 VALUES ($1, $2) => 所以前端傳來的 title, completed 會依序寫入到 $1, $2
	err := pool.QueryRow(ctx, query, title, completed).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to insert todo: %w", err)
	}

	return &todo, nil
}
