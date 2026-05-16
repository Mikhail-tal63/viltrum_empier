package board

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardService struct {
	repo *BoardRepository
}

func NewBoardService(repo *BoardRepository) *BoardService {
	return &BoardService{
		repo: repo,
	}
}

func (s *BoardService) CreateBoard(ctx context.Context, payload *CreateBoardPayload, workspaceId primitive.ObjectID) (*Board, error) {
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
	return createdBoard, nil
}

func (s *BoardService) GetWorkspaceBoards(ctx context.Context, workspaceID primitive.ObjectID) ([]*Board, error) {
	return s.repo.GetWorkspaceBoards(ctx, workspaceID)

}

/*colmuns ****************/

func (s *BoardService) CreateColmun(ctx context.Context, boardID primitive.ObjectID, payload *ColumnPayload) (*Column, error) {
	now := time.Now()
	colmunID := primitive.NewObjectID()

	colmun := &Column{
		ID:         colmunID,
		BoardID:    boardID,
		Name:       payload.Name,
		Position:   0,
		Color:      "",
		IsArchived: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	colmun, err := s.repo.CreateColmun(ctx, colmun)
	if err != nil {
		return nil, err
	}
	return colmun, nil

}
