-- +goose Up
CREATE TABLE feeds(
    id UUID PRIMARY KEY,
    user_id UUID not null references users(id) on delete cascade,
    name text not null,
    url text unique not null,
    created_at TIMESTAMP NOT NULL default now(),
    updated_at TIMESTAMP NOT NULL default now()
);

-- +goose Down
DROP TABLE feeds;
