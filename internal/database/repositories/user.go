package repositories // import "github.com/catwashere/microservice/internal/database/repositories"

import (
	"context"
	"github.com/catwashere/microservice/internal/database"
	"github.com/catwashere/microservice/internal/entity"
	"github.com/catwashere/microservice/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
)

const UserCollectionName = "users"

func NewUser(ctx context.Context, datastore database.IDatabase) (repositories.UserRepository, error) {
	uc, err := newGeneric[entity.User](ctx, datastore, UserCollectionName)
	if err != nil {
		return nil, err
	}

	return uc.(Repository[entity.User]), nil
}

func (r Repository[T]) GetUserByEmail(ctx context.Context, email *string) (entity.User, error) {
	var entry entity.User
	res := r.Collection.FindOne(ctx, bson.M{"$and": []bson.M{
		bson.M{"email": email},
		NotSoftDeleted,
	}})
	err := res.Decode(&entry)
	return entry, err
}
