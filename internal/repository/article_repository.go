package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"todo_api/internal/models"
	"todo_api/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// GetArticles
// 用途：
// 這支 function 負責查文章列表資料。
// 目前第一版先支援
// /articles?page=1&pageSize=5
// /articles?page=1&pageSize=5&difficulty=beginner
func GetArticles(pool *pgxpool.Pool, page int, pageSize int, tag string, difficulty string) (*models.ArticleListResponse, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	// 設定 DB 查詢逾時時間，避免請求卡太久
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	utils.PerformOperation(ctx)

	// 從第幾筆資料開始取
	// page=1, pageSize=5 -> offset=0（從第1筆開始）
	// page=2, pageSize=5 -> offset=5（從第6筆開始）
	offset := (page - 1) * pageSize

	conditions := []string{"a.status = 'published'"}

	args := []interface{}{}
	argIndex := 1
	// difficulty 不存在=> whereClause = "a.status = 'published'"
	// difficulty 存在 => whereClause = "a.status = 'published' AND a.difficulty = $1"
	if difficulty != "" {
		conditions = append(conditions, fmt.Sprintf("a.difficulty = $%d", argIndex)) //  "a.difficulty = $1"
		args = append(args, difficulty)
		argIndex++
	}

	// 把多個條件組成 WHERE 子句
	whereClause := strings.Join(conditions, " AND ")

	// =========================
	// 1. 先查總筆數
	// =========================
	var totalCount int
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM articles a
		WHERE %s
	`, whereClause)

	err := pool.QueryRow(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("查詢 articles 總數失敗: %w", err)
	}

	// =========================
	// 2. 再查當前頁資料
	// 先跳過前 OFFSET 筆的資料（也就是「從第幾筆開始」查
	// 再取最多 LIMIT 筆資料（也就是「最多查幾筆」
	// =========================
	listQuery := fmt.Sprintf(`
		SELECT
			a.id,
			a.title,
			a.summary,
			a.difficulty,
			a.published_at,
			a.like_count,
			a.comment_count,
			a.view_count,
			u.id AS author_id,
			u.email AS author_email,
			u.image_url AS author_image_url
		FROM articles a
		JOIN users u ON a.author_id = u.id
		WHERE %s
		ORDER BY a.published_at DESC, a.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	// 例子，假設：
	// args 原先是 [difficulty]
	// pageSize = 5
	// offset = 0
	// listArgs := append(args, pageSize, offset)
	// listArgs = []interface{}{"advanced", 5, 0}
	// listArgs: [advanced 5 0]
	listArgs := append(args, pageSize, offset)
	// 再把 listArgs 傳給 Query 執行 SQL 查詢
	rows, err := pool.Query(ctx, listQuery, listArgs...)

	if err != nil {
		return nil, fmt.Errorf("查詢 articles 失敗: %w", err)
	}
	defer rows.Close()

	var articles []models.ArticleListItem

	for rows.Next() {
		var article models.ArticleListItem

		if err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Summary,
			&article.Difficulty,
			&article.PublishedAt,
			&article.LikeCount,
			&article.CommentCount,
			&article.ViewCount,
			&article.AuthorID,
			&article.AuthorEmail,
			&article.AuthorImage,
		); err != nil {
			return nil, fmt.Errorf("讀取 article 失敗: %w", err)
		}

		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("讀取 articles 失敗: %w", err)
	}

	// 向上取整，算總頁數 => 每一頁顯示多少筆 (pageSize) 跟 知道總筆數(totalCount)，可以知道總頁數，
	// 根據不同 whereClause 條件，總筆數都不同，所以總頁數也會不同。
	totalPages := (totalCount + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	response := &models.ArticleListResponse{
		Items:      articles,
		Page:       page, // 第幾頁
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}

	return response, nil
}

// 佔位符
// %s  → 字串
// %d  → 整數
// %f  → 浮點數
// %v  → 自動判斷類型

// fmt.Sprintf 動態建構一個 SQL 查詢字串
