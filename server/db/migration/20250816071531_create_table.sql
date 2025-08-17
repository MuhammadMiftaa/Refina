-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    role VARCHAR(100) DEFAULT 'user'::VARCHAR NOT NULL,
    email_verfied_at timestamp without time zone
);

CREATE TABLE categories (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    parent_id uuid,
    parent_label VARCHAR(50),
    name VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL
);

CREATE TABLE wallet_types (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL
);

CREATE TABLE investment_types (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name VARCHAR(50) NOT NULL,
    unit text
);

CREATE TABLE wallets (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id uuid NOT NULL,
    wallet_type_id uuid NOT NULL,
    name VARCHAR(50) NOT NULL,
    balance numeric(18,2) NOT NULL,
    number VARCHAR(50) NOT NULL
);

CREATE TABLE investments (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    investment_type_id uuid NOT NULL,
    user_id uuid NOT NULL,
    name VARCHAR(50) NOT NULL,
    amount numeric(18,2) NOT NULL,
    quantity numeric(18,2) NOT NULL,
    investment_date timestamp without time zone NOT NULL,
    description text
);

CREATE TABLE transactions (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    wallet_id uuid NOT NULL,
    category_id uuid NOT NULL,
    amount numeric(18,2) NOT NULL,
    transaction_date timestamp without time zone NOT NULL,
    description text
);

CREATE TABLE attachments (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    transaction_id uuid NOT NULL,
    image text,
    format text,
    size bigint
);

CREATE TABLE reports (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id uuid NOT NULL,
    from_date timestamp with time zone NOT NULL,
    to_date timestamp with time zone NOT NULL,
    request_at timestamp with time zone NOT NULL,
    next_request_at timestamp with time zone NOT NULL,
    status report_status NOT NULL,
    file_url text,
    file_size bigint,
    generated_at timestamp with time zone
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reports;
DROP TABLE IF EXISTS attachments;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS investments;
DROP TABLE IF EXISTS wallets;
DROP TABLE IF EXISTS wallet_types;
DROP TABLE IF EXISTS investment_types;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
