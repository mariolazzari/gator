-- name: CreateFeed :one
INSERT INTO feeds (id, user_id, name, url, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT f.name, f.url, u.name as user
FROM feeds f
inner join users as u on u.id = f.user_id;

-- name: GetFeedByUrl :one
SELECT *
FROM feeds 
where url = $1;
