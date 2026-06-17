
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

-- name: DeleteDecks :exec
  DELETE FROM decks;

-- name: GetDecksByUser :many
SELECT * FROM decks
WHERE user_id = $1;
