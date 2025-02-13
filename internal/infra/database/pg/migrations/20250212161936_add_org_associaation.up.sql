ALTER TABLE "teams"
ADD COLUMN "org_id" VARCHAR(255) NOT NULL,
ADD CONSTRAINT "organization_id_fkey" FOREIGN KEY ("org_id") REFERENCES "organizations" ("id") ON DELETE CASCADE;

ALTER TABLE "roles"
ADD COLUMN "org_id" VARCHAR(255) NOT NULL,
ADD CONSTRAINT "organization_id_fkey" FOREIGN KEY ("org_id") REFERENCES "organizations" ("id") ON DELETE CASCADE;

ALTER TABLE "permissions"
ADD COLUMN "org_id" VARCHAR(255) NOT NULL,
ADD CONSTRAINT "organization_id_fkey" FOREIGN KEY ("org_id") REFERENCES "organizations" ("id") ON DELETE CASCADE;
