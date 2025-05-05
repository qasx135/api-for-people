-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS persons(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255),
    age INT,
    gender VARCHAR(255),
    nationality VARCHAR(255),
    isDeleted BOOLEAN DEFAULT false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE persons;
-- +goose StatementEnd
