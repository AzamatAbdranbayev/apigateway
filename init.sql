CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT gen_random_uuid () UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    middle_name VARCHAR(255) NOT NULL,
    login VARCHAR(255),
    group_number VARCHAR(255) NOT NULL,
    balance float,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

-- CREATE TABLE IF NOT EXISTS tasks (
--     id UUID DEFAULT gen_random_uuid () UNIQUE,
--     description VARCHAR(255)
--     cost  float  NOT NULL,
--     created_at TIMESTAMPTZ NOT NULL,
--     updated_at TIMESTAMPTZ NOT NULL
-- );