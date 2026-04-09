package models

// 文章列表單筆資料長什麼樣
// 分頁 response 長什麼樣
// 先定義好，讓後面的 handler / repository 都有共同依據。

import "time"

// Article
// 用途：
// 這個 struct 比較偏向「文章主資料模型」，之後如果要做 create article / get article detail / update article，通常都會以這個 struct 當基礎。
type Article struct {
	ID           string     `json:"id" db:"id"`
	AuthorID     string     `json:"author_id" db:"author_id"`
	Title        string     `json:"title" db:"title"`
	Summary      string     `json:"summary" db:"summary"`
	Content      string     `json:"content" db:"content"`
	Difficulty   string     `json:"difficulty" db:"difficulty"`
	Status       string     `json:"status" db:"status"`
	PublishedAt  *time.Time `json:"published_at" db:"published_at"`
	LikeCount    int        `json:"like_count" db:"like_count"`
	CommentCount int        `json:"comment_count" db:"comment_count"`
	ViewCount    int        `json:"view_count" db:"view_count"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// ArticleListItem
// 用途：
// 列表頁通常不需要完整 content ，也不一定需要所有欄位，所以把列表頁需要的欄位獨立成一個 response model，會比較清楚。
type ArticleListItem struct {
	ID           string     `json:"id" db:"id"`
	Title        string     `json:"title" db:"title"`
	Summary      string     `json:"summary" db:"summary"`
	Difficulty   string     `json:"difficulty" db:"difficulty"`
	PublishedAt  *time.Time `json:"published_at" db:"published_at"`
	LikeCount    int        `json:"like_count" db:"like_count"`
	CommentCount int        `json:"comment_count" db:"comment_count"`
	ViewCount    int        `json:"view_count" db:"view_count"`

	// 先把作者資訊直接攤平在列表 item 裡
	AuthorID    string  `json:"author_id" db:"author_id"`
	AuthorEmail string  `json:"author_email" db:"author_email"`
	AuthorImage *string `json:"author_image_url" db:"author_image_url"`
}

// ArticleListResponse
// 用途：
// 列表頁的 response 格式，除了 items 以外，還會有分頁資訊。
type ArticleListResponse struct {
	Items      []ArticleListItem `json:"items"`
	Page       int               `json:"page"`
	PageSize   int               `json:"pageSize"`
	TotalCount int               `json:"totalCount"`
	TotalPages int               `json:"totalPages"`
}
