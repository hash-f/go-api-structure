-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    password_hash,
    api_key
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: UpdateUserAPIKey :one
UPDATE users
SET api_key = $1,
    updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: GetUserByAPIKey :one
SELECT * FROM users
WHERE api_key = $1;
