-- name: CreateUser :one
INSERT INTO users (username, role) VALUES ($1, $2) RETURNING *;


-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id 
LIMIT $1
OFFSET $2;

-- name: UpdateUser :exec
UPDATE users SET username = $1, role = $2
WHERE id = $3;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;