package repositories // import "github.com/catwashere/microservice/internal/database/repositories"

import (
	"context"
	"github.com/catwashere/microservice/internal/database"
	"github.com/catwashere/microservice/internal/entity"
	"github.com/catwashere/microservice/internal/repositories"
)

const SessionCollectionName = "sessions"

func NewSession(ctx context.Context, datastore database.IDatabase) (repositories.SessionRepository, error) {
	return newGeneric[entity.Session](ctx, datastore, SessionCollectionName)
}
