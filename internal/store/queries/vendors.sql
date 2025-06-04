-- name: CreateVendor :one
INSERT INTO vendors (
    user_id,
    name,
    description
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetVendorByID :one
SELECT * FROM vendors
WHERE id = $1;

-- name: ListVendorsByUserID :many
SELECT * FROM vendors
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateVendor :one
UPDATE vendors
SET
    name = COALESCE(sqlc.narg(name), name),
    description = COALESCE(sqlc.narg(description), description),
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteVendor :exec
DELETE FROM vendors
WHERE id = $1 AND user_id = $2;
