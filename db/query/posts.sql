-- name: CreatePost :one
INSERT INTO "posts" ("title", "body", "user_id", "status") 
VALUES ($1, $2, $3, $4) RETURNING *;


-- name: GetPost :one
SELECT * FROM "posts"
WHERE id = $1 LIMIT 1;

-- name: ListPosts :many
SELECT * FROM "posts"
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: LikePost :exec
UPDATE "posts" SET "likes" = "likes" + 1
WHERE id = $1 RETURNING *;