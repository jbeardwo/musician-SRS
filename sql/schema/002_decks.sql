-- +goose Up
CREATE TABLE decks(
  id UUID PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  total_reviews INTEGER NOT NULL
);

-- +goose Down
DROP TABLE decks;
