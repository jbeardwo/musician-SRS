-- +goose Up
CREATE TABLE cards(
  id UUID PRIMARY KEY,
  front_content TEXT NOT NULL,
  back_content TEXT NOT NULL,
  interval INTEGER NOT NULL,
  ease_factor FLOAT NOT NULL,
  repetitions_count INTEGER NOT NULL,
  last_reviewed_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL,
  deck_id UUID NOT NULL REFERENCES decks(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE cards;
