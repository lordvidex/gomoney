-- name: CreateAccount :one
INSERT INTO "accounts" (title, description, currency, user_id) VALUES ($1, $2, $3, $4) RETURNING id;

-- name: GetUserAccounts :many
SELECT * FROM "accounts" WHERE user_id = $1;

-- name: GetAccount :one
SELECT * FROM "accounts" WHERE id = $1;

-- name: DeleteAccount :exec
DELETE FROM "accounts" WHERE id = $1;