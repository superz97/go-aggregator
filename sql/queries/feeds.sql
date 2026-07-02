-- name: CreateFeed :one
insert into feeds (id, created_at, updated_at, name, url, user_id)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: ListFeeds :many
select feeds.name, feeds.url, users.name as user_name
from feeds
join users on feeds.user_id = users.id;

-- name: GetFeedByURL :one
select * from feeds where url = $1;

-- name: CreateFeedFollow :one
with inserted_feed_follow as (
    insert into feed_follows (id, created_at, updated_at, user_id, feed_id)
    values ($1, $2, $3, $4, $5)
    returning *
)
select inserted_feed_follow.*,
       feeds.name as feed_name,
       users.name as user_name
from inserted_feed_follow
inner join feeds on inserted_feed_follow.feed_id = feeds.id
inner join users on inserted_feed_follow.user_id = users.id;

-- name: GetFeedFollowsForUser :many
select feed_follows.*, feeds.name as feed_name, users.name as user_name
from feed_follows
inner join feeds on feed_follows.feed_id = feeds.id
inner join users on feed_follows.user_id = users.id
where feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
delete from feed_follows
where user_id = $1 and feed_id = $2;

-- name: MarkFeedFetched :exec
update feeds
set last_fetched_at = now(), updated_at = now()
where id = $1;

-- name: GetNextFeedsToFetch :many
select * from feeds
order by last_fetched_at nulls first
limit $1;