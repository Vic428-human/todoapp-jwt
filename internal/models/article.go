package models

import "time"

type Article struct {
	ID           string     `json:"id" db:"id"`
	AuthorID     string     `json:"author_id" db:"author_id"`
	Title        string     `json:"title" db:"title"`
	Summary      string     `json:"summary" db:"summary"`
	Content      string     `json:"content" db:"content"`
	Difficulty   string     `json:"difficulty" db:"difficulty"`
	Status       string     `json:"status" db:"status"`
	PublishedAt  *time.Time `json:"published_at" db:"published_at"` // 把 PublishedAt 設成 *time.Time，因為 draft 狀態可能是 NULL
	LikeCount    int        `json:"like_count" db:"like_count"`
	CommentCount int        `json:"comment_count" db:"comment_count"`
	ViewCount    int        `json:"view_count" db:"view_count"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}
