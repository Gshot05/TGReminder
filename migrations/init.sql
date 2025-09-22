CREATE TABLE IF NOT EXISTS reminders (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    frequency INTERVAL NOT NULL,
    chat_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
