CREATE INDEX IF NOT EXISTS users_textsearch_idx ON users USING gin (
	(
		to_tsvector ('simple', full_name) || to_tsvector ('simple', username)
	)
);