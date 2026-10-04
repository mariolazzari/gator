-- +goose Up
CREATE TABLE feed_follows(
    id UUID PRIMARY KEY,
    user_id UUID not null references users(id) on delete cascade,
    feed_id UUID not null references feeds(id) on delete cascade,    
    created_at TIMESTAMP NOT NULL default now(),
    updated_at TIMESTAMP NOT NULL default now(),
    CONSTRAINT feed_follows_user_feed_unique UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;
