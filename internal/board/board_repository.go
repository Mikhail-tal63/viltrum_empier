package board

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *BoardRepository) DeleteBoard(ctx context.Context, boardID primitive.ObjectID) error {
	filter := bson.M{"_id": boardID}

	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (r *BoardRepository) GetBoardByID(ctx context.Context, boardID primitive.ObjectID) (*Board, error) {
	filter := bson.M{"_id": boardID}
	var board Board
	err := r.collection.FindOne(ctx, filter).Decode(&board)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &board, err

}

/*culmuns **************************************************************************************************************************************
*************************************************************************************************************************************************
*************************************************************************************************************************************************/

func (r *BoardRepository) CreateColmun(ctx context.Context, colmun *Column) (*Column, error) {
	_, err := r.column.InsertOne(ctx, colmun)
	if err != nil {
		return nil, err
	}
	return colmun, nil
}

func (r *BoardRepository) GetBoardColumns(ctx context.Context, boardID primitive.ObjectID) ([]*Column, error) {
	filter := bson.M{"board_id": boardID}
	cursor, err := r.column.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	columns := []*Column{}
	if err := cursor.All(ctx, &columns); err != nil {
		return nil, err
	}
	return columns, nil
}

func (r *BoardRepository) GetWorkspaceIDByColumn(
	ctx context.Context,
	columnID primitive.ObjectID,
) (primitive.ObjectID, error) {

	var col struct {
		BoardID primitive.ObjectID `bson:"board_id"`
	}
	err := r.column.FindOne(
		ctx,
		bson.M{"_id": columnID},
		options.FindOne().SetProjection(bson.M{"board_id": 1}),
	).Decode(&col)
	if err != nil {
		return primitive.NilObjectID, err
	}

	var b struct {
		WorkspaceID primitive.ObjectID `bson:"workspace_id"`
	}
	err = r.collection.FindOne(
		ctx,
		bson.M{"_id": col.BoardID},
		options.FindOne().SetProjection(bson.M{"workspace_id": 1}),
	).Decode(&b)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return b.WorkspaceID, nil
}

func (r *BoardRepository) CountBoardColumns(
	ctx context.Context,
	boardID primitive.ObjectID,
) (int64, error) {

	filter := bson.M{
		"board_id":    boardID,
		"is_archived": false,
	}

	count, err := r.column.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *BoardRepository) GetColumnByID(ctx context.Context, colmunID primitive.ObjectID) (*Column, error) {
	filter := bson.M{"_id": colmunID}
	var column Column
	err := r.collection.FindOne(ctx, filter).Decode(&column)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &column, nil
}

func (r *BoardRepository) DeleteColumn(ctx context.Context, columnID primitive.ObjectID) error {
	filter := bson.M{"_id": columnID}

	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (r *BoardRepository) IncrementPositionsInRange(ctx context.Context, boardID primitive.ObjectID, fromPosition int) error {
	filter := bson.M{
		"board_id": boardID,
		"position": bson.M{
			"$gte": fromPosition,
		},
	}

	update := bson.M{
		"$inc": bson.M{
			"position": 1,
		},
	}

	_, err := r.column.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
func (r *BoardRepository) ShiftPositions(
	ctx context.Context,
	boardID primitive.ObjectID,
	fromPosition int,
	toPosition int,
	delta int,
) error {

	filter := bson.M{
		"board_id": boardID,
		"position": bson.M{
			"$gte": fromPosition,
			"$lte": toPosition,
		}, "is_archived": false,
	}

	update := bson.M{
		"$inc": bson.M{
			"position": delta,
		},
	}

	_, err := r.collection.UpdateMany(
		ctx,
		filter,
		update,
	)

	return err
}
func (r *BoardRepository) DecrementPositionsInRange(ctx context.Context, boardID primitive.ObjectID, oldPosition int) error {

	filter := bson.M{
		"board_id": boardID,
		"position": bson.M{
			"$gt": oldPosition,
		},
	}

	update := bson.M{
		"$inc": bson.M{
			"position": -1,
		},
	}

	_, err := r.column.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	return err

}

func (r *BoardRepository) UpdateColumnLocation(ctx context.Context, columnID primitive.ObjectID, position int) error {
	filter := bson.M{"_id": columnID}

	update := bson.M{
		"$set": bson.M{
			"position": position,
		},
	}

	_, err := r.column.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return err
}

func (r *BoardRepository) UpdateColumnDetails(ctx context.Context, columnID primitive.ObjectID, fields bson.M) error {
	filter := bson.M{"_id": columnID}

	update := bson.M{
		"$set": fields,
	}

	_, err := r.column.UpdateOne(ctx, filter, update)
	return err
}
