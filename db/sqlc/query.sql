-- name: GetPublicFiles :many
SELECT
  *
FROM
  file
WHERE
  public = true
ORDER BY
  updated_at DESC;

-- name: GetPrivateFiles :many
SELECT
  *
FROM
  file
WHERE
  user_id = $1
ORDER BY
  updated_at DESC;

-- name: CreateFile :one
INSERT INTO file (
  file_name,
  file_url,
  user_id,
  public
) VALUES (
  $1,
  $2,
  $3,
  $4
)
RETURNING *;

-- name: UpdateFile :one
UPDATE file
SET
  file_name = $2,
  public = $3,
  updated_at = CURRENT_TIMESTAMP
WHERE
  file_id = $1
RETURNING *;

-- name: DeleteFile :one
DELETE FROM file
WHERE
  file_id = $1
RETURNING *;
