
-- name: CreateDeck :one
INSERT INTO decks (
 id,
 title,
 description,
 created_at,
 user_id,
 total_reviews
) VALUES (
  gen_random_uuid(),
	$1,
	$2,
  NOW(),
  $3,
  0
) RETURNING 
 id,
 title,
 description,
 created_at,
 user_id,
 total_reviews;

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
    total_reviews  = $3
WHERE id = $4
RETURNING *;
