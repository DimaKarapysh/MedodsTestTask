-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS token
(

    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token_hash TEXT NOT NULL,
    client_ip TEXT NOT NULL,
    issued_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,




    created_at timestamp(0) NOT NULL DEFAULT now(),
    updated_at timestamp(0) NOT NULL DEFAULT now(),
    deleted_at timestamp(0)          DEFAULT NULL,

    constraint fk_user_id foreign key (user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS token;
-- +goose StatementEnd
