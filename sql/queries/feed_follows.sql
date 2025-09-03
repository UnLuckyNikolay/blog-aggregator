-- name: CreateFeedFollow :one
WITH new_feed_follow AS (
    INSERT INTO feed_follows (
        id,
        created_at,
        updated_at,
        user_id,
        feed_id
    )
    VALUES (
        $1,
        $2,
        $3,
        $4,
        (
            SELECT id
            FROM feeds
            WHERE feeds.url = sqlc.arg('feed_url')
        )
    )
    RETURNING * 
)
SELECT 
    new_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM new_feed_follow
INNER JOIN feeds
ON new_feed_follow.feed_id = feeds.id
INNER JOIN users
ON new_feed_follow.user_id = users.id;

-- name: DeleteFeedFollow :one
WITH del AS (
    DELETE
    FROM feed_follows ff
    USING feeds f
    WHERE ff.feed_id = f.id
        AND ff.user_id = sqlc.arg('user_id')
        AND f.url = sqlc.arg('feed_url')
    RETURNING ff.feed_id
)
SELECT f.*
FROM feeds f
JOIN del
ON del.feed_id = f.id;

-- name: GetFeedFollowsForUser :many
SELECT feeds.name
FROM feed_follows
INNER JOIN feeds
ON feed_follows.feed_id = feeds.id
INNER JOIN users
ON feed_follows.user_id = users.id
WHERE users.name = sqlc.arg('user_name');
