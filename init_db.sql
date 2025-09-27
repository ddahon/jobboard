CREATE TABLE companies (
	id INTEGER,
	name TEXT,
	description TEXT,
	website TEXT,
	shortname TEXT
);

CREATE TABLE jobs (
	description TEXT,
	title TEXT,
	link TEXT,
	location TEXT,
	category INTEGER,
	created_at DATETIME DEFAULT (CURRENT_TIMESTAMP),
	updated_at DATETIME DEFAULT (CURRENT_TIMESTAMP),
	company_id INTEGER
);
