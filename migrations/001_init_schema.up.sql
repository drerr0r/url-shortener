-- Миграция для создания таблицы URLs
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code VARCHAR(12) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    access_count INTEGER DEFAULT 0
);

-- 🟡 ДОБАВЛЕНО: Индексы для улучшения производительности
CREATE UNIQUE INDEX idx_urls_short_code ON urls(short_code);
CREATE INDEX idx_urls_created_at ON urls(created_at);
CREATE INDEX idx_urls_original_url ON urls(original_url);

-- 🟡 ДОБАВЛЕНО: Комментарии к таблице и колонкам для документации
COMMENT ON TABLE urls IS 'Таблица для хранения сокращенных URL';
COMMENT ON COLUMN urls.original_url IS 'Оригинальный URL';
COMMENT ON COLUMN urls.short_code IS 'Сокращенный код URL';
COMMENT ON COLUMN urls.created_at IS 'Время создания записи';
COMMENT ON COLUMN urls.access_count IS 'Количество переходов по ссылке';