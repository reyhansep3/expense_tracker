-- +migrate Up
-- +migrate statementBegin

CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    name VARCHAR(60) NOT NULL,
    password VARCHAR(60) NOT NULL,
    email VARCHAR(60) NOT NULL,
    token VARCHAR(255)
);

-- +migrate statementEnd