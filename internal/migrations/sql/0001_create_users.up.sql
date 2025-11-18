CREATE TABLE IF NOT EXISTS users (
                                     id          UUID PRIMARY KEY,
                                     name        TEXT        NOT NULL,
                                     email       TEXT        NOT NULL UNIQUE,
                                     phone       TEXT        NOT NULL DEFAULT '',
                                     created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users (created_at);