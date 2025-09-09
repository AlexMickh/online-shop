CREATE TYPE token_type AS ENUM(
    'email-verify'
);

CREATE TABLE IF NOT EXISTS tokens(
    token TEXT PRIMARY KEY,
    user_id UUID REFERENCES users (id),
    type token_type
);