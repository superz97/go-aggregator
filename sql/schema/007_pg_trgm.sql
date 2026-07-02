-- +goose Up
create extension if not exists pg_trgm;

-- +goose Down
drop extension if exists pg_trgm;