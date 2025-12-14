-- name: ListAllTasks :many
SELECT * FROM task;

-- name: GetTask :one
SELECT * FROM task
WHERE id = ? LIMIT 1;

-- name: CreateTask :one
INSERT INTO task (
    title, complete
) VALUES (
    ?, ?
) RETURNING *;

-- name: UpdateTask :one
UPDATE task
SET
    title = ?,
    complete = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM task
WHERE id = ?;
