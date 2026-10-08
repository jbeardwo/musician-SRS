-- +goose Up
CREATE TABLE cards(
  id UUID PRIMARY KEY,
  front_content TEXT NOT NULL,
  back_content TEXT NOT NULL,
  interval INTEGER NOT NULL DEFAULT 1,
  target INTEGER NOT NULL,
  ease_factor FLOAT NOT NULL DEFAULT 2.5,
  repetitions_count INTEGER NOT NULL DEFAULT 0,
  last_reviewed_at TIMESTAMP,
  last_reviewed_num INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deck_id UUID NOT NULL REFERENCES decks(id) ON DELETE CASCADE,
	tempo INTEGER NOT NULL,
  mastered_tempo INTEGER NOT NULL DEFAULT 0,
  mastery_stage TEXT NOT NULL DEFAULT 'playable'
    CHECK (mastery_stage IN ('challenge', 'playable', 'retained')),
  mastery_stage_time_started TIMESTAMP,
  challenge_again_count INTEGER NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE cards;
