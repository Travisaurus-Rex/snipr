CREATE SCHEMA IF NOT EXISTS analytics;

CREATE TABLE analytics.clicks (
    id SERIAL PRIMARY KEY,
    url_id INT NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
    clicked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    referrer TEXT,
    user_agent TEXT,
    ip_address INET
)