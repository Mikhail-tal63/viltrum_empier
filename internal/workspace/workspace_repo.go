package workspace

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkspaceRepo struct {
	collection *mongo.Collection
}

func NewWorkspaceRepo(db *mongo.Database) *WorkspaceRepo {
	return &WorkspaceRepo{
		collection: db.Collection("workspaces"),
	}
}

func (r *WorkspaceRepo) CreateWorkspace(workspace *Workspace) (*Workspace, error) {
	res, err := r.collection.InsertOne(context.TODO(), workspace)
	if err != nil {
		return nil, err
	}
	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to convert inserted ID to ObjectID")
	}

	workspace.ID = insertedID
	return workspace, nil
}
