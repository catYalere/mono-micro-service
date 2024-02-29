package repositories // import "github.com/catwashere/microservice/internal/database/repositories"

import (
	"context"
	"errors"
	"fmt"
	"github.com/catwashere/microservice/internal/database"
	"github.com/catwashere/microservice/internal/database/utils"
	"github.com/catwashere/microservice/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository[T any] struct {
	Collection *mongo.Collection
}

var (
	SoftDelete     = bson.M{"$set": bson.M{"deleted": true}}
	NotSoftDeleted = bson.M{"deleted": bson.M{"$ne": true}}
)

func newGeneric[T any](_ context.Context, datastore database.IDatabase, collectionName string) (repositories.IRepository[T], error) {
	collection := datastore.GetDatabase().Collection(fmt.Sprintf("%s", collectionName))

	if collection == nil {
		return nil, errors.New("error creating collection") //TODO move to error package
	}

	return Repository[T]{
		Collection: collection,
	}, nil
}

func (r Repository[T]) setID(entity *T, id primitive.ObjectID) error {
	reflection, err := utils.GetReflection[T](entity, "bson", "_id")
	if err != nil {
		return err
		//return errors.New("cannot set Object ID")
	}

	s := id.Hex()
	if id.IsZero() {
		reflection.SetValue(nil)
	} else {
		reflection.SetValue(&s)
	}

	return nil
}

// Get a list of resource
// The function is simply getting all entries in r.collection for the sake of example simplicity
func (r Repository[T]) Get(ctx context.Context, params map[string]interface{}) ([]T, error) {
	filter := []bson.M{
		NotSoftDeleted,
	}

	for k, v := range params {
		filter = append(filter, bson.M{k: v})
	}

	cur, err := r.Collection.Find(ctx, bson.M{"$and": filter})

	if err != nil {
		return nil, err
	}

	var result []T
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var entry T
		if err = cur.Decode(&entry); err != nil {
			return nil, err
		}
		result = append(result, entry)
	}

	return result, nil
}

// GetOne resource based on its ID
func (r Repository[T]) GetOne(ctx context.Context, id string) (T, error) {
	var entry T
	_id, _ := primitive.ObjectIDFromHex(id)
	res := r.Collection.FindOne(ctx, bson.M{"$and": []bson.M{
		bson.M{"_id": _id},
		NotSoftDeleted,
	}})
	err := res.Decode(&entry)
	return entry, err
}

// Create a new resource
func (r Repository[T]) Create(ctx context.Context, entity *T) error {
	result, err := r.Collection.InsertOne(ctx, entity)
	if err != nil {
		return err
	}

	r.setID(entity, result.InsertedID.(primitive.ObjectID))

	return nil
}

// Update a resource
func (r Repository[T]) Update(ctx context.Context, id string, entity *T) error {
	_id, err := primitive.ObjectIDFromHex(id)
	r.setID(entity, primitive.NilObjectID)
	if err != nil {
		return err
	}

	_, err = r.Collection.UpdateOne(ctx, bson.M{"_id": _id}, bson.M{"$set": entity})
	if err != nil {
		return err
	}

	r.setID(entity, _id)

	return nil
}

// Delete a resource, soft delete by marking it as {"deleted": true}
func (r Repository[T]) Delete(ctx context.Context, id string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.Collection.UpdateOne(ctx, bson.M{"_id": _id}, SoftDelete)
	if err != nil {
		return err
	}

	return nil
}
