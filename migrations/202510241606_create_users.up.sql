CREATE TYPE log_level AS ENUM ('INFO', 'WARN', 'ERROR');

CREATE TABLE IF NOT EXISTS users_mails (
    tg_id BIGINT PRIMARY KEY,
    email VARCHAR(100) NOT NULL UNIQUE
    -- при желании можно добавить ещё поля (имя, дата регистрации и т.д.)
);

CREATE TABLE IF NOT EXISTS users_filters (
    tg_id BIGINT PRIMARY KEY REFERENCES users_mails(tg_id),
    service TEXT,
    level log_level
);
