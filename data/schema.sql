\c dev

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Profile service
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    auth_id TEXT NOT NULL UNIQUE,
    personal_email TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    bio TEXT NOT NULL,
    profile_picture TEXT,
    email TEXT,
    phone TEXT,
    website TEXT,
    linkedin TEXT
);

CREATE TABLE invites (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    -- Upon migrating to microservices remove "user_id" as a forein key
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL
);

-- Connections service
CREATE TABLE custom_connections (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    -- Upon migrating to microservices remove "user_id" as a forein key
    user_id UUID NOT NULL references users(id) ON DELETE CASCADE,
    connected_at TIMESTAMP NOT NULL DEFAULT NOW(),
    first_name TEXT,
    last_name TEXT,
    notes TEXT,
    email TEXT,
    phone TEXT,
    website TEXT,
    linkedin TEXT
);