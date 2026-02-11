package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"
	"todo_api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newTraceID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b) // 失敗就用全 0 也無所謂，trace 不該影響主流程
	return hex.EncodeToString(b)
}

// 中文版 trace
func performOperation(ctx context.Context, traceID string, op string, fields map[string]any) func(err error) {
	start := time.Now()

	// ⬇️ 開始紀錄
	log.Printf("🔍 [追蹤編號:%s] 開始執行 %s | 參數:%v", traceID, op, fields)

	return func(err error) {
		耗時 := time.Since(start)
		ctx錯誤 := ctx.Err()

		if err != nil {
			log.Printf("❌ [追蹤編號:%s] %s 執行失敗 | 耗時:%s | 錯誤:%v | Context狀態:%v",
				traceID, op, 耗時, err, ctx錯誤)
			return
		}

		log.Printf("✅ [追蹤編號:%s] %s 執行成功 | 耗時:%s | Context狀態:%v",
			traceID, op, 耗時, ctx錯誤)
	}
}

// 建立物件 → 寫入資料庫 → 回傳完整物件
// 傳入的是 todo 結構體對應的json的key名稱 	(上層)todo, err := repository.CreateTodo(pool, input.Title, input.Completed)
func CreateTodo(pool *pgxpool.Pool, title string, completed bool) (*models.Todo, error) {
	// 建立帶有背景上下文的連線池
	var ctx context.Context
	var cancel context.CancelFunc
	// 帶有 5 秒 timeout 的 context，避免查詢卡住。
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // 釋放記憶體

	traceID := newTraceID()

	traceDone := performOperation(ctx, traceID, "新增待辦事項(CreateTodo)", map[string]any{
		"title":     title,
		"completed": completed,
		"timeout秒數": 5 * time.Second,
	})

	// 在資料表名稱 todos 中，對 表 的欄位新增一筆資料
	query := `INSERT INTO todos (title, completed) VALUES ($1, $2) RETURNING id, title, completed, created_at, updated_at`

	var todo models.Todo
	// 其實是在做「執行 SQL（只拿一筆結果）→ 把回傳欄位塞進 todo 這個 struct」
	// title, completed：會依序對應到 SQL 裡的 $1, $2，也就是 VALUES ($1, $2) => 所以前端傳來的 title, completed 會依序寫入到 $1, $2
	err := pool.QueryRow(ctx, query, title, completed).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	// ✅ trace：結束（記錄成功/失敗、耗時）
	traceDone(err)

	if err != nil {
		return nil, fmt.Errorf("新增 todo 失敗: %w", err)
	}

	return &todo, nil
}
