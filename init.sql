DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'pvz_service') THEN
        CREATE DATABASE pvz_service;
    END IF;
END $$; 