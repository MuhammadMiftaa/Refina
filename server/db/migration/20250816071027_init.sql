-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_cron";
CREATE TYPE report_status AS ENUM ('processing', 'completed', 'failed');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS report_status;
DROP EXTENSION IF EXISTS "pg_cron";
DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd
