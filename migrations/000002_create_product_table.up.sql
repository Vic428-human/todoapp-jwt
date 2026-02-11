-- 資料表的欄位名稱 owner_id, title, game, platform, username, views, monthly_views, price, description, verified, country, featured ...等
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    owner_id VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    game VARCHAR(100) NOT NULL,
    platform VARCHAR(50) NOT NULL,
    username VARCHAR(50) NOT NULL,
    views INTEGER DEFAULT 0,
    monthly_views INTEGER DEFAULT 0,
    price INTEGER NOT NULL,
    description TEXT,
    verified BOOLEAN DEFAULT FALSE,
    country VARCHAR(50),
    featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
