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
		Name:        payload.Name,
		Description: payload.Description,
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

func (s *WorkspaceService) ListWorkspaces(ctx context.Context, userID primitive.ObjectID) ([]*Workspace, error) {
	return s.workrepo.ListWorkspaces(ctx, userID)
}

func (s *WorkspaceService) InviteUser(
	ctx context.Context,
	invitedBy primitive.ObjectID,
	invitedUserID primitive.ObjectID,
	workspaceID primitive.ObjectID,
) (*WorkspaceInvite, error) {
	now := time.Now()
	inviteID := primitive.NewObjectID()
	invite := &WorkspaceInvite{
		ID: inviteID,

		WorkspaceID: workspaceID,

		InvitedBy: invitedBy,

		InvitedUserID: invitedUserID,

		Status: "pending",

		CreatedAt: now,
	}
	created, err := s.workrepo.InviteUser(invite)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *WorkspaceService) GetInviteByUserID(ctx context.Context, UserID primitive.ObjectID) ([]*WorkspaceInvite, error) {
	return s.workrepo.GetInviteByUserID(ctx, UserID)
}

func (s *WorkspaceService) AcceptInvite(ctx context.Context, workspaceID, userID, inviteID primitive.ObjectID) (*WorkspaceMember, error) {
	err := s.workrepo.AcceptInvite(ctx, inviteID)
	if err != nil {
		return nil, err
	}

	NewMember := &WorkspaceMember{
		ID:          primitive.NewObjectID(),
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        "member",
		JoinedAt:    time.Now(),
	}
	_, err = s.workrepo.CreateWorkspaceMember(NewMember)

	if err != nil {
		return nil, err
	}
	return NewMember, nil

}
