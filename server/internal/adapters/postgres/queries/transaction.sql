-- name: Withdraw :exec
UPDATE "accounts" SET balance = balance - sqlc.arg(amount) WHERE id = sqlc.arg(id);

-- name: Deposit :exec
UPDATE "accounts" SET balance = balance + sqlc.arg(amount) WHERE id = sqlc.arg(id);

-- name: SaveTransaction :exec
INSERT INTO "transactions" (id, created_at, from_account_id, to_account_id, amount, type) VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetTransactions :many
SELECT * from "transactions" WHERE from_account_id=$1 OR to_account_id=$1 ORDER BY created_at DESC;