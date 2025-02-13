ALTER TABLE "teams"
DROP CONSTRAINT IF EXISTS "organization_id_fkey",
DROP COLUMN IF EXISTS "organization_id";

ALTER TABLE "roles"
DROP CONSTRAINT IF EXISTS "organization_id_fkey",
DROP COLUMN IF EXISTS "organization_id";

ALTER TABLE "permissions"
DROP CONSTRAINT IF EXISTS "organization_id_fkey",
DROP COLUMN IF EXISTS "organization_id";
