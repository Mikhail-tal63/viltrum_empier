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

func (r *Messagerepo) CreateGroupChat(ctx context.Context, Group *GroupChat) (*GroupChat, error) {
	_, err := r.groupRepo.InsertOne(ctx, Group)
	if err != nil {
		return nil, err
	}
	return Group, nil
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

func (r *Messagerepo) CreateMessage(ctx context.Context, message *Message) (*Message, error) {
	_, err := r.msgRepo.InsertOne(ctx, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}
func (r *Messagerepo) DeleteMessage(ctx context.Context, msgID primitive.ObjectID) error {
	filter := bson.M{
		"_id": msgID,
	}
	_, err := r.msgRepo.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (r *Messagerepo) ListMessagesInGroupChat(ctx context.Context, gcID primitive.ObjectID) ([]*Message, error) {
	filter := bson.M{
		"groupChatId": gcID,
	}
	curser, err := r.msgRepo.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer curser.Close(ctx)
	messages := []*Message{}

	if err := curser.All(ctx, &messages); err != nil {
		return nil, err
	}
	return messages, nil

}

func (r *Messagerepo) EditMessage(ctx context.Context, msgID primitive.ObjectID, payload *EditMessagePayload) error {
	filter := bson.M{
		"_id": msgID,
	}

	update := bson.M{
		"$set": bson.M{
			"content": payload.Content,
		},
	}

	_, err := r.msgRepo.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
func (r *Messagerepo) GetMessageByID(ctx context.Context, msgID primitive.ObjectID) (*Message, error) {
	filter := bson.M{"_id": msgID}
	var msg Message
	err := r.msgRepo.FindOne(ctx, filter).Decode(&msg)
	if err != nil {

		return nil, err
	}
	return &msg, nil
}
