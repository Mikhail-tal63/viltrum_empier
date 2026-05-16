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

		CreatedAt: now,
		UpdatedAt: now,
	}

	createdBoard, err := s.repo.CreateBoard(board)
	if err != nil {
		return nil, err
	}
	return createdBoard, nil
}
