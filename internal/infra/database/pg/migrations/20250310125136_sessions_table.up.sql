CREATE TABLE
	"sessions" (
		"id" varchar(255) PRIMARY KEY NOT NULL,
		"username" varchar(255) NOT NULL,
		"refresh_token" varchar(512) NOT NULL,
		"agent" varchar(255) NOT NULL,
		"ip" varchar(255) NOT NULL,
		"revoked" boolean NOT NULL DEFAULT FALSE,
		"expires" timestamptz NOT NULL,
		"created" timestamptz NOT NULL DEFAULT now (),
		"updated" timestamptz NOT NULL DEFAULT now ()
	);

CREATE INDEX idx_refresh_token ON "sessions" ("refresh_token");

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") ON DELETE CASCADE;