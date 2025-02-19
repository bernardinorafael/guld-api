CREATE TABLE
	IF NOT EXISTS "email_validations" (
		"id" varchar(255) PRIMARY KEY NOT NULL,
		"email_id" varchar(255) NOT NULL,
		"attempts" int DEFAULT 0,
		"is_consumed" boolean DEFAULT FALSE,
		"is_valid" boolean DEFAULT FALSE,
		"created" timestamptz NOT NULL DEFAULT now (),
		"expires" timestamptz NOT NULL,
		CONSTRAINT "validations_emails_id_fkey" FOREIGN KEY ("email_id") REFERENCES "emails" ("id") ON DELETE CASCADE
	);