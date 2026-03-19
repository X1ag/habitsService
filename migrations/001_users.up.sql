CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE auth_identities (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider TEXT NOT NULL, -- 'google' | 'telegram'
    provider_user_id TEXT NOT NULL, -- google sub или telegram user id
    email TEXT,
    username TEXT,
    display_name TEXT,
    avatar_url TEXT,
    telegram_chat_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_login_at TIMESTAMPTZ,
    UNIQUE (provider, provider_user_id)
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_seen_at TIMESTAMPTZ,
    user_agent TEXT,
);

CREATE TABLE IF NOT EXISTS habits (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL REFERENCES users(id),
	key_name VARCHAR(255),
	type VARCHAR(255),
	label VARCHAR(255),
	emoji VARCHAR(255),
	order_index INT,
	active BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);