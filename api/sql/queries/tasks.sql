-- name: GetTasksByUser :many
SELECT * FROM tasks WHERE user_id = $1;

-- name: CreateTask :one
INSERT INTO tasks (id, user_id, title, description, priority, due_date) VALUES ($1, $2, $3, $4, $5, $6) RETURNING title, description, priority, due_date;
