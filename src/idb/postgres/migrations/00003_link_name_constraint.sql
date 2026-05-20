-- +goose Up
ALTER TABLE smtd.links
    ADD CONSTRAINT link_name_length CHECK
        (char_length(name) > 0 AND char_length(name) < 256);

-- +goose Down
ALTER TABLE smtd.links DROP CONSTRAINT link_name_length;
