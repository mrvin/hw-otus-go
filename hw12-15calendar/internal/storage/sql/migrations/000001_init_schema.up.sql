CREATE TABLE IF NOT EXISTS users (
	id UUID DEFAULT gen_random_uuid (),
	name TEXT NOT NULL UNIQUE,
	hash_password TEXT,
	email TEXT,
	role TEXT,
	PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_name ON users(name);

CREATE TABLE IF NOT EXISTS events (
	id serial primary key,
	title text,
	description text,
	start_time timestamptz,
	stop_time timestamptz,
	user_id UUID references users(id) on delete cascade
);
