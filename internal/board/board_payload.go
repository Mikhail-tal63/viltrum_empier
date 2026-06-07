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
