ALTER TABLE "posts" DROP CONSTRAINT "posts_user_id_fkey";

ALTER TABLE "follows" DROP CONSTRAINT "follows_followed_user_id_fkey";

ALTER TABLE "follows" DROP CONSTRAINT "follows_following_user_id_fkey";

