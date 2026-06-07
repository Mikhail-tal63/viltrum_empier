package board

import (
	"context"
	"errors"

	"time"

	"github.com/Mikhail-tal63/viltrum_empier/internal/websocket"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardService struct {
	repo        *BoardRepository
	boardHub    *websocket.Hub
	taskDeleter TaskDeleter
}
type TaskDeleter interface {
	DeleteColumnTasks(
		ctx context.Context,
		columnID primitive.ObjectID,
	) error
}

func NewBoardService(repo *BoardRepository, boardHub *websocket.Hub, taskDeleter TaskDeleter) *BoardService {
	return &BoardService{
		repo:        repo,
		boardHub:    boardHub,
		taskDeleter: taskDeleter,
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

func (s *BoardService) DeleteBoard(ctx context.Context, userid, boardID primitive.ObjectID) error {

	board, err := s.repo.GetBoardByID(ctx, boardID)
	if err != nil {
		return err
	}
	if board == nil {
		return errors.New("board not found")
	}
	var workspaceID = board.WorkspaceID
	var columns = []*Column{}

	columns, err = s.repo.GetBoardColumns(ctx, boardID)
	if err != nil {
		return err
	}
	for _, column := range columns {

		if err := s.DeleteColumn(ctx, column.ID, userid); err != nil {
			return err
		}
	}

	if err := s.repo.DeleteBoard(ctx, boardID); err != nil {
		return err
	}

	s.BroadcastBoardDelete(ctx, userid, boardID, workspaceID)

	return nil
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

	if err := s.taskDeleter.DeleteColumnTasks(ctx, columnID); err != nil {
		return err
	}

	if err := s.repo.DeleteColumn(ctx, columnID); err != nil {
		return err
	}

	s.BroadcastColumnDelete(ctx, userID, columnID, workspaceID)

	return nil
}

func (s *BoardService) DragDropColumn(ctx context.Context, columnID, boardID, userID primitive.ObjectID, newPosition int) error {
	column, err := s.repo.GetColumnByID(ctx, columnID)
	if err != nil {
		return err
	}

	oldPosition := column.Position

	if oldPosition < newPosition {

		err = s.repo.ShiftPositions(
			ctx,
			boardID,
			oldPosition+1,
			newPosition,
			-1,
		)

	} else if oldPosition > newPosition {

		err = s.repo.ShiftPositions(
			ctx,
			boardID,
			newPosition,
			oldPosition-1,
			1,
		)
	}

	if err != nil {
		return err
	}

	if err := s.repo.UpdateColumnLocation(
		ctx,
		columnID,
		newPosition,
	); err != nil {
		return err
	}

	return nil
}
