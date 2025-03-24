-- +goose Up
-- +goose StatementBegin
CREATE TABLE requests
(
    uuid    VARCHAR(255) PRIMARY KEY,
    value   INTEGER not null,
    timeout INTEGER not null,
    result  INTEGER,
    state   VARCHAR(255) default 'NEW'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE requests;
-- +goose StatementEnd
