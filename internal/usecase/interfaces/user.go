package interfaces

import (
	"context"
	"github.com/catwashere/microservice/internal/entity"
)

type IUseCaseUser interface {
	GetUser(context.Context, string) (entity.User, error)
	GetUsers(context.Context, map[string]interface{}) ([]entity.User, error)
	CreateUser(context.Context, *entity.User) error
	UpdateUser(context.Context, string, *entity.User) error
	DeleteUser(context.Context, string) error

	ValidateUser(context.Context, *entity.Credentials) (entity.User, error)
}
