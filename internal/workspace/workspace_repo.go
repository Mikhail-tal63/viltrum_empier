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
func (r *WorkspaceRepo) GetWorkspaceMember(
	ctx context.Context,
	workspaceID,
	userID primitive.ObjectID,
) (*WorkspaceMember, error) {

	filter := bson.M{
		"workspace_id": workspaceID,
		"user_id":      userID,
	}

	var member WorkspaceMember

	err := r.memberships.FindOne(ctx, filter).Decode(&member)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &member, nil
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

func (r *WorkspaceRepo) GetInviteByUserID(ctx context.Context, UserID primitive.ObjectID) ([]*WorkspaceInvite, error) {
	filter := bson.M{
		"invited_user_id": UserID,
		"status":          "pending",
	}

	cursor, err := r.invite.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	invites := []*WorkspaceInvite{}
	if err := cursor.All(ctx, &invites); err != nil {
		return nil, err
	}
	return invites, nil
}

func (r *WorkspaceRepo) AcceptInvite(ctx context.Context, inviteID primitive.ObjectID) error {
	filter := bson.M{"_id": inviteID, "status": "pending"}
	update := bson.M{"$set": bson.M{"status": "accepted"}}

	_, err := r.invite.UpdateOne(ctx, filter, update)
	return err

}
func (r *WorkspaceRepo) GetWorkspaceMembers(ctx context.Context, workspaceID primitive.ObjectID) ([]bson.M, error) {

	pipeline := mongo.Pipeline{

		{
			{"$match", bson.D{
				{"workspace_id", workspaceID},
			}},
		},

		{
			{"$lookup", bson.D{
				{"from", "users"},
				{"localField", "user_id"},
				{"foreignField", "_id"},
				{"as", "user"},
			}},
		},

		{
			{"$unwind", "$user"},
		},

		{
			{"$project", bson.D{
				{"_id", 1},
				{"role", 1},
				{"workspace_id", 1},
				{"user_id", 1},

				{"username", "$user.username"},
				{"email", "$user.email"},
			}},
		},
	}

	cursor, err := r.memberships.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var results []bson.M

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
