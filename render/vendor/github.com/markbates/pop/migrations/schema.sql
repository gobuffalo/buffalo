CREATE TABLE "schema_migration" (
"version" TEXT NOT NULL
);
CREATE UNIQUE INDEX "version_idx" ON "schema_migration" (version);
CREATE TABLE "songs" (
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"id" TEXT PRIMARY KEY,
"title" TEXT NOT NULL
);
CREATE TABLE "users" (
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"name" TEXT NOT NULL,
"alive" NUMERIC,
"birth_date" DATETIME,
"bio" text,
"price" numeric DEFAULT '1.00',
"email" TEXT NOT NULL DEFAULT 'foo@example.com',
"id" INTEGER PRIMARY KEY AUTOINCREMENT
);
CREATE TABLE "good_friends" (
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"first_name" TEXT NOT NULL,
"last_name" TEXT NOT NULL,
"id" INTEGER PRIMARY KEY AUTOINCREMENT
);
CREATE TABLE "validatable_cars" (
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"name" TEXT NOT NULL,
"id" INTEGER PRIMARY KEY AUTOINCREMENT
);
CREATE TABLE "not_validatable_cars" (
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"name" TEXT NOT NULL,
"id" INTEGER PRIMARY KEY AUTOINCREMENT
);
CREATE TABLE "callbacks_users" (
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"before_s" TEXT NOT NULL,
"before_c" TEXT NOT NULL,
"before_u" TEXT NOT NULL,
"before_d" TEXT NOT NULL,
"after_s" TEXT NOT NULL,
"after_c" TEXT NOT NULL,
"after_u" TEXT NOT NULL,
"after_d" TEXT NOT NULL,
"id" INTEGER PRIMARY KEY AUTOINCREMENT
);
CREATE TABLE "books" (
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"title" TEXT NOT NULL,
"user_id" char(36) NOT NULL,
"isbn" TEXT NOT NULL,
"id" INTEGER PRIMARY KEY AUTOINCREMENT
);
