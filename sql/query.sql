-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING id, email, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, password_hash, created_at, updated_at
FROM users
WHERE id = $1;

-- name: StoreRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING id, user_id, created_at, expires_at, is_revoked;

-- name: GetRefreshToken :one
SELECT id, user_id, token_hash, created_at, expires_at, is_revoked
FROM refresh_tokens
WHERE id = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET is_revoked = TRUE
WHERE id = $1;

-- name: RevokeAllUserTokens :exec
UPDATE refresh_tokens
SET is_revoked = TRUE
WHERE user_id = $1;

-- name: CreateSession :one
INSERT INTO sessions (user_id, refresh_token_id, user_agent, ip_address, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, refresh_token_id, created_at, expires_at, is_revoked;

-- name: GetSessionByID :one
SELECT id, user_id, refresh_token_id, user_agent, ip_address, created_at, expires_at, is_revoked
FROM sessions
WHERE id =$1;

-- name: RevokeSession :exec
UPDATE sessions
SET is_revoked = TRUE
WHERE id = $1;

-- name: RevokeAllUserSessions :exec
UPDATE sessions
SET is_revoked = TRUE
WHERE user_id = $1;