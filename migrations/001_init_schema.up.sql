--migrations/001_init_schema.up.sql

--Создаем таблицу для хранения сокращенных ссылок
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code VARCHAR(10) NOT NULL UNIQUE
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    click_count BEGINT NOT NULL DEFAULT 0
);
-- Создаем индекс для быстрого поиска по short_code
CREATE INDEX idx_urls_short_code ON urls(short_code);

-- Создаем индекс для оригинальных URL
CREATE INDEX idx_urls_original_url ON urls(original_url)
