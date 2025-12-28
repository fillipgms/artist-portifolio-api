-- name: ListIllustrations :many
SELECT * 
FROM illustrations 
ORDER BY finished_at DESC, created_at DESC
LIMIT $1 OFFSET $2;

-- name: FindIllustrationById :one
SELECT * FROM illustrations WHERE id = $1;

-- name: FindIllustrationByName :one
SELECT * FROM illustrations WHERE slug = $1 LIMIT 1;

-- name: CreateIllustration :one
INSERT INTO illustrations (
    title,
    slug,
    description,
    imageURL,
    imageHeight,
    imageWidth,
    imageMimeType,
    imageFileSize,
    post,
    finished_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;

-- name: UpdateSlug :one
UPDATE illustrations SET slug = $1 WHERE id = $2 RETURNING *; 

-- name: FindIllustrationsCount :one
SELECT COUNT(*) FROM illustrations;