package task

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskService struct {
	repo *TaskRepository
}

func NewTaskService(repo *TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
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
	return createdTask, nil
}
