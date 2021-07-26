-- Create articles table if not exists.
CREATE TABLE IF NOT EXISTS articles (
	slug VARCHAR(255) NOT NULL,
	title VARCHAR(255) NOT NULL,
	content MEDIUMTEXT NOT NULL,
	preview TEXT NOT NULL,
	categories VARCHAR(64) NOT NULL,
	tags VARCHAR(64) NOT NULL,
	pinned BOOLEAN DEFAULT FALSE,
	draft BOOLEAN DEFAULT FALSE,
	published_at DATETIME NOT NULL,
	PRIMARY KEY (slug)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
