-- name: CreateCard :one
INSERT INTO cards(
  id,
  front_content,
  back_content,
  interval,
  ease_factor,
  repetitions_count,
  last_reviewed_at,
  created_at,
  deck_id
) VALUES (
  gen_random_uuid(),
  $1,
  $2,
  1,
  2.5,
  0,
  NULL,
  NOW(),
  $3
) RETURNING 
  id,
  front_content,
  back_content,
  interval,
  ease_factor,
  repetitions_count,
  last_reviewed_at,
  created_at,
  deck_id;
-- name: DeleteCards :exec
  DELETE FROM cards;
