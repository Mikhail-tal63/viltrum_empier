package task

type CreateTaskPayload struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Priority        string   `json:"priority"`
	DueDate         string   `json:"due_date"`
	Labels          []string `json:"labels"`
	AssignedMembers []string `json:"assigned_members"`
}
