package repositories

import (
	"context"
	"github.com/catwashere/microservice/internal/entity"
)

type UserRepository interface {
	IRepository[entity.User]

	GetUserByEmail(context.Context, *string) (entity.User, error)
}
