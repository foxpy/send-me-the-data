-- +goose Up
ALTER TABLE smtd.links ADD max_file_size BIGINT;
UPDATE smtd.links SET max_file_size = 4_294_967_296; -- 4 GiB
ALTER TABLE smtd.links ALTER max_file_size SET NOT NULL;

ALTER TABLE smtd.links
    ADD CONSTRAINT max_file_size CHECK
        (max_file_size > 0);

-- +goose Down
ALTER TABLE smtd.links DROP CONSTRAINT max_file_size;

ALTER TABLE smtd.links DROP max_file_size;
