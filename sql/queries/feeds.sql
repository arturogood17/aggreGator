-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: ListFeeds :many
SELECT name, url, user_id FROM feeds;

-- name: FeedByURL :one
SELECT * FROM feeds
WHERE url = $1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET updated_at = NOW(), last_fetched_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * from feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;