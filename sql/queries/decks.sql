
-- name: CreateDeck :one
INSERT INTO decks (
 id,
 title,
 description,
 created_at,
 user_id,
 total_reviews,
 tempo_interval_up,
 tempo_interval_dn,
 perfect_threshold,
 bad_threshold
) VALUES (
  gen_random_uuid(),
	$1,
	$2,
  NOW(),
  $3,
  0,
  $4,
  $5,
  $6,
  $7
) RETURNING 
 id,
 title,
 description,
 created_at,
 user_id,
 total_reviews,
 tempo_interval_up,
 tempo_interval_dn,
 perfect_threshold,
 bad_threshold;

-- name: GetDeckById :one
  SELECT * FROM decks
  WHERE id = $1;

-- name: DeleteDecks :exec
  DELETE FROM decks;

-- name: DeleteDeck :exec
  DELETE FROM decks
  WHERE id = $1;

-- name: GetDecksByUser :many
  SELECT * FROM decks
  WHERE user_id = $1
  ORDER BY created_at ASC;

-- name: UpdateDeck :one
UPDATE decks
SET title  = $1,
    description  = $2,
    total_reviews  = $3,
    tempo_interval_up  = $4,
    tempo_interval_dn  = $5,
    perfect_threshold  = $6,
    bad_threshold = $7
WHERE id = $8
RETURNING *;
