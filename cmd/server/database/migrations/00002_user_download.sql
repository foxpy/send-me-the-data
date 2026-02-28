-- +goose Up
-- +goose StatementBegin
ALTER TABLE smtd.links ADD user_downloadable BOOLEAN NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE smtd.links DROP user_downloadable;
-- +goose StatementEnd
