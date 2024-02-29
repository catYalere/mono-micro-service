package repositories

import (
	"context"
	database "github.com/catwashere/microservice/internal/mocks/database"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

type Something struct {
	ID        *string `bson:"_id,omitempty"`
	Something *string `bson:"something,omitempty"`
}

func Test_newGeneric(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	collectionName := "tests"

	mt.Run("init", func(mt *mtest.T) {
		datastore := new(database.IDatabase)
		datastore.EXPECT().GetDatabase().Return(mt.DB)

		collection, _ := newGeneric[Something](context.Background(), datastore, collectionName)
		assert.NotNil(t, collection)
	})
}

func TestRepository_Get(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	ctx := context.Background()
	collectionName := "tests"
	mongoId := "123456789012345678901234"
	something := "something"

	mt.Run("success", func(mt *mtest.T) {
		datastore := new(database.IDatabase)
		datastore.EXPECT().GetDatabase().Return(mt.DB)

		collection, _ := newGeneric[Something](context.Background(), datastore, collectionName)
		assert.NotNil(t, collection)

		ns := mt.Coll.Database().Name() + "." + collectionName
		mt.AddMockResponses(mtest.CreateCursorResponse(1, ns, mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: mongoId + "1"},
				{Key: "something", Value: something + "1"},
			}, bson.D{
				{Key: "_id", Value: mongoId + "2"},
				{Key: "something", Value: something + "2"},
			}))

		results, _ := collection.Get(ctx, map[string]interface{}{"id": "1234"})
		assert.Len(t, results, 2)
		assert.Equal(t, mongoId+"1", *results[0].ID)
		assert.Equal(t, something+"1", *results[0].Something)
		assert.Equal(t, mongoId+"2", *results[1].ID)
		assert.Equal(t, something+"2", *results[1].Something)
	})

	mt.Run("no records", func(mt *mtest.T) {
		datastore := new(database.IDatabase)
		datastore.EXPECT().GetDatabase().Return(mt.DB)

		collection, _ := newGeneric[Something](context.Background(), datastore, collectionName)
		assert.NotNil(t, collection)

		ns := mt.Coll.Database().Name() + "." + collectionName
		mt.AddMockResponses(mtest.CreateCursorResponse(1, ns, mtest.FirstBatch))

		results, err := collection.Get(ctx, map[string]interface{}{"id": "1234"})
		assert.Len(t, results, 0)
		assert.Nil(t, err)
	})

	mt.Run("failure", func(mt *mtest.T) {
		datastore := new(database.IDatabase)
		datastore.EXPECT().GetDatabase().Return(mt.DB)

		collection, _ := newGeneric[Something](context.Background(), datastore, collectionName)
		assert.NotNil(t, collection)

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    200,
			Message: "generic error in mongo",
		}))

		results, err := collection.Get(ctx, map[string]interface{}{"id": "1234"})
		assert.Len(t, results, 0)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "generic error in mongo")
	})
}

func TestRepository_GetOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	ctx := context.Background()
	collectionName := "tests"
	mongoId := "123456789012345678901234"
	something := "something"

	mt.Run("success", func(mt *mtest.T) {
		datastore := new(database.IDatabase)
		datastore.EXPECT().GetDatabase().Return(mt.DB)

		collection, _ := newGeneric[Something](context.Background(), datastore, collectionName)
		assert.NotNil(t, collection)

		ns := mt.Coll.Database().Name() + "." + collectionName
		mt.AddMockResponses(mtest.CreateCursorResponse(1, ns, mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: mongoId},
				{Key: "something", Value: something},
			}))

		result, _ := collection.GetOne(ctx, mongoId)
		assert.NotNil(t, result)
		assert.Equal(t, mongoId, *result.ID)
		assert.Equal(t, something, *result.Something)
	})

	mt.Run("no records", func(mt *mtest.T) {
		datastore := new(database.IDatabase)
		datastore.EXPECT().GetDatabase().Return(mt.DB)

		collection, _ := newGeneric[Something](context.Background(), datastore, collectionName)
		assert.NotNil(t, collection)

		ns := mt.Coll.Database().Name() + "." + collectionName
		mt.AddMockResponses(mtest.CreateCursorResponse(1, ns, mtest.FirstBatch))

		_, err := collection.GetOne(ctx, mongoId)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "no responses remaining")
	})

	mt.Run("failure", func(mt *mtest.T) {
		datastore := new(database.IDatabase)
		datastore.EXPECT().GetDatabase().Return(mt.DB)

		collection, _ := newGeneric[Something](context.Background(), datastore, collectionName)
		assert.NotNil(t, collection)

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    200,
			Message: "generic error in mongo",
		}))

		_, err := collection.GetOne(ctx, mongoId)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "generic error in mongo")
	})
}

