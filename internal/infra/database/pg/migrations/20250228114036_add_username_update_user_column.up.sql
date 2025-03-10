ALTER TABLE "users"
ADD COLUMN "ignore_password_policy" BOOLEAN DEFAULT FALSE,
ADD COLUMN "username_last_updated" TIMESTAMPTZ,
ADD COLUMN "username_lockout_end" TIMESTAMPTZ;
