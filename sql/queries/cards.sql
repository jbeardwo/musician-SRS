-- name: CreateCard :one
INSERT INTO cards(
  id,
  front_content,
  back_content,
  target,
  deck_id,
  tempo
) VALUES (
  gen_random_uuid(),
  $1,
  $2,
  $3,
  $4,
  $5
)
RETURNING *;

-- name: GetCardById :one
SELECT * FROM cards
WHERE id = $1;

-- name: DeleteCards :exec
DELETE FROM cards;

-- name: DeleteCard :exec
DELETE FROM cards
WHERE id = $1;

-- name: GetCardsByDeck :many
SELECT * FROM cards
WHERE deck_id = $1;

-- name: UpdateCard :one
UPDATE cards
SET  front_content = $1,
  back_content = $2,
  interval = $3,
  target = $4,
  ease_factor = $5,
  repetitions_count = $6,
  last_reviewed_at = $7,
  last_reviewed_num = $8,
  tempo = $9,
  mastered_tempo = $10,
  mastery_stage = $11,
  mastery_stage_time_started = $12,
  challenge_again_count = $13
WHERE id = $14
RETURNING *;
