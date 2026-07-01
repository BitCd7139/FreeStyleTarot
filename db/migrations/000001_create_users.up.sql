CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email           VARCHAR(64)  NOT NULL UNIQUE,
  nickname        VARCHAR(32)  NOT NULL DEFAULT '',
  password_hash   VARCHAR(72),

  tier            VARCHAR(16)  NOT NULL DEFAULT 'free',
  vip_expires_at  TIMESTAMPTZ,

  last_predict_at TIMESTAMPTZ,

  created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users (email);
