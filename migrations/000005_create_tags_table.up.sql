CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),   -- 唯一識別碼，使用 UUID，自動生成
    name VARCHAR(100) NOT NULL,                      -- 標籤名稱，必填，最多 100 字元
    slug VARCHAR(100) NOT NULL UNIQUE,               -- 標籤的唯一代稱（通常用於 URL），必填且唯一
    is_active BOOLEAN NOT NULL DEFAULT TRUE,         -- 標籤是否啟用，預設為 TRUE
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- 建立時間，預設為當前時間
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP  -- 更新時間，預設為當前時間
);

-- 模擬前端點選，跟 tags 後端表的交互關係，以及如何對 article_tags 進行資料的篩選跟撈出
-- 點選 1 個 chip 的情況
-- [前端點擊 #Practical]
--         │
--         ▼
-- 送出請求
-- GET /articles?tags=practical
--         │
--         ▼
-- 後端先去 tags 找 slug = 'practical'
--         │
--         ▼
-- 找到 tags.id = T1
--         │
--         ▼
-- 再去 article_tags 找 tag_id = T1
--         │
--         ▼
-- 找到 article_id:
-- A1, A3
--         │
--         ▼
-- 再去 articles 撈出
-- A1, A3 的文章資料
--         │
--         ▼
-- 回傳給前端

-- 點選 2 個以上 chip 的情況，只要有其中一個 tag 就算符合，使用 OR 邏輯
-- GET /articles?tags=practical,beginner-friendly
-- A1 -> T1 (practical)
-- A1 -> T2 (beginner-friendly)
-- A3 -> T1 (practical)