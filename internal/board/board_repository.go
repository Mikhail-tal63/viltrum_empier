package board

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type BoardRepository struct {
	collection *mongo.Collection
}

func NewBoardRepository(db *mongo.Database) *BoardRepository {
	return &BoardRepository{
		collection: db.Collection("boards"),
	}
}

func (r *BoardRepository) CreateBoard(board *Board) (*Board, error) {
	_, err := r.collection.InsertOne(context.TODO(), board)
	if err != nil {
		return nil, err
	}
	return board, nil
}
