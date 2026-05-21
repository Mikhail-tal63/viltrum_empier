package task

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) *TaskRepository {
	return &TaskRepository{
		collection: db.Collection("tasks"),
	}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task *Task) (*Task, error) {
	res, err := r.collection.InsertOne(ctx, task)
	if err != nil {
		return nil, err
	}
	taskID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to convert inserted ID to ObjectID")
	}
	task.ID = taskID
	return task, nil
}
func (r *TaskRepository) CountTaskColumns(
	ctx context.Context,
	colmunID primitive.ObjectID,
) (int64, error) {

	filter := bson.M{
		"column_id":   colmunID,
		"is_archived": false,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}
