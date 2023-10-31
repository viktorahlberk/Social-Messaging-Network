-- +migrate Up
CREATE TABLE IF NOT EXISTS almost_private (
    "user_id" VARCHAR(255) not null,
    "post_id" VARCHAR(255) not null
);

-- +migrate Down
DROP TABLE almost_private;