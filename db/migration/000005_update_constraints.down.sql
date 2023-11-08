-- -- find the foriegn key constraint name between follows and accounts and drop it
-- ALTER TABLE "follows" DROP CONSTRAINT "follows_following_user_id_fkey";
-- -- find the foriegn key constraint name between follows and accounts and drop it
-- ALTER TABLE "follows" DROP CONSTRAINT "follows_followed_user_id_fkey";

-- reverse this 

ALTER TABLE "follows" ADD FOREIGN KEY ("following_user_id") REFERENCES "accounts" ("id");

ALTER TABLE "follows" ADD FOREIGN KEY ("followed_user_id") REFERENCES "accounts" ("id");