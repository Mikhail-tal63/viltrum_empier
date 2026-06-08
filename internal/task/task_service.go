package task

import (
	"context"
	"errors"

	"time"

	"github.com/Mikhail-tal63/viltrum_empier/internal/board"
	"github.com/Mikhail-tal63/viltrum_empier/internal/permission"
	"github.com/Mikhail-tal63/viltrum_empier/internal/websocket"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskService struct {
	repo              *TaskRepository
	hub               *websocket.Hub
	boardRepo         *board.BoardRepository
	permissionChecker PermissionChecker
}
type PermissionChecker interface {
	HasPermission(
		ctx context.Context,
		userID, workspaceID primitive.ObjectID,
		PermDeleteBoard string,
	) bool
}

func NewTaskService(repo *TaskRepository, hub *websocket.Hub, boardRepo *board.BoardRepository, permissionChecker PermissionChecker) *TaskService {
	return &TaskService{
		repo:              repo,
		hub:               hub,
		boardRepo:         boardRepo,
		permissionChecker: permissionChecker,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, payload *CreateTaskPayload, user, columnID primitive.ObjectID) (*Task, error) {

	var worspaceID, err = s.boardRepo.GetWorkspaceIDByColumn(ctx, columnID)
	if err != nil {
		return nil, err
	}

	allowed := s.permissionChecker.HasPermission(
		ctx,
		user,
		worspaceID,
		permission.PermCreateTask,
	)

	if !allowed {

		return nil, errors.New("permission denied")
	}

	now := time.Now()
	taskID := primitive.NewObjectID()
	count, err := s.repo.CountTaskColumns(ctx, columnID)
	if err != nil {
		return nil, err
	}
	var dueDate *time.Time

	if payload.DueDate != "" {

		parsed, err := time.Parse(time.RFC3339, payload.DueDate)
		if err != nil {
			return nil, err
		}

		dueDate = &parsed
	}
	assignedMembers := []primitive.ObjectID{}

	for _, id := range payload.AssignedMembers {

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}

		assignedMembers = append(assignedMembers, objectID)
	}
	task := &Task{
		ID:              taskID,
		ColumnID:        columnID,
		Title:           payload.Title,
		Description:     payload.Description,
		Priority:        payload.Priority,
		Position:        int(count),
		CreatedBy:       user,
		AssignedMembers: assignedMembers,
		Labels:          payload.Labels,
		DueDate:         dueDate,
		IsArchived:      false,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	createdTask, err := s.repo.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}

	s.broadcastTaskCreated(ctx, columnID, createdTask, user)

	return createdTask, nil
}

func (s *TaskService) ListTasks(ctx context.Context, columnID primitive.ObjectID) ([]*Task, error) {
	return s.repo.ListTasks(ctx, columnID)
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID, userID primitive.ObjectID) error {

	task, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return err
	}
	worspaceId, err := s.boardRepo.GetWorkspaceIDByColumn(ctx, task.ColumnID)
	if err != nil {
		return err
	}

	allowed := s.permissionChecker.HasPermission(
		ctx,
		userID,
		worspaceId,
		permission.PermDeleteTask,
	)

	if !allowed {

		return errors.New("permission denied")
	}

	if err := s.repo.DecrementPositionsInRange(ctx, task.ColumnID, task.Position); err != nil {
		return err
	}
	if err := s.repo.DeleteTask(ctx, taskID); err != nil {
		return err
	}
	s.BroadcastTaskDeleted(ctx, userID, task)
	return nil

}

func (s *TaskService) DragDropTaskToColumn(
	ctx context.Context,
	newPosition int,
	taskID, userID primitive.ObjectID,
	newColumn primitive.ObjectID,
) error {

	task, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return err
	}

	if task == nil {
		return errors.New("task not found")
	}

	worspaceId, err := s.boardRepo.GetWorkspaceIDByColumn(ctx, task.ColumnID)
	if err != nil {
		return err
	}

	allowed := s.permissionChecker.HasPermission(
		ctx,
		userID,
		worspaceId,
		permission.PermEditTask,
	)

	if !allowed {

		return errors.New("permission denied")
	}

	oldPosition := task.Position
	oldColumn := task.ColumnID

	if oldColumn == newColumn {

		if oldPosition < newPosition {
			err = s.repo.ShiftPositions(ctx, oldColumn, oldPosition+1, newPosition, -1)
		} else if oldPosition > newPosition {
			err = s.repo.ShiftPositions(ctx, oldColumn, newPosition, oldPosition-1, 1)
		}

		if err != nil {
			return err
		}

	} else {

		err = s.repo.DecrementPositionsInRange(ctx, oldColumn, oldPosition)
		if err != nil {
			return err
		}

		err = s.repo.IncrementPositionsInRange(ctx, newColumn, newPosition)
		if err != nil {
			return err
		}
	}

	if err := s.repo.UpdateTaskLocation(ctx, taskID, newColumn, newPosition); err != nil {
		return err
	}

	updatedTask, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return err
	}

	s.BroadcastTaskMoved(ctx, userID, updatedTask)

	return nil
}

func (s *TaskService) EditTaskDetails(ctx context.Context, payload *EditTaskPayload, taskid, userid primitive.ObjectID) error {
	task, err := s.repo.GetTaskByID(ctx, taskid)
	if err != nil {
		return err
	}

	if task == nil {
		return errors.New("task not found")
	}

	worspaceId, err := s.boardRepo.GetWorkspaceIDByColumn(ctx, task.ColumnID)
	if err != nil {
		return err
	}

	allowed := s.permissionChecker.HasPermission(
		ctx,
		userid,
		worspaceId,
		permission.PermEditTask,
	)

	if !allowed {

		return errors.New("permission denied")
	}
	now := time.Now()
	assignedMembers := []primitive.ObjectID{}
	for _, id := range payload.AssignedMembers {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		assignedMembers = append(assignedMembers, objectID)
	}
	task = &Task{
		Title:           payload.Title,
		Description:     payload.Description,
		Priority:        payload.Priority,
		AssignedMembers: assignedMembers,
		Labels:          payload.Labels,
		IsArchived:      payload.IsArchived,
		UpdatedAt:       now,
	}

	err = s.repo.UpdateTask(task, ctx, taskid)
	if err != nil {
		return err
	}

	updatedTask, err := s.repo.GetTaskByID(ctx, taskid)
	if err != nil {
		return err
	}
	s.BroadcastTaskUpdated(ctx, userid, updatedTask)

	return nil
}
