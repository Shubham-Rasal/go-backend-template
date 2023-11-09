-- add constrants fk between account user_id and others
ALTER TABLE "follows" ADD FOREIGN KEY ("followed_user_id") REFERENCES "accounts" ("user_id");

ALTER TABLE "follows" ADD FOREIGN KEY ("following_user_id") REFERENCES "accounts" ("user_id");

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "accounts" ("user_id");