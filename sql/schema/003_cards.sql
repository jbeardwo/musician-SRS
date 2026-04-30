-- +goose Up
CREATE TABLE cards(
  id UUID PRIMARY KEY,
  front_content TEXT,
  back_content TEXT,
  interval INTEGER,
  ease_factor FLOAT,
  repetitions_count INTEGER,
  last_reviewed_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL,
  deck_id UUID NOT NULL REFERENCES decks(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE cards;
