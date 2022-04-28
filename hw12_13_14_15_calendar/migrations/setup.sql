CREATE DATABASE event-db;

CREATE TABLE users (
	id serial primary key,
	name text
);

CREATE TABLE events (
	id serial primary key,
	title text,
	description text,
	start_time timestamp,
	stop_time timestamp,
	user_id integer references users(id) on delete cascade
--	send_time bigint,
--	created_at timestamp,
--	updated_at timestamp
);

