package repository

import (
	"context"
	"fmt"
	"time"
	"todo_api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// 超時檢查
func performOperation(ctx context.Context) {
	select {
	// 5秒才超時，所以會顯示完成
	case <-time.After(2 * time.Second):
		fmt.Println("Operation completed")
	// 如超時時間設定1秒，則會顯示超時
	case <-ctx.Done():
		fmt.Println("Operation timed out")
	}
}

// 建立物件 → 寫入資料庫 → 回傳完整物件
// 傳入的是 todo 結構體對應的json的key名稱 	(上層)todo, err := repository.CreateTodo(pool, input.Title, input.Completed)
func CreateTodo(pool *pgxpool.Pool, title string, completed bool) (*models.Todo, error) {
	// 建立帶有背景上下文的連線池
	var ctx context.Context
	var cancel context.CancelFunc
	// 資料庫查詢超時，超過5秒算超時
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // 釋放記憶體

	performOperation(ctx)

	// 在資料表名稱 todos 中，對 表 的欄位新增一筆資料
	query := `INSERT INTO todos (title, completed) VALUES ($1, $2) RETURNING id, title, completed, created_at, updated_at`

	var todo models.Todo
	// 其實是在做「執行 SQL（只拿一筆結果）→ 把回傳欄位塞進 todo 這個 struct」
	// title, completed：會依序對應到 SQL 裡的 $1, $2，也就是 VALUES ($1, $2) => 所以前端傳來的 title, completed 會依序寫入到 $1, $2
	err := pool.QueryRow(ctx, query, title, completed).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("新增 todo 失敗: %w", err)
	}

	return &todo, nil
}

func GetAllTodos(pool *pgxpool.Pool) ([]models.Todo, error) {

	// 建立帶有背景上下文的連線池
	var ctx context.Context
	var cancel context.CancelFunc
	// 資料庫查詢超時，超過5秒算超時
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // 釋放記憶

	performOperation(ctx)

	// 在資料表名稱 todos 中，對 表 的欄位新增一筆資料
	var query string = `
		SELECT id, title, completed, created_at, updated_at
		FROM todos
		ORDER BY created_at DESC
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("查詢 todos 失敗: %w", err)
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, fmt.Errorf("讀取 todo 失敗: %w", err)
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("讀取 todos 失敗: %w", err)
	}

	return todos, nil
}
