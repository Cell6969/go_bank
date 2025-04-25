CREATE TABLE "sessions" (
    "id" uuid PRIMARY KEY,
    "username" VARCHAR NOT NULL,
    "refresh_token" VARCHAR NOT NULL,
    "user_agent" VARCHAR NOT NULL,
    "client_ip" VARCHAR NOT NULL,
    "is_blocked" BOOLEAN NOT NULL,
    "expired_at" timestamp NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");