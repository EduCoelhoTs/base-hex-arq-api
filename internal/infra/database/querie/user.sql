-- name: CreateUser :exec
INSERT INTO auth.users (id, first_name, last_name, email, password, birth_date, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW());

-- name: GetUserById :one
SELECT id, first_name, last_name, email, password, birth_date, created_at, updated_at
FROM auth.users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, first_name, last_name, email, password, birth_date, created_at, updated_at
FROM auth.users
WHERE email = $1;

-- name: GetAllUsers :many
SELECT id, first_name, last_name, email, password, birth_date, created_at, updated_at
FROM auth.users
ORDER BY created_at DESC;

-- name: UpdateUser :exec
UPDATE auth.users
SET first_name = $2,
    last_name = $3,
    email = $4,
    password = $5,
    birth_date = $6,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM auth.users
WHERE id = $1;