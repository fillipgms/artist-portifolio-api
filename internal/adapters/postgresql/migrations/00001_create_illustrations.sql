-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS illustrations (
    id              BIGSERIAL PRIMARY key,
    title           TEXT NOT NULL,
    slug            TEXT UNIQUE,
    description     TEXT NOT NULL,
    imageURL        TEXT NOT NULL,
    imageHeight     INT,
    imageWidth      INT,
    imageMimeType   TEXT,
    imageFileSize   INT,
    post            TEXT,
    finished_at      TIMESTAMPTZ DEFAULT now(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS softwares (
    id              BIGSERIAL PRIMARY key,
    name            TEXT NOT NULL,
    image           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS illustration_softwares (
    illustration_id BIGINT REFERENCES illustrations(id) ON DELETE CASCADE,
    software_id     BIGINT REFERENCES softwares(id) ON DELETE CASCADE,
    PRIMARY KEY (illustration_id, software_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS illustrations;
-- +goose StatementEnd
