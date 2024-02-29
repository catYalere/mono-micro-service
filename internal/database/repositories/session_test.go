package repositories

import (
	"context"
	database "github.com/catwashere/microservice/internal/mocks/database"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func TestNewSession(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("init", func(mt *mtest.T) {
		datastore := new(database.IDatabase)
		datastore.EXPECT().GetDatabase().Return(mt.DB)

		collection, _ := NewSession(context.Background(), datastore)
		assert.NotNil(t, collection)
	})
}
