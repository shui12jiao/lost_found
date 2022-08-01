-- name: AddUsr :one
INSERT INTO usr (
    openid,
    name,
    student_id,
    avatar
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUsr :one
SELECT * FROM usr
WHERE openid  = $1 LIMIT 1;


