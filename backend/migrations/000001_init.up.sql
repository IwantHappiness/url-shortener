CREATE SCHEMA url_shortener;

CREATE TABLE url_shortener.users (
    id SERIAL PRIMARY KEY,
    version INTEGER NOT NULL DEFAULT 1,
    nickname VARCHAR(20) UNIQUE NOT NULL CHECK (char_length(nickname) BETWEEN 3 AND 20),
    email VARCHAR(254) UNIQUE NOT NULL CHECK (char_length(email) >= 3),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE url_shortener.urls (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES url_shortener.users(id) ON DELETE CASCADE,
    original_url TEXT NOT NULL CHECK (char_length(original_url) >= 3),
    short_url VARCHAR(64) UNIQUE NOT NULL CHECK (char_length(short_url) >= 3),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
