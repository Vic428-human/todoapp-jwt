package handlers

// 它負責：
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

// 補上你的 handler function，例如：
import (
	"net/http"
)

func GetArticlesHandler(w http.ResponseWriter, r *http.Request) {
	// 實作邏輯
}
