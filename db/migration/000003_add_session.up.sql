CREATE TABLE "sessions" (
    "id" uuid NOT NULL,
    "username" varchar NOT NULL,
    "refresh_token" varchar NOT NULL,
    "is_blocked" boolean NOT NULL DEFAULT false,
    "user_agent" varchar NOT NULL,
    "client_ip" varchar NOT NULL,
    "expires_at" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);
ALTER TABLE "sessiions"
ADD FOREIGN KEY ("username") REFERENCES "users" ("username");