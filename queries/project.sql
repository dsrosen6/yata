-- name: ListAllProjects :many
SELECT * FROM project;

-- name: ListProjectsByParentProjectID :many
SELECT * FROM project
WHERE parent_project_id = ?;

-- name: GetProject :one
SELECT * FROM project
WHERE id = ? LIMIT 1;

-- name: CreateProject :one
INSERT INTO project (
    title,
    parent_project_id
) VALUES (
    ?, ?
) RETURNING *;

-- name: UpdateProject :one
UPDATE project
SET
    title = ?,
    parent_project_id = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM project
WHERE id = ?;
