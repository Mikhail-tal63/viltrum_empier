package task

type CreateTaskPayload struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Priority        string   `json:"priority"`
	DueDate         string   `json:"due_date"`
	Labels          []string `json:"labels"`
	AssignedMembers []string `json:"assigned_members"`
}
type EditTaskPayload struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Priority        string   `json:"priority"`
	Labels          []string `json:"labels"`
	AssignedMembers []string `json:"assigned_members"`
}
type DragPayload struct {
	Position int `json:"position"`

	NewColumnID string `json:"new_column_id"`
}
