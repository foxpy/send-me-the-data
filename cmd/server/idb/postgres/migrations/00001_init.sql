-- +goose Up
CREATE SCHEMA smtd;
COMMENT ON SCHEMA smtd IS 'send me the data acronym';

CREATE TABLE smtd.links (
    link_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    -- XXX: maybe this field should actually be the primary key and `link_id` should be dropped since it's only used in the database?
    -- food for thought: external ID can be encoded into a 72-bit number
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
