package workspace

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkspaceRepo struct {
	collection  *mongo.Collection
	memberships *mongo.Collection
	invite      *mongo.Collection
}

func NewWorkspaceRepo(db *mongo.Database) *WorkspaceRepo {
	return &WorkspaceRepo{
		collection:  db.Collection("workspaces"),
		memberships: db.Collection("workspace_members"),
		invite:      db.Collection("workspace_invites"),
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

func (r *WorkspaceRepo) CreateWorkspaceMember(member *WorkspaceMember) (*WorkspaceMember, error) {

	res, err := r.memberships.InsertOne(context.TODO(), member)
	if err != nil {
		return nil, err
	}
	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if ok {
		member.ID = insertedID
	}
	return member, nil

}

func (r *WorkspaceRepo) ListWorkspaces(ctx context.Context, userID primitive.ObjectID) ([]*Workspace, error) {

	filter := bson.M{
		"user_id": userID,
	}

	cursor, err := r.memberships.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	memberships := []*WorkspaceMember{}
	if err := cursor.All(ctx, &memberships); err != nil {
		return nil, err
	}

	if len(memberships) == 0 {
		return []*Workspace{}, nil
	}

	workspaceIDs := make([]primitive.ObjectID, 0, len(memberships))

	for _, member := range memberships {
		workspaceIDs = append(workspaceIDs, member.WorkspaceID)
	}

	workspaceFilter := bson.M{
		"_id": bson.M{
			"$in": workspaceIDs,
		},
	}

	workspaceCursor, err := r.collection.Find(ctx, workspaceFilter)
	if err != nil {
		return nil, err
	}
	defer workspaceCursor.Close(ctx)

	workspaces := []*Workspace{}
	if err := workspaceCursor.All(ctx, &workspaces); err != nil {
		return nil, err
	}

	return workspaces, nil
}

func (r *WorkspaceRepo) GetWorkspaceByID(ctx context.Context, workspaceID primitive.ObjectID) (*Workspace, error) {
	filter := bson.M{"_id": workspaceID}
	var workspace Workspace
	err := r.collection.FindOne(ctx, filter).Decode(&workspace)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &workspace, nil

}

func (r *WorkspaceRepo) InviteUser(member *WorkspaceInvite) (*WorkspaceInvite, error) {
	res, err := r.invite.InsertOne(context.TODO(), member)
	if err != nil {
		return nil, err
	}
	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if ok {
		member.ID = insertedID
	}
	return member, nil
}
