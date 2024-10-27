CREATE TABLE IF NOT EXISTS "user" (
    "id" bigserial PRIMARY KEY,
    "username" text NOT NULL,
    "hashed_password" text NOT NULL,
    "email" text
);