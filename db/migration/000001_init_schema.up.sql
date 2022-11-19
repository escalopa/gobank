CREATE TYPE "currency" AS ENUM (
  'EGP',
  'RUB',
  'USD'
);

CREATE TABLE IF NOT EXISTS "users" (
                                       "id" bigserial PRIMARY KEY NOT NULL,
                                       "owner" varchar NOT NULL,
                                       "username" varchar UNIQUE NOT NULL,
                                       "salt" varchar NOT NULL,
                                       "password" varchar NOT NULL,
                                       "created_at" timestamptz NOT NULL DEFAULT (now())
    );

CREATE TABLE IF NOT EXISTS "accounts" (
                                          "id" bigserial PRIMARY KEY NOT NULL ,
                                          "user_id" bigint NOT NULL,
                                          "balance" double PRECISION NOT NULL,
                                          "currency" currency NOT NULL,
                                          "created_at" timestamptz NOT NULL DEFAULT (now())
    );

CREATE TABLE IF NOT EXISTS "entries" (
                                         "id" bigserial PRIMARY KEY NOT NULL,
                                         "account_id" bigint NOT NULL,
                                         "amount" double PRECISION NOT NULL,
                                         "created_at" timestamptz NOT NULL DEFAULT (now())
    );

CREATE TABLE IF NOT EXISTS "transfers" (
                                           "id" bigserial PRIMARY KEY NOT NULL ,
                                           "from_account_id" bigint NOT NULL ,
                                           "to_account_id" bigint NOT NULL ,
                                           "amount" bigint NOT NULL,
                                           "reverted" bool DEFAULT false,
                                           "created_at" timestamptz DEFAULT (now())
    );

--- INDEXES

CREATE INDEX ON "accounts" ("user_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

--- FOREIGN KEYS

ALTER TABLE "accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

--- CONSTRAINTS

ALTER TABLE "transfers" ADD CONSTRAINT "from_account_id_to_account_id" CHECK ("from_account_id" != "to_account_id");

ALTER TABLE "transfers" ADD  CONSTRAINT "transfers_amount_positive" CHECK ("amount" > 0);

ALTER TABLE "accounts" ADD  CONSTRAINT "unique_currency_per_user" UNIQUE ("user_id", "currency");

-- COMMENTS

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';
