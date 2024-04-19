-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- Create the 'users' table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL
    );

-- Create the 'tigers' table
CREATE TABLE IF NOT EXISTS tigers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    last_seen TIMESTAMP NOT NULL,
    lat DOUBLE PRECISION NOT NULL,
    long DOUBLE PRECISION NOT NULL
    );

-- Create the 'sightings' table
CREATE TABLE IF NOT EXISTS tiger_sightings (
    id SERIAL PRIMARY KEY,
    tiger_id INTEGER REFERENCES tigers(id) ON DELETE CASCADE,
    timestamp TIMESTAMP NOT NULL,
    lat DOUBLE PRECISION NOT NULL,
    long DOUBLE PRECISION NOT NULL,
    image BYTEA,
    reporter_email VARCHAR(255)
    );

-- Indexes for improved performance (optional, but recommended for large datasets)
CREATE INDEX IF NOT EXISTS idx_tigers_last_seen ON tigers (last_seen DESC);
CREATE INDEX IF NOT EXISTS idx_tiger_sightings_timestamp ON tiger_sightings (timestamp DESC);

-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back

-- Drop the 'sightings' table
DROP TABLE IF EXISTS tiger_sightings;

-- Drop the 'tigers' table
DROP TABLE IF EXISTS tigers;

-- Drop the 'users' table
DROP TABLE IF EXISTS users;
