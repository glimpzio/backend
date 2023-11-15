\c dev

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    auth_id TEXT NOT NULL UNIQUE,
    personal_email TEXT NOT NULL,
    name TEXT NOT NULL,
    bio TEXT NOT NULL,
    profile_picture TEXT,
    email TEXT,
    phone TEXT,
    website TEXT,
    linkedin TEXT
);

CREATE TABLE links (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL
);