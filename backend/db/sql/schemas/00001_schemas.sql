-- +goose Up
CREATE    SCHEMA IF NOT EXISTS auth;

-- +goose Down
DROP      SCHEMA IF EXISTS auth CASCADE;