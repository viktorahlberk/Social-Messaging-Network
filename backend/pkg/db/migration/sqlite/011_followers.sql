-- +migrate Up
CREATE TABLE IF NOT EXISTS followers (
    "follower_id" VARCHAR(255) not null,
    "user_id" VARCHAR(255) not null
);

-- +migrate Down
DROP TABLE followers;