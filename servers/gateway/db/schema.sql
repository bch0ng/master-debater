CREATE DATABASE IF NOT EXISTS master_debater;
  USE master_debater;
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    passhash VARCHAR(255) NOT NULL,
    usernname VARCHAR(255) NOT NULL UNIQUE
);
