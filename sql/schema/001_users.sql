-- +goose Up
CREATE TABLE users(
    id UUID PRIMARY KEY,
    name text unique not null,
    created_at TIMESTAMP NOT NULL default now(),
    updated_at TIMESTAMP NOT NULL default now()

)

-- +goose Down
DROP TABLE users;
