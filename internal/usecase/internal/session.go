package internal // import "github.com/catwashere/microservice/internal/usecase/internal"

import (
	"context"
	"crypto/rsa"
	"github.com/catwashere/microservice/internal/entity"
	"github.com/catwashere/microservice/internal/repositories"
	"github.com/catwashere/microservice/internal/usecase/interfaces"
)

type UseCaseSession struct {
	repository repositories.SessionRepository
	user       interfaces.IUseCaseUser
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

var _ interfaces.IUseCaseSession = (*UseCaseSession)(nil)
var _ interfaces.IUseCaseSessionSetter = (*UseCaseSession)(nil)

func NewUseCaseSession(_ context.Context, pk *rsa.PrivateKey, k *rsa.PublicKey, repository repositories.SessionRepository) interfaces.IUseCaseSession {
	return &UseCaseSession{
		privateKey: pk,
		publicKey:  k,
		repository: repository,
	}
}

func (sc *UseCaseSession) SetUseCaseUser(user interfaces.IUseCaseUser) {
	sc.user = user
}

func (uc *UseCaseSession) GetSession(ctx context.Context, id string) (entity.Session, error) {
	session, err := uc.repository.GetOne(ctx, id)
	if err != nil {
		return entity.Session{}, err
	}

	if session.User.ID == nil {
		user, err := uc.user.GetUser(ctx, *session.UserID)
		if err != nil {
			return entity.Session{}, err
		}
		session.User = user
	}

	return session, nil
}

func (uc *UseCaseSession) GetSessions(ctx context.Context, params map[string]interface{}) ([]entity.Session, error) {
	sessions, err := uc.repository.Get(ctx, params)
	if err != nil {
		return nil, err
	}

	if sessions == nil {
		return make([]entity.Session, 0, 0), nil
	}

	for i, session := range sessions {
		if session.User.ID == nil {
			user, err := uc.user.GetUser(ctx, *session.UserID)
			if err != nil {
				return nil, err
			}
			sessions[i].User = user
		}
	}
	return sessions, nil
}

func (uc *UseCaseSession) CreateSession(ctx context.Context, session *entity.Session) error {
	// Block if not user.id
	if session.UserID == nil && session.User.ID != nil {
		session.UserID = session.User.ID
	}

	err := uc.repository.Create(ctx, session)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UseCaseSession) UpdateSession(ctx context.Context, id string, session *entity.Session) error {
	// Not sure if needed
	err := uc.repository.Update(ctx, id, session)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UseCaseSession) DeleteSession(ctx context.Context, id string) error {
	err := uc.repository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
