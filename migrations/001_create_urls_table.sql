-- +goose Up
CREATE TABLE IF NOT EXISTS urls (
    short_url VARCHAR(255) PRIMARY KEY,
    long_url TEXT NOT NULL,
    hits INT NOT NULL DEFAULT 0,
    last_hit_at TIMESTAMP,
    CONSTRAINT uc_url_combination UNIQUE (short_url, long_url)
);

-- +goose Down
DROP TABLE urls;
