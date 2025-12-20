package tui

type (
	storeErrorMsg        struct{ error }
	refreshTasksMsg      struct{ projectID int64 }
	refreshProjectsMsg   struct{}
	projectsRefreshedMsg struct{}
)
