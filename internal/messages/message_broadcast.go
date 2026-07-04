package messages

import (
	"context"
	"log"

	"github.com/Mikhail-tal63/viltrum_empier/internal/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *MessageService) BroadcastCreateGroupChat(ctx context.Context, workspaceID, userID primitive.ObjectID, GC *GroupChat) {

	payload, err := websocket.MarshalEvent(websocket.EventGroupCreated, map[string]any{
		"group_chat": GC,
		"by":         userID.Hex(),
	})
	if err != nil {
		log.Printf("ws: marshal failed: %v", err)
		return
	}
	s.messageHup.BroadcastToWorkspace(workspaceID.Hex(), payload)

}
