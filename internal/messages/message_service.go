package messages

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageService struct {
	repo *Messagerepo
}

func NewMessageService(repo *Messagerepo) *MessageService {
	return &MessageService{
		repo: repo,
	}
}

func (s *MessageService) CreateGroupChat(ctx context.Context, entityType *CreateGroupChatPayload, entityID primitive.ObjectID) (*GroupChat, error) {

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
		return err, nil
	}
	return nil, gc
}
