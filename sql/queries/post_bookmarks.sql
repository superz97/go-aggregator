-- name: CreatePostBookmark :one
with inserted_bookmark as (
    insert into post_bookmarks (id, created_at, updated_at, user_id, post_id)
    values ($1, $2, $3, $4, $5)
    returning *
)
select inserted_bookmark.*,
       posts.title as post_title,
       posts.url as post_url
from inserted_bookmark
inner join posts on inserted_bookmark.post_id = posts.id;

-- name: DeletePostBookmark :exec
delete from post_bookmarks
where user_id = $1 and post_id = $2;

-- name: GetBookmarkedPostsForUser :many
select posts.*
from post_bookmarks
inner join posts on post_bookmarks.post_id = posts.id
where post_bookmarks.user_id = $1
order by post_bookmarks.created_at desc;