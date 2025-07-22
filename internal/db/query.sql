-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: CreateUser :exec
INSERT INTO users (
  username, email, password, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5
);

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: CreateImage :exec
INSERT INTO images (
    user_id,
    url
) VALUES (
    $1,
    $2
);

-- name: GetRandomUserWithImages :one
SELECT
    u.id,
    u.username,
    u.email,
    u.created_at,
    u.updated_at,
    COALESCE(
        (SELECT json_agg(json_build_object('url', i.url))
        FROM images i
        WHERE i.user_id = u.id),
        '[]'::json
    ) AS images
FROM
    users u
ORDER BY
    RANDOM()
LIMIT 1;

-- name: LikeUser :exec
INSERT INTO likes (user_id, n_of_likes)
VALUES ($1, 1)
ON CONFLICT (user_id) DO UPDATE
SET n_of_likes = likes.n_of_likes + 1;