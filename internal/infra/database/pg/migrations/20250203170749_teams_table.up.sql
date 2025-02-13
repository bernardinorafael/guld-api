CREATE TABLE
	IF NOT EXISTS "teams" (
		"id" VARCHAR(255) PRIMARY KEY NOT NULL,
		"name" VARCHAR(255) NOT NULL,
		"slug" VARCHAR(255) UNIQUE NOT NULL,
		"owner_id" VARCHAR(255) NOT NULL,
		"logo" VARCHAR(255) NULL,
		"created" TIMESTAMPTZ NOT NULL DEFAULT NOW (),
		"updated" TIMESTAMPTZ NOT NULL DEFAULT NOW (),
		CONSTRAINT "owner_id_fkey" FOREIGN KEY ("owner_id") REFERENCES "users" ("id") ON DELETE CASCADE
	);

CREATE TABLE
	IF NOT EXISTS "roles" (
		"id" VARCHAR(255) PRIMARY KEY NOT NULL,
		"team_id" VARCHAR(255) NOT NULL,
		"name" VARCHAR(255) NOT NULL UNIQUE,
		"key" VARCHAR(255) NOT NULL UNIQUE,
		"description" VARCHAR(255) NOT NULL,
		"created" TIMESTAMPTZ NOT NULL DEFAULT NOW (),
		"updated" TIMESTAMPTZ NOT NULL DEFAULT NOW (),
		CONSTRAINT "team_id_fkey" FOREIGN KEY ("team_id") REFERENCES "teams" ("id") ON DELETE CASCADE
	);

CREATE TABLE
	IF NOT EXISTS "permissions" (
		"id" VARCHAR(255) PRIMARY KEY NOT NULL,
		"team_id" VARCHAR(255) NOT NULL,
		"name" VARCHAR(255) NOT NULL UNIQUE,
		"key" VARCHAR(255) NOT NULL UNIQUE,
		"description" VARCHAR(255) NOT NULL,
		"created" TIMESTAMPTZ NOT NULL DEFAULT NOW (),
		"updated" TIMESTAMPTZ NOT NULL DEFAULT NOW (),
		CONSTRAINT "team_id_fkey" FOREIGN KEY ("team_id") REFERENCES "teams" ("id") ON DELETE CASCADE
	);

CREATE TABLE
	IF NOT EXISTS "role_permissions" (
		"id" VARCHAR(255) PRIMARY KEY NOT NULL,
		"role_id" VARCHAR(255) NOT NULL,
		"permission_id" VARCHAR(255) NOT NULL,
		CONSTRAINT "role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON DELETE CASCADE,
		CONSTRAINT "permission_id_fkey" FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id") ON DELETE CASCADE,
		UNIQUE ("role_id", "permission_id")
	);

CREATE TABLE
	IF NOT EXISTS "team_members" (
		"id" VARCHAR(255) PRIMARY KEY NOT NULL,
		"team_id" VARCHAR(255) NOT NULL,
		"user_id" VARCHAR(255) NOT NULL,
		"role_id" VARCHAR(255) NOT NULL,
		"joined" TIMESTAMPTZ NOT NULL DEFAULT NOW (),
		UNIQUE ("team_id", "user_id"),
		CONSTRAINT "team_id_fkey" FOREIGN KEY ("team_id") REFERENCES "teams" ("id") ON DELETE CASCADE,
		CONSTRAINT "user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE,
		CONSTRAINT "role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON DELETE CASCADE
	);
