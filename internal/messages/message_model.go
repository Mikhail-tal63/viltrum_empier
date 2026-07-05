package messages

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupChat struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	EntityType ChatEntity         `bson:"entityType" json:"entityType"`
	EntityID   primitive.ObjectID `bson:"entityId" json:"entityId"`

	Visibility ChatVisibility `bson:"visibility" json:"visibility"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type Message struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	GroupChatID primitive.ObjectID `bson:"groupChatId" json:"groupChatId"`
	SenderID    primitive.ObjectID `bson:"senderId" json:"senderId"`

	Content string `bson:"content" json:"content"`

	EditedAt *time.Time `bson:"editedAt,omitempty" json:"editedAt,omitempty"`
	Deleted  bool       `bson:"deleted" json:"deleted"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

type ChatEntity string

const (
	EntityWorkspace ChatEntity = "workspace"
	EntityBoard     ChatEntity = "board"
	EntityColumn    ChatEntity = "column"
)

type ChatVisibility string

const (
	VisibilityPublic ChatVisibility = "public"
	VisibilityHidden ChatVisibility = "hidden"
)
