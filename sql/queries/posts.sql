-- name: AddPost :one
INSERT INTO posts (
    id,
    created_at,
    updated_at,
    title,
    url,
    description,
    published_at,
    feed_id
)
VALUES (
    sqlc.arg('id'),
    sqlc.arg('created_at'),
    sqlc.arg('created_at'),
    sqlc.arg('title'),
    sqlc.arg('url'),
    sqlc.arg('description'),
    sqlc.arg('published_at'),
    sqlc.arg('feed_id')
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT p.*, f.name AS feed_name
FROM posts p
INNER JOIN feeds f
ON p.feed_id = f.id
INNER JOIN feed_follows ff
ON f.id = ff.feed_id
WHERE ff.user_id = sqlc.arg('user_id')
ORDER BY COALESCE(p.published_at, p.created_at) DESC
LIMIT sqlc.arg('post_amount');