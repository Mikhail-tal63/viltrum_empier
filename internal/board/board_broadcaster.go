package board

import (
	"context"
	"log"

	"github.com/Mikhail-tal63/viltrum_empier/internal/websocket"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *BoardService) BroadcastBoardCreated(userID primitive.ObjectID, board *Board) {
	var workspaceID = board.WorkspaceID

	payload, err := websocket.MarshalEvent(websocket.EventBoardCreated, map[string]any{
		"board": board,
		"by":    userID.Hex(),
	})
	if err != nil {
		log.Printf("ws: marshal workspace failed: %v", err)
		return
	}

	s.boardHub.BroadcastToWorkspace(workspaceID.Hex(), payload)

}
func (s *BoardService) BroadcastBoardDelete(ctx context.Context, userID, boardID, workspaceID primitive.ObjectID) {

	payload, err := websocket.MarshalEvent(websocket.EventBoardDeleted, map[string]any{
		"board_id": boardID.Hex(),
		"by":       userID.Hex(),
	})
	if err != nil {
		log.Printf("ws: marshal workspace failed: %v", err)
		return
	}

	s.boardHub.BroadcastToWorkspace(workspaceID.Hex(), payload)
}

/*columns*********************************************************************************************************************************
******************************************************************************************************************************************/

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

func (s *BoardService) BroadcastColumnUpdated(ctx context.Context, userID primitive.ObjectID, column *Column) {

	workspaceID, err := s.repo.GetWorkspaceIDByColumn(ctx, column.ID)
	if err != nil {
		log.Printf("ws: resolve workspace failed: %v", err)
		return
	}

	payload, err := websocket.MarshalEvent(websocket.EventColumnUpdated, map[string]any{
		"column": column,
		"by":     userID.Hex(),
	})
	if err != nil {
		log.Printf("ws: marshal workspace failed: %v", err)
		return
	}

	s.boardHub.BroadcastToWorkspace(workspaceID.Hex(), payload)

}

func (s *BoardService) BroadcastColumnDelete(ctx context.Context, userID, colmunID, workspaceID primitive.ObjectID) {

	payload, err := websocket.MarshalEvent(websocket.EventColumnDeleted, map[string]any{
		"column": colmunID.Hex(),
		"by":     userID.Hex(),
	})
	if err != nil {
		log.Printf("ws: marshal workspace failed: %v", err)
		return
	}

	s.boardHub.BroadcastToWorkspace(workspaceID.Hex(), payload)
}

func (s *BoardService) BroadcastMoveColumn(ctx context.Context, userID primitive.ObjectID, column *Column) {
	workspaceID, err := s.repo.GetWorkspaceIDByColumn(ctx, column.ID)
	if err != nil {
		log.Printf("ws: resolve workspace failed: %v", err)
		return
	}
	payload, err := websocket.MarshalEvent(websocket.EventColumnMoved, map[string]any{
		"column": column,
		"by":     userID.Hex(),
	})
	if err != nil {
		log.Printf("ws: marshal workspace failed: %v", err)
		return
	}

	s.boardHub.BroadcastToWorkspace(workspaceID.Hex(), payload)
}
