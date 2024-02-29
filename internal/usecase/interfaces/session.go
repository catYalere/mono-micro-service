package interfaces

import (
	"context"
	"github.com/catwashere/microservice/internal/entity"
)

type IUseCaseSession interface {
	GetSession(context.Context, string) (entity.Session, error)
	GetSessions(context.Context, map[string]interface{}) ([]entity.Session, error)
	CreateSession(context.Context, *entity.Session) error
	UpdateSession(context.Context, string, *entity.Session) error
	DeleteSession(context.Context, string) error
}

type IUseCaseSessionSetter interface {
	SetUseCaseUser(IUseCaseUser)
}
