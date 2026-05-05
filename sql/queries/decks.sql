
-- name: CreateDeck :one
INSERT INTO decks (id, title, description, created_at, user_id)
VALUES (
  gen_random_uuid(),
	$1,
	$2,
  NOW(),
  $3
)
RETURNING id, title, description, created_at, user_id;

-- name: DeleteDecks :exec
  DELETE FROM decks;

-- name: GetDecksByUser :many
SELECT * FROM decks
WHERE user_id = $1;
