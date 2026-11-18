CREATE TABLE users (
    id              uuid        PRIMARY     KEY,
    email           text        NOT NULL    UNIQUE,
    username        text        NOT NULL    UNIQUE,
    password_hash   text,
    is_active       boolean     NOT NULL,
    created_at      timestamptz NOT NULL    DEFAULT now(),
    updated_at      timestamptz NOT NULL    DEFAULT now()
);