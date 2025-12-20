-- name: ListAllTasks :many
SELECT * FROM task;

-- name: ListTasksByProjectID :many
SELECT * FROM task
WHERE project_id = ?;

-- name: ListTasksByParentTaskID :many
SELECT * FROM task
WHERE parent_task_id = ?;

-- name: GetTask :one
SELECT * FROM task
WHERE id = ? LIMIT 1;

-- name: CreateTask :one
INSERT INTO task (
    title,
    parent_task_id,
    project_id,
    complete,
    due_at
) VALUES (
    ?, ?, ?, ?, ?
) RETURNING *;

-- name: UpdateTask :one
UPDATE task
SET
    title = ?,
    parent_task_id = ?,
    project_id = ?,
    complete = ?,
    due_at = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM task
WHERE id = ?;
