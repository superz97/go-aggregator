-- +goose Up
create table post_likes(
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    user_id uuid not null references users(id) on delete cascade,
    post_id uuid not null references posts(id) on delete cascade,
    unique(user_id, post_id)
);

-- +goose Down
drop table post_likes;