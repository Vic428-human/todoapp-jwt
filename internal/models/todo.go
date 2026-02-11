package models

import (
	"time"
)

/*
前端收到：
{
  "id": 123,
  "title": "代辦1"   ← 前端看到的 JSON
}

資料庫儲存：
id = 123
title = "代辦1"    ← DB 看到的欄位名
*/
// json 對應前端 API
// db 對應資料庫表格中的
type Todo struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title"  db:"title"`
	Completed bool      `json:"completed" db:"completed"`
	CreatedAt time.Time `json:"created_at" db:"created_at" `
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
