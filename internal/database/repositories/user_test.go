package repositories

import (
	"context"
	database "github.com/catwashere/microservice/internal/mocks/database"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func TestNewUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("init", func(mt *mtest.T) {
		datastore := new(database.IDatabase)
		datastore.EXPECT().GetDatabase().Return(mt.DB)

		collection, _ := NewUser(context.Background(), datastore)
		assert.NotNil(t, collection)
	})
}

func TestRepository_GetUserByEmail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	ctx := context.Background()
	collectionName := "tests"
	mongoId := "123456789012345678901234"
	email := "test@gmail.com"

	mt.Run("success", func(mt *mtest.T) {
		datastore := new(database.IDatabase)
		datastore.EXPECT().GetDatabase().Return(mt.DB)

		collection, _ := NewUser(context.Background(), datastore)
		assert.NotNil(t, collection)

		ns := mt.Coll.Database().Name() + "." + collectionName
		mt.AddMockResponses(mtest.CreateCursorResponse(1, ns, mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: mongoId},
				{Key: "email", Value: email},
			}))

		result, _ := collection.GetUserByEmail(ctx, &email)
		assert.NotNil(t, result)
		assert.Equal(t, mongoId, *result.ID)
		assert.Equal(t, email, *result.Email)
	})

	mt.Run("no records", func(mt *mtest.T) {
		datastore := new(database.IDatabase)
		datastore.EXPECT().GetDatabase().Return(mt.DB)

		collection, _ := NewUser(context.Background(), datastore)
		assert.NotNil(t, collection)

		ns := mt.Coll.Database().Name() + "." + collectionName
		mt.AddMockResponses(mtest.CreateCursorResponse(1, ns, mtest.FirstBatch))

		_, err := collection.GetUserByEmail(ctx, &email)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "no responses remaining")
	})

	mt.Run("failure", func(mt *mtest.T) {
		datastore := new(database.IDatabase)
		datastore.EXPECT().GetDatabase().Return(mt.DB)

		collection, _ := NewUser(context.Background(), datastore)
		assert.NotNil(t, collection)

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    200,
			Message: "generic error in mongo",
		}))

		_, err := collection.GetUserByEmail(ctx, &email)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "generic error in mongo")
	})
}
