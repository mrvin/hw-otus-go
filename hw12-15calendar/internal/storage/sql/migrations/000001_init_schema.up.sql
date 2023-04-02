CREATE TABLE IF NOT EXISTS users (
	id serial primary key,
	name text,
	email text
);

CREATE TABLE IF NOT EXISTS events (
	id serial primary key,
	title text,
	description text,
	start_time timestamptz,
	stop_time timestamptz,
	user_id integer references users(id) on delete cascade
);
