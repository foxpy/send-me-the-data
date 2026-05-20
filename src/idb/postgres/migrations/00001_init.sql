-- +goose Up
CREATE SCHEMA smtd;
COMMENT ON SCHEMA smtd IS 'send me the data acronym';

CREATE TABLE smtd.links (
    link_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    external_key TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE smtd.file_journal (
    link_id BIGINT NOT NULL REFERENCES smtd.links,
    name TEXT NOT NULL,
    PRIMARY KEY (link_id, name)
);

-- +goose Down
DROP TABLE smtd.file_journal;
DROP TABLE smtd.links;
DROP SCHEMA smtd;
