package messages

type CreateGroupChatPayload struct {
	EntityType ChatEntity `json:"entityType"`
}

type CreateMessagePayload struct {
	Content string `json:"content"`
}

type EditMessagePayload struct {
	Content string `json:"content"`
}
