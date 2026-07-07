-- name: CreatePostLike :one
with inserted_like as (
    insert into post_likes (id, created_at, updated_at, user_id, post_id)
    values ($1, $2, $3, $4, $5)
    returning *
)
select inserted_like.*,
       posts.title as post_title,
       posts.url as post_url
from inserted_like
inner join posts on inserted_like.post_id = posts.id;

-- name: DeletePostLike :exec
delete from post_likes
where user_id = $1 and post_id = $2;

-- name: GetLikedPostsForUser :many
select posts.*
from post_likes
inner join posts on post_likes.post_id = posts.id
where post_likes.user_id = $1
order by post_likes.created_at desc;