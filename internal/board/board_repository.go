package board

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BoardRepository struct {
	collection *mongo.Collection
	column     *mongo.Collection
}

func NewBoardRepository(db *mongo.Database) *BoardRepository {
	return &BoardRepository{
		collection: db.Collection("boards"),
		column:     db.Collection("columns"),
	}
}

func (r *BoardRepository) CreateBoard(ctx context.Context, board *Board) (*Board, error) {
	_, err := r.collection.InsertOne(ctx, board)
	if err != nil {
		return nil, err
	}
	return board, nil
}

func (r *BoardRepository) GetWorkspaceBoards(ctx context.Context, workspaceID primitive.ObjectID) ([]*Board, error) {
	filter := bson.M{"workspace_id": workspaceID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	boards := []*Board{}
	if err := cursor.All(ctx, &boards); err != nil {
		return nil, err
	}
	return boards, nil
}

/*culmuns ****/

func (r *BoardRepository) CreateColmun(ctx context.Context, colmun *Column) (*Column, error) {
	_, err := r.collection.InsertOne(ctx, colmun)
	if err != nil {
		return nil, err
	}
	return colmun, nil
}
