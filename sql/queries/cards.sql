-- name: CreateCard :one
INSERT INTO cards(
  id,
  front_content,
  back_content,
  interval,
  target,
  ease_factor,
  repetitions_count,
  last_reviewed_at,
  last_reviewed_num,
  created_at,
  deck_id,
  tempo,
  perfect_streak,
  bad_streak
) VALUES (
  gen_random_uuid(),
  $1,
  $2,
  1,
  $3,
  2.5,
  0,
  NULL,
  0,
  NOW(),
  $4,
  $5,
  0,
  0
) RETURNING 
  id,
  front_content,
  back_content,
  interval,
  target,
  ease_factor,
  repetitions_count,
  last_reviewed_at,
  last_reviewed_num,
  created_at,
  deck_id,
  tempo,
  perfect_streak,
  bad_streak;

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
  perfect_streak = $10,
  bad_streak = $11
WHERE id = $12
RETURNING *;
