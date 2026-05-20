-- +goose Up
ALTER TABLE smtd.links RENAME external_key TO public_id;

-- +goose Down
ALTER TABLE smtd.links RENAME public_id TO external_key;
