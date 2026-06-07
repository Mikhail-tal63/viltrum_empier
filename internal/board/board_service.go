package board

import (
	"context"

	"time"

	"github.com/Mikhail-tal63/viltrum_empier/internal/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardService struct {
	repo     *BoardRepository
	boardHub *websocket.Hub
}

func NewBoardService(repo *BoardRepository, boardHub *websocket.Hub) *BoardService {
	return &BoardService{
		repo:     repo,
		boardHub: boardHub,
	}
}

func (s *BoardService) CreateBoard(ctx context.Context, payload *CreateBoardPayload, userID, workspaceId primitive.ObjectID) (*Board, error) {
	now := time.Now()

	boardID := primitive.NewObjectID()

	board := &Board{
		ID:          boardID,
		WorkspaceID: workspaceId,
		Name:        payload.Name,
		Position:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	createdBoard, err := s.repo.CreateBoard(ctx, board)
	if err != nil {
		return nil, err
	}

	s.BroadcastBoardCreated(userID, createdBoard)

	return createdBoard, nil
}

func (s *BoardService) GetWorkspaceBoards(ctx context.Context, workspaceID primitive.ObjectID) ([]*Board, error) {
	return s.repo.GetWorkspaceBoards(ctx, workspaceID)

}

/*colmuns ****************/

func (s *BoardService) CreateColmun(ctx context.Context, boardID, userid primitive.ObjectID, payload *ColumnPayload) (*Column, error) {
	now := time.Now()
	colmunID := primitive.NewObjectID()

	count, err := s.repo.CountBoardColumns(ctx, boardID)
	if err != nil {
		return nil, err
	}

	colmun := &Column{
		ID:         colmunID,
		BoardID:    boardID,
		Name:       payload.Name,
		Position:   int(count),
		Color:      "",
		IsArchived: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	colmun, err = s.repo.CreateColmun(ctx, colmun)
	if err != nil {
		return nil, err
	}

	s.BroadcastColumnCreated(ctx, userid, colmun)
	return colmun, nil

}

func (s *BoardService) GetBoarderColumns(ctx context.Context, workspaceID primitive.ObjectID) ([]*Column, error) {
	return s.repo.GetBoardColumns(ctx, workspaceID)
}

func (s *BoardService) DeleteColumn(ctx context.Context, columnID, userID primitive.ObjectID) error {

	workspaceID, err := s.repo.GetWorkspaceIDByColumn(ctx, columnID)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteColumn(ctx, columnID); err != nil {
		return err
	}

	s.BroadcastColumnDelete(ctx, userID, columnID, workspaceID)

	return nil
}
