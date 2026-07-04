package messages

import (
	"context"
	"errors"

	"time"

	"github.com/Mikhail-tal63/viltrum_empier/internal/board"
	"github.com/Mikhail-tal63/viltrum_empier/internal/permission"
	"github.com/Mikhail-tal63/viltrum_empier/internal/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageService struct {
	repo              *Messagerepo
	permissionChecker PermissionChecker
	messageHup        *websocket.Hub
	boardRepo         *board.BoardRepository
}

type PermissionChecker interface {
	HasPermission(
		ctx context.Context,
		userID, workspaceID primitive.ObjectID,
		PermDeleteBoard string,
	) bool
}

func NewMessageService(repo *Messagerepo, permissionChecker PermissionChecker, message *websocket.Hub, boardRepo *board.BoardRepository) *MessageService {
	return &MessageService{
		repo:              repo,
		permissionChecker: permissionChecker,
		messageHup:        message,
		boardRepo:         boardRepo,
	}
}

func (s *MessageService) CreateGroupChat(ctx context.Context, entityType *CreateGroupChatPayload, entityID, userID primitive.ObjectID) (*GroupChat, error) {

	groupChat := &GroupChat{
		ID:         primitive.NewObjectID(),
		EntityType: entityType.EntityType,
		EntityID:   entityID,
		Visibility: VisibilityPublic,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	gc, err := s.repo.CreateGroupChat(ctx, groupChat)
	if err != nil {
		return nil, err
	}

	var workspaceID primitive.ObjectID

	switch entityType.EntityType {
	case EntityWorkspace:
		workspaceID = entityID

	case EntityBoard:
		board, err := s.boardRepo.GetBoardByID(ctx, entityID)
		if err != nil {
			return nil, err
		}
		workspaceID = board.WorkspaceID

	case EntityColumn:
		id, err := s.boardRepo.GetWorkspaceIDByColumn(ctx, entityID)
		if err != nil {
			return nil, err
		}
		workspaceID = id
	default:
		return nil, errors.New("invalid entity type")
	}

	s.BroadcastCreateGroupChat(ctx, userID, workspaceID, gc)

	return gc, nil
}

func (s *MessageService) DeleteGroupChat(ctx context.Context, gID primitive.ObjectID) error {

	if err := s.repo.DeleteGroupChat(ctx, gID); err != nil {
		return err
	}
	return nil
}

func (s *MessageService) CreateMessage(ctx context.Context, payload *CreateMessagePayload, senderId, gID primitive.ObjectID) (*Message, error) {
	message := &Message{
		ID: primitive.NewObjectID(),

		GroupChatID: gID,

		SenderID: senderId,

		Content: payload.Content,

		EditedAt: nil,
		Deleted:  false,

		CreatedAt: time.Now(),
	}

	msg, err := s.repo.CreateMessage(ctx, message)
	if err != nil {
		return nil, err
	}
	return msg, nil

}

func (s *MessageService) DeleteMessage(ctx context.Context, mID, deleterID, workspaceid primitive.ObjectID) error {

	msg, err := s.repo.GetMessageByID(ctx, mID)
	if err != nil {
		return err
	}

	if deleterID != msg.SenderID {

		allawed := s.permissionChecker.HasPermission(ctx, deleterID, workspaceid, permission.PermDeleteMessage)

		if !allawed {
			return errors.New("access denied")
		}
	}

	if err := s.repo.DeleteMessage(ctx, mID); err != nil {
		return err
	}
	return nil
}

func (s *MessageService) ListMessagesInGroupChat(ctx context.Context, gID primitive.ObjectID) ([]*Message, error) {
	messages, err := s.repo.ListMessagesInGroupChat(ctx, gID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *MessageService) EditMessage(ctx context.Context, msgID primitive.ObjectID, payload *EditMessagePayload) error {
	if err := s.repo.EditMessage(ctx, msgID, payload); err != nil {
		return err
	}
	return nil
}
