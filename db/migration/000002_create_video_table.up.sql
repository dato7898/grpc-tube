CREATE TABLE IF NOT EXISTS "video" (
    "id" text PRIMARY KEY,
    "title" text,
    "description" text,
    "views" bigint DEFAULT 0
);