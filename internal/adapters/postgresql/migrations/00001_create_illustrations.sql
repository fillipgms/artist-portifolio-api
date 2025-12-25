-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS illustrations (
    id              BIGSERIAL PRIMARY key,
    title           TEXT NOT NULL,
    description     TEXT NOT NULL,
    imageURL        TEXT NOT NULL,
    post            TEXT,
    finished_at      TIMESTAMPTZ DEFAULT now(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS illustrations;
-- +goose StatementEnd
