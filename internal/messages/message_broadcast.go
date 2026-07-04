package messages

import (
	"context"
	"log"

	"github.com/Mikhail-tal63/viltrum_empier/internal/websocket"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *MessageService) BroadcastCreateGroupChat(ctx context.Context, userID, workspaceID primitive.ObjectID, GC *GroupChat) {

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

func (s *MessageService) BroadcastMessageCreated(ctx context.Context, userID, gcID primitive.ObjectID, message *Message) {
	payload, err := websocket.MarshalEvent(websocket.EventMessageCreated, map[string]any{
		"message": message,
		"by":      userID.Hex(),
	})
	if err != nil {
		log.Printf("ws: marshal failed: %v", err)
		return
	}
	s.messageHup.BroadcastToGroupChat(gcID.Hex(), payload)

}
func (s *MessageService) BroadcastMessageDeleted(ctx context.Context, userID, gcID  ,  messageID primitive.ObjectID,) {
	payload, err := websocket.MarshalEvent(websocket.EventMessageDeleted, map[string]any{
    "messageId": messageID.Hex(),
    "by":        userID.Hex(),
})
	if err != nil {
		log.Printf("ws: marshal failed: %v", err)
		return
	}
	s.messageHup.BroadcastToGroupChat(gcID.Hex(), payload)
}

func (s *MessageService) BroadcastMessageEdited(ctx context.Context, userID, gcID primitive.ObjectID, message *Message) {
	payload, err := websocket.MarshalEvent(websocket.EventMessageUpdated, map[string]any{
		"message": message,
		"by":      userID.Hex(),
	})
	if err != nil {
		log.Printf("ws: marshal failed: %v", err)
		return
	}
	s.messageHup.BroadcastToGroupChat(gcID.Hex(), payload)

}
