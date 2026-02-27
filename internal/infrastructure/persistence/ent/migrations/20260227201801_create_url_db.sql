-- Create "urls" table
CREATE TABLE "public"."urls" (
  "id" character varying(32) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  "deleted_at" timestamp NULL,
  "destination" text NOT NULL,
  "user_id" bigint NOT NULL,
  "expires_at" timestamp NULL,
  PRIMARY KEY ("id")
);
-- Create index "ix_urls_user_active" to table: "urls"
CREATE INDEX "ix_urls_user_active" ON "public"."urls" ("user_id", "updated_at") WHERE (deleted_at IS NULL);
