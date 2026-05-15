package workspace

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkspaceService struct {
	workrepo *WorkspaceRepo
}

func NewWorkspaceService(workrepo *WorkspaceRepo) *WorkspaceService {
	return &WorkspaceService{
		workrepo: workrepo,
	}
}

func (s *WorkspaceService) CreateWorkspace(ctx context.Context, payload *WorkspacePayload, userID primitive.ObjectID) (*Workspace, error) {

	now := time.Now()

	workspaceID := primitive.NewObjectID()

	workspace := &Workspace{
		ID:          workspaceID,
		Name:        payload.name,
		Description: payload.description,
		OwnerID:     userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	ownerMember := &WorkspaceMember{
		ID:          primitive.NewObjectID(),
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        "owner",
		JoinedAt:    now,
	}

	createdWorkspace, err := s.workrepo.CreateWorkspace(workspace)
	if err != nil {
		return nil, err
	}

	_, err = s.workrepo.CreateWorkspaceMember(ownerMember)
	if err != nil {
		return nil, err
	}

	return createdWorkspace, nil
}
