-- Enable the dblink extension
CREATE EXTENSION IF NOT EXISTS dblink;

-- Your existing logic
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'socialnetwork') THEN
        PERFORM dblink_exec('dbname=postgres', 'CREATE DATABASE socialnetwork');
    END IF;
END $$;

GRANT ALL PRIVILEGES ON DATABASE socialnetwork TO admin;
