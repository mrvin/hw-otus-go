CREATE TABLE IF NOT EXISTS users (
	id serial PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	hash_password TEXT,
	role TEXT,
	email TEXT
);
CREATE INDEX IF NOT EXISTS idx_name ON users(name);

CREATE TABLE IF NOT EXISTS events (
	id serial primary key,
	title text,
	description text,
	start_time timestamptz,
	stop_time timestamptz,
	user_id integer references users(id) on delete cascade
);
