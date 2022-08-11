/*
found
*/
-- name: AddFound :one
INSERT INTO found (
    picker_openid,
    found_date,
    time_bucket,
    location_id,
    location_info,
    location_status,
    type_id,
    item_info,
    image,
    image_key,
    owner_info,
    addtional_info
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING *;

-- name: GetFound :one
SELECT * FROM found
WHERE id = $1 LIMIT 1;

-- name: ListFound :many
SELECT * FROM found
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListFoundByPicker :many
SELECT * FROM found
WHERE picker_openid = $1
ORDER BY id;

-- name: DeleteFound :exec
DELETE FROM found
WHERE id = $1;


/*
lost
*/
-- name: AddLost :one
INSERT INTO lost (
    owner_openid,
    lost_date,
    time_bucket,
    type_id,
    location_id,
    location_id1,
    location_id2
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetLost :one
SELECT * FROM lost
WHERE id = $1 LIMIT 1;

-- name: ListLost :many
SELECT * FROM lost
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListLostByOwner :many
SELECT * FROM lost
WHERE owner_openid = $1
ORDER BY id;

-- name: DeleteLost :exec
DELETE FROM lost
WHERE id = $1;


/*
match
*/
-- name: AddMatch :one
INSERT INTO match (
    picker_openid,
    owner_openid,
    found_date,
    lost_date,
    type_id,
    item_info,
    image,
    image_key,
    comment
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetMatch :one
SELECT * FROM match
WHERE id = $1 LIMIT 1;

-- name: ListMatch :many
SELECT * FROM match
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListMatchByPicker :many
SELECT * FROM match
WHERE picker_openid = $1
ORDER BY id;

-- name: ListMatchByOwner :many
SELECT * FROM match
WHERE owner_openid = $1
ORDER BY id;

-- name: DeleteMatch :exec
DELETE FROM match
WHERE id = $1;