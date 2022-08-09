-- name: AddUsr :one
INSERT INTO usr (
    openid,
    name,
    phone,
    student_id,
    avatar_url,
    avatar
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUsr :one
SELECT * FROM usr
WHERE openid  = $1 LIMIT 1;

-- name: DeleteUsr :exec
DELETE FROM usr
WHERE openid = $1;

-- name: UpdateUsrName :one
UPDATE usr
SET name = $2

WHERE openid = $1
RETURNING *;

-- name: UpdateUsrAvatar :one
UPDATE usr
SET avatar_url = $2,
    avatar = $3
WHERE openid = $1
RETURNING *;

-- name: UpdateUsr :one
UPDATE usr
SET name = $2,
    phone = $3,
    student_id = $4,
    avatar_url = $5,
    avatar = $6
WHERE openid = $1
RETURNING *; 