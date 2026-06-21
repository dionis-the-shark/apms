package authz

type Action string

const (
	ProjectCreate Action = "project.create"
	TaskCreate    Action = "task.create"
	TaskRead      Action = "task.read"
	TaskUpdate    Action = "task.update"
	TaskDelete    Action = "task.delete"
)
