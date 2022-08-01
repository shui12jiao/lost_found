/*
type_wide
*/
-- name: AddTypeWide :one
INSERT INTO type_wide (
    name
) VALUES (
  $1
) RETURNING *;

-- name: GetTypeWide :one
SELECT * FROM type_wide
WHERE id  = $1 LIMIT 1;

-- name: ListTypeWide :many
SELECT * FROM type_wide;

-- name: DeleteTypeWide :exec
DELETE FROM type_wide
WHERE id = $1;


/*
type_narrow
*/
-- name: AddTypeNarrow :one
INSERT INTO type_narrow (
    name,
    wide_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetTypeNarrow :one
SELECT * FROM type_narrow
WHERE id  = $1 LIMIT 1;

-- name: ListTypeNarrow :many
SELECT * FROM type_narrow;

-- name: DeleteTypeNarrow :exec
DELETE FROM type_narrow
WHERE id = $1;