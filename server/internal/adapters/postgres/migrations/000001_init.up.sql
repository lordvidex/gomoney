CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE currency AS ENUM ('USD', 'EUR', 'RUB');

CREATE TABLE IF NOT EXISTS "users" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS "accounts" (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    balance NUMERIC(10,2) DEFAULT 0,
    currency Currency NOT NULL,
    is_blocked BOOLEAN DEFAULT FALSE,

    user_id UUID REFERENCES "users" (id) ON DELETE CASCADE
);

