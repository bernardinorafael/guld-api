CREATE TABLE
	IF NOT EXISTS "organization_settings" (
		"id" VARCHAR(255) PRIMARY KEY NOT NULL,
		"org_id" VARCHAR(255) NOT NULL,
		"is_active" BOOLEAN NOT NULL,
		"default_membership_password" VARCHAR(255) NOT NULL,
		"max_allowed_memberships" INTEGER NOT NULL,
		"max_allowed_roles" INTEGER NOT NULL,
		"use_master_password" BOOLEAN NOT NULL,
		"created" TIMESTAMPTZ NOT NULL DEFAULT NOW (),
		"updated" TIMESTAMPTZ NOT NULL DEFAULT NOW (),
		CONSTRAINT "organization_id_fkey" FOREIGN KEY ("org_id") REFERENCES "organizations" ("id") ON DELETE CASCADE
	);