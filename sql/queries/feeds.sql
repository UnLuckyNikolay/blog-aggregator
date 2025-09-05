-- name: AddFeed :one
INSERT INTO feeds(
    id,
    created_at,
    updated_at,
    name,
    url,
    user_id
)
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
WITH 
    follow_count AS (
        SELECT 
            feed_id,
            COUNT(user_id) AS follower_count
        FROM feed_follows
        GROUP BY feed_id
),
    post_count AS (
        SELECT 
            feed_id,
            COUNT(id) AS post_count
        FROM posts
        GROUP BY feed_id
)
SELECT 
    f.id, 
    f.name, 
    f.url, 
    f.user_id, 
    u.name AS user_name,
    fc.follower_count,
    pc.post_count
FROM feeds f
INNER JOIN users u
ON f.user_id = u.id
LEFT JOIN follow_count fc
ON f.id = fc.feed_id
LEFT JOIN post_count pc
ON f.id = pc.feed_id;

-- name: GetFeedsToFetch :many
SELECT
    feeds.id,
    feeds.name,
    feeds.url
FROM feeds;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET 
    updated_at = sqlc.arg('time'),
    last_fetched_at = sqlc.arg('time')
WHERE id = sqlc.arg('feed_id');