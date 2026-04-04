-- +goose Up
ALTER TABLE smtd.links ADD user_downloadable BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE smtd.links DROP user_downloadable;
