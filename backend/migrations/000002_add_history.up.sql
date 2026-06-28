CREATE TABLE url_shortener.clicks (
    id          SERIAL PRIMARY KEY,
    short_url   VARCHAR(64) NOT NULL REFERENCES url_shortener.urls(short_url) ON DELETE CASCADE,
    ip          VARCHAR(45) NOT NULL,
    clicked_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_clicks_short_url ON url_shortener.clicks(short_url);
