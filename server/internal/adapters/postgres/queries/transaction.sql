-- name: Withdraw :exec
UPDATE "accounts"
SET balance = balance - sqlc.arg(amount)
WHERE id = sqlc.arg(id);

-- name: Deposit :exec
UPDATE "accounts"
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id);

-- name: SaveTransaction :exec
INSERT INTO "transactions" (id, created_at, from_account_id, to_account_id, amount, type)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetTransactions :many
SELECT tx.id, tx.amount, tx.type, tx.created_at,
       fr.id AS from_id, fr.title AS from_title, fr.description AS from_description, fr.balance AS from_balance, fr.currency AS from_currency, fr.is_blocked AS from_is_blocked, fr.user_id AS from_user_id,
       t.id AS to_id, t.title AS to_title, t.description AS to_description, t.balance AS to_balance, t.currency AS to_currency, t.is_blocked AS to_is_blocked, t.user_id AS to_user_id
from (SELECT *
      from "transactions"
      WHERE from_account_id = $1
         OR to_account_id = $1
      ORDER BY created_at DESC
      LIMIT sqlc.narg('limit')) tx
         LEFT JOIN "accounts" fr ON tx.from_account_id = fr.id
         LEFT JOIN "accounts" t ON tx.to_account_id = t.id;