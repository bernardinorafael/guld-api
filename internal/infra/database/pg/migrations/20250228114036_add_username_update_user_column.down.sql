ALTER TABLE users
DROP COLUMN "ignore_password_policy",
DROP COLUMN "username_last_updated",
DROP COLUMN "username_lockout_end";
