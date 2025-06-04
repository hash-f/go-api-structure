-- name: CreateMerchant :one
INSERT INTO merchants (
    user_id,
    name,
    description
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetMerchantByID :one
SELECT * FROM merchants
WHERE id = $1;

-- name: ListMerchantsByUserID :many
SELECT * FROM merchants
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateMerchant :one
UPDATE merchants
SET
    name = COALESCE(sqlc.narg(name), name),
    description = COALESCE(sqlc.narg(description), description),
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteMerchant :exec
DELETE FROM merchants
WHERE id = $1 AND user_id = $2;
