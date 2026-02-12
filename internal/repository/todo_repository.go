package repository

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"todo_api/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
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

// repository層: 建立物件 → 寫入資料庫 → 回傳完整物件

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

func GetTodoByID(pool *pgxpool.Pool, id int) (*models.Todo, error) {
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
		WHERE id = $1
	`

	var todo models.Todo
	// 其實是在做「執行 SQL（只拿一筆結果）→ 把回傳欄位塞進 todo 這個 struct」
	err := pool.QueryRow(ctx, query, id).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("新增 todo 失敗: %w", err)
	}

	return &todo, nil
}

// TODO:　這邊有模擬過 read-only 的情況，將來有機會再另外整理
/*
如果你還想進一步減少 call（可選加強）

後端直接 403/429 回應
當 read-only 開啟時，對寫入相關的 method（POST/PUT/PATCH/DELETE）直接回 403 Forbidden 或 429 Too Many Requests，並附上訊息「此功能暫時唯讀」。
加上 rate limit 或 IP 暫時封鎖（進階）
如果是特定異常狀況很嚴重，可以針對該 API 做更嚴格的限流。
前端 + 後端雙重防護（推薦）
前端 read-only 防一般使用者，後端 read-only 防繞過者，兩層都做最保險。

目前這個做法已經很務實了，先把「資料不壞」守住是最重要的，其他 call 的問題相對次要（除非你已經看到有大量異常呼叫在打）。
*/
func UpdateTodo(pool *pgxpool.Pool, id int, title string, completed bool, readonlyTest bool) (*models.Todo, error) {
	const maxRetries = 1

	for attempt := 0; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 如果今天遇到突發狀況，不希望用戶對某一支特定API寫入資料，可以將該API設置為只讀模式，有可能是察覺有贓資料狀況發生，可以對
		// 前端API多添加 http://localhost:3000/todos/2?readonly_test=1 這樣的模式，強制為只讀取模式
		// ─── 測試專用：模擬 read-only ───────────────────────────────
		if readonlyTest {
			_, err := pool.Exec(ctx, "SET default_transaction_read_only = on")
			if err != nil {
				log.Printf("[TEST] Failed to set read-only: %v", err)
			} else {
				log.Printf("[TEST] Set default_transaction_read_only = on (attempt %d)", attempt+1)
			}
		}
		// ──────────────────────────────────────────────────────────────

		var query = `
            UPDATE todos
            SET title = $1, completed = $2, updated_at = CURRENT_TIMESTAMP
            WHERE id = $3
            RETURNING id, title, completed, created_at, updated_at
        `

		var todo models.Todo
		err := pool.QueryRow(ctx, query, title, completed, id).Scan(
			&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt,
		)

		if err == nil {
			return &todo, nil
		}

		// 你的 read-only 偵測邏輯...
		errStr := strings.ToLower(err.Error())
		isReadOnly := strings.Contains(errStr, "read-only") ||
			strings.Contains(errStr, "cannot execute") && strings.Contains(errStr, "read-only")

		if isReadOnly && attempt < maxRetries {
			log.Printf("偵測到 read-only 錯誤 (可能是 failover)，執行 pool.Reset() 並重試一次")
			pool.Reset()
			time.Sleep(500 * time.Millisecond)
			continue
		}

		return nil, err
	}

	return nil, fmt.Errorf("超過重試次數")
}

// 更精確的 retry 判斷（可再擴充）
func isRetryablePostgresError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())

	// 常見可重試錯誤關鍵字
	if strings.Contains(errStr, "read-only") ||
		strings.Contains(errStr, "cannot execute") && strings.Contains(errStr, "read-only") ||
		strings.Contains(errStr, "connection reset") ||
		strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "no pg_hba.conf") || // 極端情況
		strings.Contains(errStr, "tls") { // 某些 tls 重新交握失敗
		return true
	}

	// 如果你用的是 pgx/v5，也可以檢查 pgconn 錯誤碼
	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {
		case "57P03", // cannot execute during recovery
			"08006", // connection failure
			"08001": // unable to connect
			return true
		}
	}

	return false
}

/*
func (p *Pool) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	c, err := p.Acquire(ctx)  // 1. 從連接池獲取連接
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	defer c.Release()         // 2. 執行完自動還回連接池

	return c.Exec(ctx, sql, arguments...)  // 3. 執行SQL
}

*/

func DeletTodo(pool *pgxpool.Pool, id int) error {
	// 建立帶有背景上下文的連線池
	var ctx context.Context
	var cancel context.CancelFunc
	// 資料庫查詢超時，超過5秒算超時
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // 釋放記憶

	// 在資料表名稱 todos 中，對 表 的欄位新增一筆資料
	var query string = `
		DELETE FROM todos
		WHERE id = $1
	`

	/* 搜關鍵字找得到 :　how to delete item in db by using pgxpool for golang range
		deleteItem deletes a record from the "users" table by ID
	func deleteItem(ctx context.Context, pool *pgxpool.Pool, id int) (int64, error) {
		// Use Exec for non-SELECT queries
		cmdTag, err := pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	*/
	_, err := pool.Exec(ctx, query, id)
	if err != nil {
	}
	return err
}
