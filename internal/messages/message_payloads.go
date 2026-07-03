package messages

type CreateGroupChatPayload struct {
	EntityType ChatEntity `json:"entityType"`
}

type CreateMessagePayload struct {
	Content string `json:"content"`
}
