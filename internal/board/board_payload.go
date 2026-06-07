package board

type CreateBoardPayload struct {
	Name string `json:"name"`
}

type ColumnPayload struct {
	Name string `json:"name"`
}
type DragColumnPayload struct {
	Position int `json:"position"`
}
type PatchColumnPayload struct {
	Name       *string `json:"name"`
	Color      *string `json:"color"`
	IsArchived *bool   `json:"is_archived"`
}
