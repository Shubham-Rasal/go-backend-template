-- name: CreateAccount :one
INSERT INTO Accounts (username, role , user_id) VALUES ($1, $2, $3) RETURNING *;


-- name: GetAccount :one
SELECT * FROM Accounts
WHERE user_id = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM Accounts
ORDER BY id 
LIMIT $1
OFFSET $2;

-- name: UpdateReputation :exec
UPDATE Accounts SET reputation = reputation + $1
WHERE user_id = $2;

-- name: DeleteAccount :exec
DELETE FROM Accounts
WHERE user_id = $1 RETURNING user_id;