package workspace

import (
	"context"
	"errors"
	"time"

	"github.com/Mikhail-tal63/viltrum_empier/internal/permission"
	"go.mongodb.org/mongo-driver/bson"
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

func (s *WorkspaceService) GetWorkspacesUsers(ctx context.Context, workspaceID primitive.ObjectID) ([]bson.M, error) {
	return s.workrepo.GetWorkspaceMembers(ctx, workspaceID)
}
func (s *WorkspaceService) ChangeMembersRole(ctx context.Context, workspaceID, userID, actorID primitive.ObjectID, role *ChangeMembersRolePayload,
) error {
	allowed := s.HasPermission(ctx, actorID, workspaceID, permission.PermManageRoles)
	if !allowed {
		return errors.New("permission denied")
	}
	return s.workrepo.ChangeMembersRole(ctx, workspaceID, userID, role.Role)
}

func (s *WorkspaceService) EditUserPermetions(ctx context.Context, workspaceID, userID, actorID primitive.ObjectID, payload *EditUserPermetionspayload) error {
	allowed := s.HasPermission(ctx, actorID, workspaceID, permission.PermManageRoles)
	if !allowed {
		return errors.New("permission denied")
	}
	return s.workrepo.EditUserPermetions(ctx, workspaceID, userID, payload.Permissions)
}

// permissios*************************************************************************************
func (s *WorkspaceService) HasPermission(
	ctx context.Context,
	userID,
	workspaceID primitive.ObjectID,
	perm string,
) bool {
	member, err := s.workrepo.GetWorkspaceMember(ctx, workspaceID, userID)
	if err != nil || member == nil {
		return false
	}
	if member.Role == "owner" {
		return true
	}
	for _, p := range member.Permissions {
		if p == perm {
			return true
		}
	}

	return false
}
