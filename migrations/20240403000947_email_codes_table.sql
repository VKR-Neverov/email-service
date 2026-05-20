-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS email_codes(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    code TEXT NOT NULL,
    sent_to TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT NOW(),
    confirmed_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE email_codes;
-- +goose StatementEnd
