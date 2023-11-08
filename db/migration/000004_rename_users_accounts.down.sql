
ALTER TABLE "accounts" DROP CONSTRAINT "accounts_user_id_fkey";
DROP TABLE IF EXISTS "users";
ALTER TABLE "accounts" DROP COLUMN "user_id";
ALTER TABLE "accounts" RENAME TO "users";

