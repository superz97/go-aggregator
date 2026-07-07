-- name: CreatePost :one
insert into posts(id, created_at, updated_at, title, url, description, published_at, feed_id)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning *;

-- name: GetPostsForUser :many
select posts.*
from posts
inner join feed_follows on posts.feed_id = feed_follows.feed_id
where feed_follows.user_id = $1
order by posts.published_at desc nulls last
limit $2;

-- name: GetPostsForUserByPublishedAtDesc :many
select posts.*
from posts
inner join feed_follows on posts.feed_id = feed_follows.feed_id
inner join feeds on feed_follows.feed_id = feeds.id
where feed_follows.user_id = $1
and (sqlc.narg('feed_name')::text is null or feeds.name = sqlc.narg('feed_name'))
order by posts.published_at desc nulls last
limit $2
offset $3;

-- name: GetPostsForUserByPublishedAtAsc :many
select posts.*
from posts
inner join feed_follows on posts.feed_id = feed_follows.feed_id
inner join feeds on feed_follows.feed_id = feeds.id
where feed_follows.user_id = $1
and (sqlc.narg('feed_name')::text is null or feeds.name = sqlc.narg('feed_name'))
order by posts.published_at asc nulls last
limit $2
offset $3;

-- name: GetPostsForUserByTitleDesc :many
select posts.*
from posts
inner join feed_follows on posts.feed_id = feed_follows.feed_id
inner join feeds on feed_follows.feed_id = feeds.id
where feed_follows.user_id = $1
and (sqlc.narg('feed_name')::text is null or feeds.name = sqlc.narg('feed_name'))
order by posts.title desc
limit $2
offset $3;

-- name: GetPostsForUserByTitleAsc :many
select posts.*
from posts
inner join feed_follows on posts.feed_id = feed_follows.feed_id
inner join feeds on feed_follows.feed_id = feeds.id
where feed_follows.user_id = $1
and (sqlc.narg('feed_name')::text is null or feeds.name = sqlc.narg('feed_name'))
order by posts.title asc
limit $2
offset $3;

-- name: SearchPostsForUser :many
select posts.*
from posts
inner join feed_follows on posts.feed_id = feed_follows.feed_id
where feed_follows.user_id = $1
and similarity(posts.title, $2) > 0.1
order by similarity(posts.title, $2) desc
limit $3;

-- name: GetPostByURL :one
select * from posts where url = $1;

-- name: GetPostsForTUI :many
select posts.*, feeds.name as feed_name
from posts
inner join feed_follows on posts.feed_id = feed_follows.feed_id
inner join feeds on feed_follows.feed_id = feeds.id
where feed_follows.user_id = $1
order by posts.published_at desc nulls last
limit $2;