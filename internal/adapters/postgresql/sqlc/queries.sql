-- name: ListIllustrations :many
SELECT * 
FROM illustrations 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: FindIllustrationById :one
SELECT * FROM illustrations WHERE id = $1;

-- name: CreateIllustration :one
INSERT INTO illustrations (
    title,
    description,
    imageURL,
    post,
    finished_at
) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: FindIllustrationsCount :one
SELECT COUNT(*) FROM illustrations;