package task

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *TaskRepository) ListTasks(
	ctx context.Context,
	columnID primitive.ObjectID,
) ([]*Task, error) {

	filter := bson.M{
		"column_id":   columnID,
		"is_archived": false,
	}

	opt := options.Find().SetSort(
		bson.M{"position": 1},
	)

	cursor, err := r.collection.Find(ctx, filter, opt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	tasks := []*Task{}

	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) DeleteTask(ctx context.Context, taskID primitive.ObjectID) error {
	firlter := bson.M{"_id": taskID}

	_, err := r.collection.DeleteOne(ctx, firlter)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) UpdateTask(task *Task, ctx context.Context, taskID primitive.ObjectID) error {
	filter := bson.M{"_id": taskID}

	update := bson.M{
		"$set": bson.M{
			"title":            task.Title,
			"description":      task.Description,
			"priority":         task.Priority,
			"labels":           task.Labels,
			"assigned_members": task.AssignedMembers,
			"is_archived":      task.IsArchived,
			"updated_at":       task.UpdatedAt,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) UpdateTaskLocation(
	ctx context.Context,
	taskID primitive.ObjectID,
	columnID primitive.ObjectID,
	position int,
) error {

	filter := bson.M{
		"_id": taskID,
	}

	update := bson.M{
		"$set": bson.M{
			"column_id": columnID,
			"position":  position,
		},
	}

	_, err := r.collection.UpdateOne(
		ctx,
		filter,
		update,
	)

	return err
}

func (r *TaskRepository) DecrementPositionsInRange(
	ctx context.Context,
	columnID primitive.ObjectID,
	oldPosition int,
) error {

	filter := bson.M{
		"column_id": columnID,
		"position": bson.M{
			"$gt": oldPosition,
		},
		"is_archived": false,
	}

	update := bson.M{
		"$inc": bson.M{
			"position": -1,
		},
	}

	_, err := r.collection.UpdateMany(
		ctx,
		filter,
		update,
	)

	return err
}

func (r *TaskRepository) IncrementPositionsInRange(
	ctx context.Context,
	columnID primitive.ObjectID,
	fromPosition int,
) error {

	filter := bson.M{
		"column_id": columnID,
		"position": bson.M{
			"$gte": fromPosition,
		},
		"is_archived": false,
	}

	update := bson.M{
		"$inc": bson.M{
			"position": 1,
		},
	}

	_, err := r.collection.UpdateMany(
		ctx,
		filter,
		update,
	)

	return err
}

func (r *TaskRepository) ShiftPositions(
	ctx context.Context,
	columnID primitive.ObjectID,
	fromPosition int,
	toPosition int,
	delta int,
) error {

	filter := bson.M{
		"column_id": columnID,
		"position": bson.M{
			"$gte": fromPosition,
			"$lte": toPosition,
		},
		"is_archived": false,
	}

	update := bson.M{
		"$inc": bson.M{
			"position": delta,
		},
	}

	_, err := r.collection.UpdateMany(
		ctx,
		filter,
		update,
	)

	return err
}

func (r *TaskRepository) GetTaskByID(ctx context.Context, taskid primitive.ObjectID) (*Task, error) {
	filter := bson.M{"_id": taskid}
	var task Task
	err := r.collection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}
