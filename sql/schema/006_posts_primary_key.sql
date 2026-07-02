-- +goose Up
alter table posts add primary key (id);

-- +goose Down
alter table posts drop constraint posts_pkey;