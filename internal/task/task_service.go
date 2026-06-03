package task

import (
	"context"
	"log"
	"time"

	"github.com/Mikhail-tal63/viltrum_empier/internal/board"
	"github.com/Mikhail-tal63/viltrum_empier/internal/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskService struct {
	repo      *TaskRepository
	hub       *websocket.Hub
	boardRepo *board.BoardRepository
}

func NewTaskService(repo *TaskRepository, hub *websocket.Hub, boardRepo *board.BoardRepository) *TaskService {
	return &TaskService{
		repo:      repo,
		hub:       hub,
		boardRepo: boardRepo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, payload *CreateTaskPayload, user, columnID primitive.ObjectID) (*Task, error) {
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

func (s *TaskService) broadcastTaskCreated(
	ctx context.Context,
	columnID primitive.ObjectID,
	task *Task,
	userID primitive.ObjectID,
) {
	workspaceID, err := s.boardRepo.GetWorkspaceIDByColumn(ctx, columnID)
	if err != nil {
		log.Printf("ws: resolve workspace for column %s failed: %v", columnID.Hex(), err)
		return
	}
	payload, err := websocket.MarshalEvent(websocket.EventTaskCreated, map[string]any{
		"task": task,
		"by":   userID.Hex(),
	})
	if err != nil {
		log.Printf("ws: marshal task_created failed: %v", err)
		return
	}
	s.hub.BroadcastToWorkspace(workspaceID.Hex(), payload)
}
