-- find the foriegn key constraint name between follows and accounts and drop it
ALTER TABLE "follows" DROP CONSTRAINT "follows_followed_user_id_fkey";
-- find the foriegn key constraint name between follows and accounts and drop it
ALTER TABLE "follows" DROP CONSTRAINT "follows_following_user_id_fkey";

ALTER TABLE "posts" DROP CONSTRAINT "posts_user_id_fkey";

-- add constrants fk between account user_id and others
ALTER TABLE "follows" ADD FOREIGN KEY ("followed_user_id") REFERENCES "accounts" ("user_id");

ALTER TABLE "follows" ADD FOREIGN KEY ("following_user_id") REFERENCES "accounts" ("user_id");

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "accounts" ("user_id");