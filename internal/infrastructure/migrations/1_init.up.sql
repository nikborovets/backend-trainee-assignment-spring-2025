DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'pvz_service') THEN
        CREATE DATABASE pvz_service;
    END IF;
END $$; 
-- users table migration
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL,
    registration_date TIMESTAMPTZ NOT NULL,
    password_hash TEXT NOT NULL
);

-- pvz table migration
CREATE TABLE IF NOT EXISTS pvz (
    id UUID PRIMARY KEY,
    registration_date TIMESTAMPTZ NOT NULL,
    city TEXT NOT NULL
);

-- reception table migration
CREATE TABLE IF NOT EXISTS reception (
    id UUID PRIMARY KEY,
    pvz_id UUID NOT NULL REFERENCES pvz(id),
    status TEXT NOT NULL,
    date_time TIMESTAMPTZ NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS one_open_reception_per_pvz ON reception(pvz_id) WHERE status = 'in_progress';

-- product table migration
CREATE TABLE IF NOT EXISTS product (
    id UUID PRIMARY KEY,
    reception_id UUID NOT NULL REFERENCES reception(id),
    type TEXT NOT NULL,
    date_time TIMESTAMPTZ NOT NULL
);