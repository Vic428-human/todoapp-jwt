CREATE TABLE IF NOT EXISTS articles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),   -- 多台伺服器同時建立文章，不需要擔心 ID 重複，但在 B-Tree 索引中，UUID 的隨機性會導致資料分布不連續，影響查詢效能。
    author_id UUID NOT NULL,                         -- 作者的 ID，必須存在，連結到 users 表
    title VARCHAR(255) NOT NULL,                     -- 文章標題，限制 255 字元
    summary TEXT NOT NULL,                           -- 文章摘要
    content TEXT NOT NULL,                           -- 文章完整內容
    difficulty VARCHAR(50) NOT NULL,                 -- 難度標記 (beginner/intermediate/advanced)
    status VARCHAR(50) NOT NULL DEFAULT 'draft',     -- 狀態 (draft/published/archived)，預設 draft
    published_at TIMESTAMP WITH TIME ZONE,           -- 發布時間，含時區
    like_count INT NOT NULL DEFAULT 0,               -- 按讚數，預設 0
    comment_count INT NOT NULL DEFAULT 0,            -- 留言數，預設 0
    view_count INT NOT NULL DEFAULT 0,               -- 瀏覽數，預設 0
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- 建立時間
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- 更新時間

    -- 外鍵約束：作者必須存在於 users 表
    -- ON DELETE CASCADE 的作用 → 當作者被刪除時，相關文章也會自動刪除
    CONSTRAINT fk_articles_author
        FOREIGN KEY (author_id)
        REFERENCES users(id)
        ON DELETE CASCADE, 

    -- 檢查約束：限制難度只能是三種
    CONSTRAINT chk_articles_difficulty
        CHECK (difficulty IN ('beginner', 'intermediate', 'advanced')),

    -- 檢查約束：限制狀態只能是三種
    CONSTRAINT chk_articles_status
        CHECK (status IN ('draft', 'published', 'archived')),

    -- 檢查約束：避免負數
    CONSTRAINT chk_articles_like_count
        CHECK (like_count >= 0),

    CONSTRAINT chk_articles_comment_count
        CHECK (comment_count >= 0),

    CONSTRAINT chk_articles_view_count
        CHECK (view_count >= 0)
);

CREATE INDEX IF NOT EXISTS idx_articles_author_id ON articles(author_id);       -- 快速查詢某作者的文章
CREATE INDEX IF NOT EXISTS idx_articles_status ON articles(status);             -- 快速查詢不同狀態的文章
CREATE INDEX IF NOT EXISTS idx_articles_difficulty ON articles(difficulty);     -- 快速查詢不同難度的文章
CREATE INDEX IF NOT EXISTS idx_articles_published_at ON articles(published_at DESC); -- 依發布時間排序
CREATE INDEX IF NOT EXISTS idx_articles_created_at ON articles(created_at DESC);     -- 依建立時間排序

