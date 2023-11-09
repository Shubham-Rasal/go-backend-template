CREATE TABLE "follows" (
  "following_user_id" integer,
  "followed_user_id" integer,
  "created_at" timestamp
);

CREATE TABLE "users" (
  "id" bigserial NOT NULL UNIQUE,
  "username" varchar PRIMARY KEY,
  "password" varchar NOT NULL,
  "email" varchar NOT NULL
);

CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "user_id" integer,
  "username" varchar,
  "role" varchar,
  "reputation" integer,
  "created_at" timestamp
);

CREATE TABLE "posts" (
  "id" integer PRIMARY KEY,
  "title" varchar,
  "body" text,
  "user_id" integer,
  "status" varchar,
  "created_at" timestamp
);

COMMENT ON COLUMN "posts"."body" IS 'Content of the post';

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "accounts" ("user_id");

ALTER TABLE "follows" ADD FOREIGN KEY ("following_user_id") REFERENCES "accounts" ("user_id");

ALTER TABLE "follows" ADD FOREIGN KEY ("followed_user_id") REFERENCES "accounts" ("user_id");

ALTER TABLE "users" ADD FOREIGN KEY ("id") REFERENCES "accounts" ("user_id");
