package models

import "time"

/*
前端收到：
{
  "id": 123,
  "ownerId": "owner001"   ← 前端看到的 JSON
}

資料庫儲存：
id = 123
owner_id = "owner001"    ← DB 看到的欄位名
*/
// json 對應前端 API
// db 對應資料庫表格中的
type Product struct {
	ID           int       `json:"id" db:"id"`
	OwnerID      string    `json:"ownerId" db:"owner_id"`
	Title        string    `json:"title" db:"title"`
	Game         string    `json:"game" db:"game"`
	Platform     string    `json:"platform" db:"platform"`
	Username     string    `json:"username" db:"username"`
	Views        int       `json:"views" db:"views"`
	MonthlyViews int       `json:"monthly_views" db:"monthly_views"`
	Price        int       `json:"price" db:"price"`
	Description  string    `json:"description" db:"description"`
	Verified     bool      `json:"verified" db:"verified"`
	Country      string    `json:"country" db:"country"`
	Featured     bool      `json:"featured" db:"featured"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}
