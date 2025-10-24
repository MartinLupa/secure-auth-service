-- Creating table for users
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Seeding users table with hashed passwords
INSERT INTO users (full_name, email, password) VALUES
-- password is the bcrypt hash of 'user'
    ('Martin Lupa', 'user@example.com', '8143e4d37076c5892453f262a3f349e2d273525b3fa096290f7db073e35e3472')
ON CONFLICT (email) DO NOTHING;