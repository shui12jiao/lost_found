-- name: AddManager :one
INSERT INTO manager (
    usr_openid,
    permission
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetManager :one
SELECT * FROM manager
WHERE id  = $1 LIMIT 1;

-- name: GetManagerByOpenid :one
SELECT * FROM manager
WHERE usr_openid = $1 LIMIT 1;

-- name: ListManager :many
SELECT * FROM manager
ORDER BY id;

-- name: DeleteManager :exec
DELETE FROM manager
WHERE id = $1;

-- name: UpdateManager :one
UPDATE manager
SET permission = $2
WHERE id = $1
RETURNING *;