BEGIN;

CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    link_id VARCHAR(64) NOT NULL,
    destination_url VARCHAR(255) NOT NULL,
    parameters JSONB NOT NULL,
    allow_parameters_override BOOLEAN NOT NULL
);

CREATE UNIQUE INDEX links_link_id_key ON links (link_id);

COMMIT;
