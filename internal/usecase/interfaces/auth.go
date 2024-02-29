package interfaces

import (
	"context"
	"github.com/catwashere/microservice/internal/entity"
)

type IUseCaseAuth interface {
	Login(context.Context, *entity.Credentials) (*entity.Session, error)
}

type IUseCaseAuthSetter interface {
	SetUseCaseUser(IUseCaseUser)
	SetUseCaseSession(IUseCaseSession)
}
