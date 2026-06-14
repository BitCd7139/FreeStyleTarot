CREATE TABLE predict_history (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  question   VARCHAR(2000) NOT NULL,
  model      VARCHAR(16)  NOT NULL DEFAULT 'default',
  card_size  JSONB        NOT NULL,
  cards      JSONB        NOT NULL,
  answer     TEXT         NOT NULL,
  created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_predict_history_user_created
  ON predict_history (user_id, created_at DESC);
