ALTER TABLE "users" RENAME TO "accounts";

-- create a new table users with the following columns
-- id, username, password, email
-- id is a serial primary key
-- username is a varchar not null
-- password is a varchar not null
-- email is a varchar not null
-- ALTER TABLE "accounts" ADD COLUMN "user_id" integer;

-- CREATE TABLE "users" (
--   "id" bigserial PRIMARY KEY,
--   "username" varchar NOT NULL,
--   "password" varchar NOT NULL,
--   "email" varchar NOT NULL
-- );

-- -- add a constraint into account fk with user_id
-- ALTER TABLE "accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");


