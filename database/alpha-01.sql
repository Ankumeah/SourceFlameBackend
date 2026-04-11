CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  created_at INT NOT NULL,
  password_hash BYTEA NOT NULL,
  salt BYTEA NOT NULL
);

CREATE TABLE IF NOT EXISTS repos (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  private BOOLEAN NOT NULL,
  owner_id INT REFERENCES users(id) NOT NULL,
  created_at INT NOT NULL,
  stars INT NOT NULL
);

CREATE TABLE IF NOT EXISTS pats (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  hash BYTEA NOT NULL,
  owner_id INT REFERENCES users(id) NOT NULL,
  created_at INT NOT NULL,
);
