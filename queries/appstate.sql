-- name: GetAppState :one
SELECT * FROM app_state
LIMIT 1;

-- name: CreateAppState :one
INSERT INTO app_state (
    selected_project_id
) VALUES (
    ?
)
RETURNING *;

-- name: UpdateAppState :one
UPDATE app_state
SET
    selected_project_id = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = 1
RETURNING *;
