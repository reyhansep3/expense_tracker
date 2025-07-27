-- +migrate Up
-- +migrate statementBegin

CREATE TABLE expenses (
    id BIGINT PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    category_id BIGINT REFERENCES categories(id) ON DELETE CASCADE,
    payment_method VARCHAR(60),
    title VARCHAR(255),
    amount INT NOT NULL,
    description VARCHAR(255),
    expense_date DATE,
    create_at TIMESTAMP,
    create_by BIGINT,
    modified_at TIMESTAMP,
    modified_by BIGINT
);
-- +migrate statementEnd