ALTER TABLE "posts" DROP CONSTRAINT "posts_user_id_fkey";

ALTER TABLE "follows" DROP CONSTRAINT "follows_followed_user_id_fkey";

ALTER TABLE "follows" DROP CONSTRAINT "follows_following_user_id_fkey";

SELECT
            *
        FROM
            information_schema.table_constraints AS tc
        WHERE tc.table_name = 'posts'
            AND tc.constraint_type = 'FOREIGN KEY'
            AND tc.table_schema = 'public';