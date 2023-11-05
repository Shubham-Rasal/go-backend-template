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

-- name: UpdateReputation :exec
UPDATE users SET reputation = reputation + $1
WHERE id = $2;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;