CREATE TABLE IF NOT EXISTS "user" (
    "id" bigserial PRIMARY KEY,
    "username" text UNIQUE NOT NULL,
    "hashed_password" text NOT NULL,
    "email" text UNIQUE NOT NULL,
    "verified" boolean NOT NULL DEFAULT false
);