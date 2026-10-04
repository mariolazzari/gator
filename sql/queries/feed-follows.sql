-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows(id, user_id, feed_id) 
        values($1, $2, $3)
    RETURNING *
)
SELECT iff.*,
    f.name AS feed_name,
    u.name AS user_name
FROM inserted_feed_follow as iff
INNER JOIN users as u on u.id = iff.user_id
INNER JOIN feeds as f on f.id = iff.feed_id;

-- name: GetFeedFollowsForUser :many
select f.name as feed_name, u.name as user_name
from feed_follows as ff
inner join users as u on u.id = ff.user_id
inner join feeds as f on f.id = ff.feed_id
where u.name = $1;

-- name: DeleteFeedFollow :exec
delete from feed_follows
where feed_id = $1 and user_id = $2;