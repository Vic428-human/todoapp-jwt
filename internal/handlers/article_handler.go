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