-- name: CreateUser :one
INSERT INTO "users" (name, phone) VALUES ($1, $2) RETURNING *;

-- name: GetUserByPhone :one
SELECT * FROM "users" WHERE phone = $1;
