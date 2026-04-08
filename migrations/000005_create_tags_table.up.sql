CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),   -- 唯一識別碼，使用 UUID，自動生成
    name VARCHAR(100) NOT NULL,                      -- 標籤名稱，必填，最多 100 字元
    slug VARCHAR(100) NOT NULL UNIQUE,               -- 標籤的唯一代稱（通常用於 URL），必填且唯一
    is_active BOOLEAN NOT NULL DEFAULT TRUE,         -- 標籤是否啟用，預設為 TRUE
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- 建立時間，預設為當前時間
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP  -- 更新時間，預設為當前時間
);
