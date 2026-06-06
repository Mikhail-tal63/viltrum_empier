package board

import (
	"context"
	"log"

	"github.com/Mikhail-tal63/viltrum_empier/internal/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *BoardService) BroadcastColumnCreated(ctx context.Context, userID primitive.ObjectID, column *Column) {
	workspaceID, err := s.repo.GetWorkspaceIDByColumn(ctx, column.ID)
	if err != nil {
		log.Printf("ws: resolve workspace failed: %v", err)
		return
	}

	payload, err := websocket.MarshalEvent(websocket.EventColumnCreated, map[string]any{
		"column": column,
		"by":     userID.Hex(),
	})
	if err != nil {
		log.Printf("ws: marshal workspace failed: %v", err)
		return
	}

	s.boardHub.BroadcastToWorkspace(workspaceID.Hex(), payload)
}
