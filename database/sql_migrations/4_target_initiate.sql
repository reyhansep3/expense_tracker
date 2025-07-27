-- +migrate Up
-- +migrate statementBegin
CREATE TABLE target (
    id BIGINT PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    file VARCHAR(255),
    title VARCHAR(255),
    payment_method VARCHAR(60),
    description VARCHAR(255),
    amount INT,
    total_amount INT,
    start_date DATE,
    end_date DATE,
    create_at TIMESTAMP,
    create_by BIGINT,
    modified_at TIMESTAMP,
    modified_by BIGINT
);
-- +migrate statementEnd