package internal // import "github.com/catwashere/microservice/internal/usecase/internal"

import (
	"context"
	"crypto/rsa"
	"github.com/catwashere/microservice/internal/entity"
	"github.com/catwashere/microservice/internal/usecase/interfaces"
)

type UseCaseAuth struct {
	user       interfaces.IUseCaseUser
	session    interfaces.IUseCaseSession
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

var _ interfaces.IUseCaseAuth = (*UseCaseAuth)(nil)
var _ interfaces.IUseCaseAuthSetter = (*UseCaseAuth)(nil)

func NewUseCaseAuth(_ context.Context, pk *rsa.PrivateKey, k *rsa.PublicKey) interfaces.IUseCaseAuth {
	return &UseCaseAuth{
		privateKey: pk,
		publicKey:  k,
	}
}

func (sc *UseCaseAuth) SetUseCaseUser(user interfaces.IUseCaseUser) {
	sc.user = user
}

func (sc *UseCaseAuth) SetUseCaseSession(session interfaces.IUseCaseSession) {
	sc.session = session
}

func (sc *UseCaseAuth) Login(ctx context.Context, credentials *entity.Credentials) (*entity.Session, error) {
	user, err := sc.user.ValidateUser(ctx, credentials)
	if err != nil {
		return nil, err
	}

	session := &entity.Session{
		UserID: user.ID,
		User:   user,
	}

	err = sc.session.CreateSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}
