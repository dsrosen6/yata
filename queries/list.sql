-- name: ListAllLists :many
SELECT * FROM list;

-- name: ListListsByParentListID :many
SELECT * FROM list
WHERE parent_list_id = ?;

-- name: GetList :one
SELECT * FROM list
WHERE id = ? LIMIT 1;

-- name: CreateList :one
INSERT INTO list (
    title, 
    parent_list_id
) VALUES (
    ?, ?
) RETURNING *;

-- name: UpdateList :one
UPDATE list
SET
    title = ?,
    parent_list_id = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteList :exec
DELETE FROM list
WHERE id = ?;
