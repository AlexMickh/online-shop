CREATE TYPE role_enum AS ENUM(
    'user',
    'admin'
);

CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(50),
    email VARCHAR(50) UNIQUE,
    password TEXT,
    role role_enum DEFAULT 'user',
    is_email_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);