CREATE TYPE service_provider_type AS ENUM ('yandex');

CREATE TABLE oauth_accounts (
    id                  uuid                    PRIMARY     KEY DEFAULT gen_random_uuid(),
    user_id             uuid                    NOT NULL    REFERENCES users(id) ON DELETE CASCADE,
    provider            service_provider_type   NOT NULL,
    provider_user_id    text                    NOT NULL,
    created_at          timestamptz             NOT NULL    DEFAULT now(),

    UNIQUE(provider, provider_user_id)
);

CREATE INDEX idx_oauth_accounts_user_id ON oauth_accounts(user_id);