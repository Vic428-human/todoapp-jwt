-- 建立 article_tags 表格，如果不存在就建立
-- 篩選功能核心就在這裡，前端點 #Practical 這種 chip，後端最後就是要靠 article_tags 去查
-- 需求是一篇文章會有多個 tag，所以中間一定要獨立一張table，它就是專門存「誰跟誰有關聯」。
CREATE TABLE IF NOT EXISTS article_tags (
    article_id UUID NOT NULL, -- 文章的唯一識別碼
    tag_id UUID NOT NULL,     -- 標籤的唯一識別碼
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- 建立關聯的時間，預設為當下

    PRIMARY KEY (article_id, tag_id), -- 主鍵：避免同一篇文章重複綁定同一個標籤

    -- 外鍵：文章刪除時，相關的標籤關聯也會自動刪除
    CONSTRAINT fk_article_tags_article
        FOREIGN KEY (article_id)
        REFERENCES articles(id)
        ON DELETE CASCADE,

    -- 外鍵：標籤刪除時，相關的文章關聯也會自動刪除
    CONSTRAINT fk_article_tags_tag
        FOREIGN KEY (tag_id)
        REFERENCES tags(id)
        ON DELETE CASCADE
);

-- 在 tag_id 欄位建立索引，加速查詢某個標籤下有哪些文章
CREATE INDEX IF NOT EXISTS idx_article_tags_tag_id ON article_tags(tag_id);

-- 在 article_id 欄位建立索引，加速查詢某篇文章有哪些標籤
CREATE INDEX IF NOT EXISTS idx_article_tags_article_id ON article_tags(article_id);

-- users
-- ┌──────────────────────────────┐
-- │ id (PK)                      │
-- │ email                        │
-- │ password                     │
-- │ image_url                    │
-- │ created_at                   │
-- │ updated_at                   │
-- └──────────────────────────────┘
--               ▲
--               │ author_id
--               │
-- articles      │
-- ┌──────────────────────────────┐
-- │ id (PK)                      │
-- │ author_id (FK -> users.id)   │
-- │ title                        │
-- │ summary                      │
-- │ content                      │
-- │ difficulty                   │
-- │ status                       │
-- │ published_at                 │
-- │ like_count                   │
-- │ comment_count                │
-- │ view_count                   │
-- │ created_at                   │
-- │ updated_at                   │
-- └──────────────────────────────┘
--               │
--               │ article_id
--               ▼
-- article_tags
-- ┌──────────────────────────────┐
-- │ article_id (PK, FK)          │
-- │ tag_id     (PK, FK)          │
-- │ created_at                   │
-- └──────────────────────────────┘
--               ▲
--               │ tag_id
--               │
-- tags
-- ┌──────────────────────────────┐
-- │ id (PK)                      │
-- │ name                         │
-- │ slug                         │
-- │ is_active                    │
-- │ created_at                   │
-- │ updated_at                   │
-- └──────────────────────────────┘

-- articles 表
-- id      title
-- A1      React Hooks for Beginners : articles(id)=> article_id
-- A2      Advanced PostgreSQL Indexing : articles(id)=> article_id
-- A3      Practical Go API Design :  articles(id)=> article_id

-- tags 表
-- id      name                    slug
-- T1      Practical               practical : tags(id) => tag_id
-- T2      Beginner Friendly       beginner-friendly  : tags(id) => tag_id
-- T3      Deep Dive               deep-dive : tags(id) => tag_id

-- article_tags 表 (這張表用來關聯 tags 跟 articles 這兩張表)
-- article_id   tag_id 欄位
-- A1           T2 
-- A1           T1
-- A2           T3
-- A3           T1