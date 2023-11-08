ALTER TABLE "follows" DROP CONSTRAINT "follows_followed_user_id_fkey";

ALTER TABLE "follows" DROP CONSTRAINT "follows_following_user_id_fkey";

ALTER TABLE "posts" DROP CONSTRAINT "posts_user_id_fkey";

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "accounts" ("id");

ALTER TABLE "follows" ADD FOREIGN KEY ("followed_user_id") REFERENCES "accounts" ("id");

ALTER TABLE "follows" ADD FOREIGN KEY ("following_user_id") REFERENCES "accounts" ("id");

--  SELECT
--             *
--         FROM
--             information_schema.table_constraints AS tc
--         WHERE tc.table_name = 'follows'
--             AND tc.constraint_type = 'FOREIGN KEY'
--             AND tc.table_schema = 'public';