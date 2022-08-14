/*
location_wide
*/
-- name: AddLocationWide :one
INSERT INTO location_wide (
    name,
    campus
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetLocationWide :one
SELECT * FROM location_wide
WHERE id = $1 LIMIT 1;

-- name: ListLocationWide :many
SELECT * FROM location_wide;

-- name: DeleteLocationWide :exec
DELETE FROM location_wide
WHERE id = $1;


/*
location_narrow
*/
-- name: AddLocationNarrow :one
INSERT INTO location_narrow (
    name,
    wide_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetLocationNarrow :one
SELECT * FROM location_narrow
WHERE id  = $1 LIMIT 1;

-- name: ListLocationNarrow :many
SELECT * FROM location_narrow
ORDER BY id;

-- name: ListLocationNarrowByWide :many
SELECT * FROM location_narrow
WHERE wide_id = $1
ORDER BY id;

-- name: DeleteLocationNarrow :exec
DELETE FROM location_narrow
WHERE id = $1;