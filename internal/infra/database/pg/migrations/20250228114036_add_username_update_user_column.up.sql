ALTER TABLE "users"
ADD COLUMN "username_last_updated" TIMESTAMPTZ,
ADD COLUMN "username_lockout_end" TIMESTAMPTZ;