CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- UUID => very long string 
    email VARCHAR(255) UNIQUE NOT NULL, -- VARCHAR => variable number of characters
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP -- forget password or change email, we need this column
)

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)