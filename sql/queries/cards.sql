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
  created_at,
  deck_id,
  tempo,
  perfect_streak,
  bad_streak;


-- name: DeleteCards :exec
  DELETE FROM cards;
