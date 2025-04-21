-- users table migration
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL,
    registration_date TIMESTAMPTZ NOT NULL,
    password_hash TEXT NOT NULL
); 