-- +goose Up
CREATE TABLE cards(
  id UUID PRIMARY KEY,
  front_content TEXT NOT NULL,
  back_content TEXT NOT NULL,
  interval INTEGER NOT NULL,
  target INTEGER NOT NULL,
  ease_factor FLOAT NOT NULL,
  repetitions_count INTEGER NOT NULL,
  last_reviewed_at TIMESTAMP,
  last_reviewed_num INTEGER NOT NULL,
  created_at TIMESTAMP NOT NULL,
  deck_id UUID NOT NULL REFERENCES decks(id) ON DELETE CASCADE,
	tempo INTEGER NOT NULL,
	perfect_streak INTEGER NOT NULL,
	bad_streak INTEGER NOT NULL 
);

-- +goose Down
DROP TABLE cards;
