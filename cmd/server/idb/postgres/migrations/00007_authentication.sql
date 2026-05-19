-- +goose Up
CREATE TABLE smtd.admins (
    admin_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash BYTEA NOT NULL
);

ALTER TABLE smtd.admins
    ADD CONSTRAINT username_length CHECK
        (char_length(username) > 0 AND char_length(username) < 256);

-- XXX: anyone with access to database can steal and reuse session tokens
CREATE TABLE smtd.session_tokens (
    session_token TEXT PRIMARY KEY,
    expires_at TIMESTAMPTZ NOT NULL,
    admin_id BIGINT NOT NULL REFERENCES smtd.admins
);

CREATE INDEX session_token_expiration ON smtd.session_tokens (expires_at);

-- +goose Down
DROP INDEX smtd.session_token_expiration;
DROP TABLE smtd.session_tokens;
ALTER TABLE smtd.admins DROP CONSTRAINT username_length;
DROP TABLE smtd.admins;
