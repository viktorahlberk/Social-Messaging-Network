-- +migrate Up
CREATE TABLE IF NOT EXISTS group_messages (
			"message_id" VARCHAR(255) not null,
			"receiver_id" VARCHAR(255) not null,
			"is_read" INT default 0
);

-- +migrate Down
DROP TABLE group_messages;
