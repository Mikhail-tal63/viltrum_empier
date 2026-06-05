package task

import (
	"context"
	"log"

	"github.com/Mikhail-tal63/viltrum_empier/internal/websocket"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *TaskService) broadcastTaskCreated(
	ctx context.Context,
	columnID primitive.ObjectID,
	task *Task,
	userID primitive.ObjectID,
) {

	workspaceID, err := s.boardRepo.GetWorkspaceIDByColumn(
		ctx,
		columnID,
	)

	if err != nil {
		log.Printf("ws: resolve workspace failed: %v", err)
		return
	}

	payload, err := websocket.MarshalEvent(
		websocket.EventTaskCreated,
		map[string]any{
			"task": task,
			"by":   userID.Hex(),
		},
	)

	if err != nil {
		log.Printf("ws: marshal failed: %v", err)
		return
	}

	s.hub.BroadcastToWorkspace(
		workspaceID.Hex(),
		payload,
	)
}

func (s *TaskService) broadcastTaskUpdated(ctx context.Context, userID primitive.ObjectID, task *Task) {
	workspaceID, err := s.boardRepo.GetWorkspaceIDByColumn(ctx, task.ColumnID)
	if err != nil {
		log.Printf("ws: resolve workspace failed: %v", err)
		return
	}

	payload, err := websocket.MarshalEvent(websocket.EventTaskUpdated, map[string]any{
		"task": task,
		"by":   userID.Hex(),
	},
	)

	if err != nil {
		log.Printf("ws: marshal failed: %v", err)
		return
	}

	s.hub.BroadcastToWorkspace(workspaceID.Hex(), payload)

}
