package workspace

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Workspace struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`

	OwnerID primitive.ObjectID `bson:"owner_id" json:"owner_id"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type WorkspaceMember struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	WorkspaceID primitive.ObjectID `bson:"workspace_id" json:"workspace_id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`

	Role string `bson:"role" json:"role"`

	JoinedAt time.Time `bson:"joined_at" json:"joined_at"`
}

type WorkspacePayload struct {
	name        string
	description string
}
