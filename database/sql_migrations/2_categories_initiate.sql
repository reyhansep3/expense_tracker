-- +migrate Up
-- +migrate statementBegin

CREATE TABLE categories(
    id BIGINT PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    category_name VARCHAR(60),
    create_at TIMESTAMP,
    create_by VARCHAR(60),
    modified_at TIMESTAMP,
    modified_by VARCHAR(60)
);

-- +statementEnd