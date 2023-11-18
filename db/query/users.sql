-- name: CreateUser :one
INSERT INTO users (username, password , email) VALUES ($1, $2, $3) RETURNING *;


-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id 
LIMIT $1
OFFSET $2;

-- name: UpdatePassword :exec
UPDATE users SET password = $1
WHERE id = $2;

-- name: UpdateEmail :exec
UPDATE users SET email = $1
WHERE id = $2
RETURNING id;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1 RETURNING id;