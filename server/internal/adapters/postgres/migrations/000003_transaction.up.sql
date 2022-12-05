CREATE TYPE transaction_type AS ENUM ('deposit', 'withdrawal', 'transfer');

CREATE TABLE IF NOT EXISTS "transactions"
(
    id              UUID PRIMARY KEY          DEFAULT uuid_generate_v4(),
    amount          NUMERIC(10, 2)   NOT NULL,
    type            transaction_type NOT NULL,
    created_at      timestamptz      NOT NULL DEFAULT NOW(),

    from_account_id BIGINT REFERENCES accounts (id),
    to_account_id   BIGINT REFERENCES accounts (id)

)