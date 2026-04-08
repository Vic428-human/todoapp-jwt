// 收 query params
// 例如讀到 tags=practical

// GET /articles?tags=practical
//         │
//         ▼
// handlers.GetArticlesHandler(...)
//         │
//         ▼
// repository.GetArticles(pool, filters)
//         │
//         ▼
// SQL 去查：
// tags.slug = 'practical'
// → article_tags.tag_id
// → article_tags.article_id
// → articles

package handlers

import (
	"net/http"
	"strconv"

	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// 用途：
// 這支 handler 是文章列表 API 的入口，它負責：
// 1. 從 query string 讀取 page、pageSize、tag、difficulty
// 2. 做基本型別轉換與防呆
// 3. 呼叫 repository 去查資料
// 4. 把結果回傳給前端

// 同時篩選標籤 + 難度
// GET /api/v1/articles?page=1&pageSize=10&tag=practical&difficulty=beginner
func GetArticlesHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageStr := c.DefaultQuery("page", "1")

		pageSizeStr := c.DefaultQuery("pageSize", "5")

		tag := c.Query("tag")

		difficulty := c.Query("difficulty")

		// page 跟 pageSize 需要從 string 轉成 int，並且做基本的防呆
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 {
			pageSize = 5
		}

		result, err := repository.GetArticles(pool, page, pageSize, tag, difficulty)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 查詢成功，直接把結果回傳給前端
		c.JSON(http.StatusOK, result)
	}
}
