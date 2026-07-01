package messages

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Messagerepo struct {
	msgRepo   *mongo.Collection
	groupRepo *mongo.Collection
}

func NewMessageRepository(db *mongo.Database) *Messagerepo {
	return &Messagerepo{
		msgRepo:   db.Collection("message"),
		groupRepo: db.Collection("group_chat"),
	}
}

func (r *Messagerepo) CreateGroupChat(ctx context.Context, Group *GroupChat) (error, *GroupChat) {
	_, err := r.groupRepo.InsertOne(ctx, Group)
	if err != nil {
		return err, nil
	}
	return nil, Group
}
func (r *Messagerepo) DeleteGroupChat(ctx context.Context, gID primitive.ObjectID) error {
	filter := bson.M{"_id": gID}

	_, err := r.groupRepo.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (r *Messagerepo) GetGroupChatByEntityID(ctx context.Context, entityId primitive.ObjectID) (*GroupChat, error) {
	filter := bson.M{
		"entityId": entityId,
	}
	var ch GroupChat

	if err := r.groupRepo.FindOne(ctx, filter).Decode(&ch); err != nil {
		return nil, err
	}
	return &ch, nil
}
