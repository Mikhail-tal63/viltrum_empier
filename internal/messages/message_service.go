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

func (s *MessageService) CreateWorkspaceGroupChat(ctx context.Context, workspaceID primitive.ObjectID) error {
	payload := &CreateGroupChatPayload{EntityType: EntityWorkspace}
	_, err := s.CreateGroupChat(ctx, payload, workspaceID, primitive.NilObjectID)
	return err
}

func (s *MessageService) CreateBoardGroupChat(ctx context.Context, boardID primitive.ObjectID) error {
	payload := &CreateGroupChatPayload{EntityType: EntityBoard}
	_, err := s.CreateGroupChat(ctx, payload, boardID, primitive.NilObjectID)
	return err
}

func (s *MessageService) CreateColumnGroupChat(ctx context.Context, columnID primitive.ObjectID) error {
	payload := &CreateGroupChatPayload{EntityType: EntityColumn}
	_, err := s.CreateGroupChat(ctx, payload, columnID, primitive.NilObjectID)
	return err
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

	s.BroadcastMessageCreated(ctx, senderId, gID, msg)

	return msg, nil

}

func (s *MessageService) DeleteMessage(ctx context.Context, mID, deleterID primitive.ObjectID) error {



	msg, err := s.repo.GetMessageByID(ctx, mID)
	if err != nil {
		return err
	}



	gID := msg.GroupChatID

	gc,err:= s.repo.GetGroupChatByID(ctx,gID)
	
	var workspaceID primitive.ObjectID

	switch gc.EntityType {
	case EntityWorkspace:
		workspaceID = gc.EntityID

	case EntityBoard:
		board, err := s.boardRepo.GetBoardByID(ctx, gc.EntityID)
		if err != nil {
			return  err
		}
		workspaceID = board.WorkspaceID

	case EntityColumn:
		id, err := s.boardRepo.GetWorkspaceIDByColumn(ctx, gc.EntityID)
		if err != nil {
			return err
		}
		workspaceID = id
	default:
		return errors.New("invalid entity type")
	}

	if deleterID != msg.SenderID {

		allawed := s.permissionChecker.HasPermission(ctx, deleterID, workspaceID, permission.PermDeleteMessage)

		if !allawed {
			return errors.New("access denied")
		}
	}

	if err := s.repo.DeleteMessage(ctx, mID); err != nil {
		return err
	}
	s.BroadcastMessageDeleted(ctx, deleterID, gID, mID)
	return nil
}

func (s *MessageService) ListMessagesInGroupChat(ctx context.Context, gID primitive.ObjectID) ([]*Message, error) {
	messages, err := s.repo.ListMessagesInGroupChat(ctx, gID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *MessageService) EditMessage(ctx context.Context, msgID,userID primitive.ObjectID, payload *EditMessagePayload) error {

	if err := s.repo.EditMessage(ctx, msgID, payload); err != nil {
		return err
	}
	message, err := s.repo.GetMessageByID(ctx, msgID)
	if err != nil {
		return err
	}
	gID := message.GroupChatID
	userid := message.SenderID
	if userID!= userid{
		return errors.New("access denied")
	}
	s.BroadcastMessageEdited(ctx, userid, gID, message)
	return nil
}
