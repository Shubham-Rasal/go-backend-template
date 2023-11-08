ALTER TABLE "users" RENAME TO "accounts";
ALTER TABLE "accounts" ADD COLUMN "user_id" integer unique;

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "password" varchar NOT NULL,
  "email" varchar NOT NULL
);

-- -- add a constraint into account fk with user_id
ALTER TABLE "accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");


