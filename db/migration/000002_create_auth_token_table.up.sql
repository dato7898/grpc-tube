CREATE TABLE IF NOT EXISTS "auth_token" (
    "id" bigserial PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "code" text NOT NULL,
    "expired_at" timestamp NOT NULL DEFAULT (now() + INTERVAL '1 hour')
);

ALTER TABLE "auth_token" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");
