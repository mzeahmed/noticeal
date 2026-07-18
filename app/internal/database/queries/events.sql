-- name: CreateEvent :one
INSERT INTO events (source, type, status, title, message, data)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING id, source, type, status, title, message, data, created_at;