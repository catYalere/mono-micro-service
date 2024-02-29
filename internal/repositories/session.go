package repositories

import "github.com/catwashere/microservice/internal/entity"

type SessionRepository interface {
	IRepository[entity.Session]
}
