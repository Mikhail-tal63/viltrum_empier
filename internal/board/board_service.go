package board

import (
	"context"
	"errors"

	"time"

	"github.com/Mikhail-tal63/viltrum_empier/internal/permission"
	"github.com/Mikhail-tal63/viltrum_empier/internal/websocket"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardService struct {
	repo              *BoardRepository
	boardHub          *websocket.Hub
	taskDeleter       TaskDeleter
	permissionChecker PermissionChecker
	groupChatCreator  GroupChatCreator
}

type GroupChatCreator interface {
	CreateBoardGroupChat(ctx context.Context, boardID primitive.ObjectID) error
	CreateColumnGroupChat(ctx context.Context, columnID primitive.ObjectID) error
}

type TaskDeleter interface {
	DeleteColumnTasks(
		ctx context.Context,
		columnID primitive.ObjectID,
	) error
}
type PermissionChecker interface {
	HasPermission(
		ctx context.Context,
		userID, workspaceID primitive.ObjectID,
		PermDeleteBoard string,
	) bool
}

func NewBoardService(repo *BoardRepository, boardHub *websocket.Hub, taskDeleter TaskDeleter, permissionChecker PermissionChecker) *BoardService {
	return &BoardService{
		repo:              repo,
		boardHub:          boardHub,
		taskDeleter:       taskDeleter,
		permissionChecker: permissionChecker,
	}
}

func (s *BoardService) SetGroupChatCreator(gc GroupChatCreator) {
	s.groupChatCreator = gc
}

func (s *BoardService) CreateBoard(ctx context.Context, payload *CreateBoardPayload, userID, workspaceId primitive.ObjectID) (*Board, error) {

	allawed := s.permissionChecker.HasPermission(ctx, userID, workspaceId, permission.PermCreateBoard)

	if !allawed {
		return nil, errors.New("permission denied")
	}
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

	if s.groupChatCreator != nil {
		if err := s.groupChatCreator.CreateBoardGroupChat(ctx, boardID); err != nil {
			return nil, err
		}
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

	allawed := s.permissionChecker.HasPermission(ctx, userid, workspaceID, permission.PermDeleteTask)

	if !allawed {
		return errors.New("permission denied")
	}

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

func (s *BoardService) UpdateBoardDetails(ctx context.Context, boardID, userID primitive.ObjectID, payload *PatchBoardPayload) error {

	checkboard, err := s.repo.GetBoardByID(ctx, boardID)
	if err != nil {
		return err
	}
	if checkboard == nil {
		return errors.New("board not found")
	}

	allowed := s.permissionChecker.HasPermission(
		ctx,
		userID,
		checkboard.WorkspaceID,
		permission.PermEditBoard,
	)

	if !allowed {
		return errors.New("permission denied")
	}

	update := bson.M{}

	if payload.Name != nil {
		update["name"] = *payload.Name
	}

	update["updated_at"] = time.Now()

	if err := s.repo.UpdateBoardDetails(ctx, boardID, update); err != nil {
		return err
	}

	board, err := s.repo.GetBoardByID(ctx, boardID)
	if err != nil {
		return err
	}

	s.BroadcastBoardUpdated(ctx, userID, board)

	return nil
}

/*colmuns ****************/

func (s *BoardService) CreateColmun(ctx context.Context, boardID, userid primitive.ObjectID, payload *ColumnPayload) (*Column, error) {

	checkboard, err := s.repo.GetBoardByID(ctx, boardID)
	if err != nil {
		return nil, err
	}
	if checkboard == nil {
		return nil, errors.New("board not found")
	}

	allowed := s.permissionChecker.HasPermission(
		ctx,
		userid,
		checkboard.WorkspaceID,
		permission.PermCreateColumn,
	)

	if !allowed {

		return nil, errors.New("permission denied")
	}

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

	if s.groupChatCreator != nil {
		if err := s.groupChatCreator.CreateColumnGroupChat(ctx, colmunID); err != nil {
			return nil, err
		}
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
	allowed := s.permissionChecker.HasPermission(
		ctx,
		userID,
		workspaceID,
		permission.PermDeleteColumn,
	)

	if !allowed {

		return errors.New("permission denied")
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

	var worspaceID, err = s.repo.GetWorkspaceIDByColumn(ctx, columnID)
	if err != nil {
		return err
	}

	allowed := s.permissionChecker.HasPermission(
		ctx,
		userID,
		worspaceID,
		permission.PermEditColumn,
	)

	if !allowed {

		return errors.New("permission denied")
	}

	column, err := s.repo.GetColumnByID(ctx, columnID)
	if err != nil {
		return err
	}
	if column == nil {
		return errors.New("column not found")
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

	updatedColumn, err := s.repo.GetColumnByID(ctx, columnID)
	if err != nil {
		return err
	}

	s.BroadcastMoveColumn(ctx, userID, updatedColumn)
	return nil
}

func (s *BoardService) UpdateColumnDetails(ctx context.Context, columnID, userID primitive.ObjectID, payload *PatchColumnPayload) error {

	var worspaceID, err = s.repo.GetWorkspaceIDByColumn(ctx, columnID)
	if err != nil {
		return err
	}

	allowed := s.permissionChecker.HasPermission(
		ctx,
		userID,
		worspaceID,
		permission.PermEditColumn,
	)

	if !allowed {

		return errors.New("permission denied")
	}

	update := bson.M{}

	if payload.Name != nil {
		update["name"] = *payload.Name
	}

	if payload.Color != nil {
		update["color"] = *payload.Color
	}

	if payload.IsArchived != nil {
		update["is_archived"] = *payload.IsArchived
	}

	update["updated_at"] = time.Now()

	if err := s.repo.UpdateColumnDetails(ctx, columnID, update); err != nil {
		return err
	}

	updated, err := s.repo.GetColumnByID(ctx, columnID)
	if err != nil {
		return err
	}

	s.BroadcastColumnUpdated(ctx, userID, updated)

	return nil
}
