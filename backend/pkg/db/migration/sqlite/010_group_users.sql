-- +migrate Up
CREATE TABLE IF NOT EXISTS group_users (
    "group_id" VARCHAR(255) not null,
    "user_id" VARCHAR(255) not null
);

-- +migrate Down
DROP TABLE group_users;