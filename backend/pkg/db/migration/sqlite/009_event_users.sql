-- +migrate Up
CREATE TABLE IF NOT EXISTS event_users (
    "event_id" VARCHAR(255) not null,
    "user_id" VARCHAR(255) not null
);

-- +migrate Down
DROP TABLE event_users;