-- +goose Up
ALTER TABLE smtd.links ADD upload_enabled BOOLEAN NOT NULL DEFAULT true;

-- +goose Down
ALTER TABLE smtd.links DROP upload_enabled;
