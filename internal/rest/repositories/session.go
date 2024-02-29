package repositories // import "github.com/catwashere/microservice/internal/rest/repositories"

import (
	"context"
	"github.com/catwashere/microservice/internal/entity"
	"github.com/catwashere/microservice/internal/repositories"
)

const SessionPath = "sessions"

func NewSession(ctx context.Context, baseUrl string) (repositories.SessionRepository, error) {
	return newGeneric[entity.Session](ctx, baseUrl, SessionPath)
}