func TestRepository_Create(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	ctx := context.Background()
	collectionName := "tests"
	something := "something"
	obj := &Something{
		Something: &something,
	}

	mt.Run("success", func(mt *mtest.T) {
		datastore := new(database.IDatabase)
		datastore.EXPECT().GetDatabase().Return(mt.DB)

		collection, _ := newGeneric[Something](context.Background(), datastore, collectionName)
		assert.NotNil(t, collection)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		err := collection.Create(ctx, obj)
		assert.Nil(t, err)
		assert.Len(t, *obj.ID, 24)
		assert.Equal(t, something, *obj.Something)
	})

	mt.Run("failure", func(mt *mtest.T) {
		datastore := new(database.IDatabase)
		datastore.EXPECT().GetDatabase().Return(mt.DB)

		collection, _ := newGeneric[Something](context.Background(), datastore, collectionName)
		assert.NotNil(t, collection)

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    200,
			Message: "generic error in mongo",
		}))
		err := collection.Create(ctx, obj)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "generic error in mongo")
	})
}

func TestRepository_Update(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	ctx := context.Background()
	collectionName := "tests"
	mongoId := "123456789012345678901234"
	something := "something"
	obj := &Something{
		Something: &something,
	}

	mt.Run("success", func(mt *mtest.T) {
		datastore := new(database.IDatabase)
		datastore.EXPECT().GetDatabase().Return(mt.DB)

		collection, _ := newGeneric[Something](context.Background(), datastore, collectionName)
		assert.NotNil(t, collection)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		err := collection.Update(ctx, mongoId, obj)
		assert.Nil(t, err)
		assert.Len(t, *obj.ID, 24)
		assert.Equal(t, mongoId, *obj.ID)
		assert.Equal(t, something, *obj.Something)
	})

	mt.Run("failure", func(mt *mtest.T) {
		datastore := new(database.IDatabase)
		datastore.EXPECT().GetDatabase().Return(mt.DB)

		collection, _ := newGeneric[Something](context.Background(), datastore, collectionName)
		assert.NotNil(t, collection)

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    200,
			Message: "generic error in mongo",
		}))
		err := collection.Update(ctx, mongoId, obj)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "generic error in mongo")
	})
}

func TestRepository_Delete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	ctx := context.Background()
	collectionName := "tests"
	mongoId := "123456789012345678901234"

	mt.Run("success", func(mt *mtest.T) {
		datastore := new(database.IDatabase)
		datastore.EXPECT().GetDatabase().Return(mt.DB)

		collection, _ := newGeneric[Something](context.Background(), datastore, collectionName)
		assert.NotNil(t, collection)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		err := collection.Delete(ctx, mongoId)
		assert.Nil(t, err)
	})

	mt.Run("failure", func(mt *mtest.T) {
		datastore := new(database.IDatabase)
		datastore.EXPECT().GetDatabase().Return(mt.DB)

		collection, _ := newGeneric[Something](context.Background(), datastore, collectionName)
		assert.NotNil(t, collection)

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    200,
			Message: "generic error in mongo",
		}))
		err := collection.Delete(ctx, mongoId)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "generic error in mongo")
	})
}
