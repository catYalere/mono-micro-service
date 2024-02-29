package internal // import "github.com/catwashere/microservice/internal/usecase/internal"

import (
	"context"
	"crypto/rsa"
	"errors"
	"github.com/catwashere/microservice/internal/entity"
	"github.com/catwashere/microservice/internal/repositories"
	"github.com/catwashere/microservice/internal/usecase/interfaces"
	"github.com/catwashere/microservice/pkg/crypto"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type UseCaseUser struct {
	repository repositories.UserRepository
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

var _ interfaces.IUseCaseUser = (*UseCaseUser)(nil)

func NewUseCaseUser(_ context.Context, pk *rsa.PrivateKey, k *rsa.PublicKey, repository repositories.UserRepository) interfaces.IUseCaseUser {
	return &UseCaseUser{
		privateKey: pk,
		publicKey:  k,
		repository: repository,
	}
}

func (uc *UseCaseUser) ValidateUser(ctx context.Context, credentials *entity.Credentials) (entity.User, error) {
	user, err := uc.repository.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		return entity.User{}, err
	}

	if uc.privateKey != nil && len(user.BPassword) > 0 {
		c, err := crypto.DecryptWithPrivateKey(user.BPassword, uc.privateKey)
		if err != nil {
			return entity.User{}, err
		}
		cs := string(c)
		user.EncryptedPassword = &cs
		user.BPassword = nil
	}

	if err = bcrypt.CompareHashAndPassword([]byte(*user.EncryptedPassword), []byte(*credentials.Password)); err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (uc *UseCaseUser) GetUser(ctx context.Context, id string) (entity.User, error) {
	user, err := uc.repository.GetOne(ctx, id)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (uc *UseCaseUser) GetUsers(ctx context.Context, params map[string]interface{}) ([]entity.User, error) {
	users, err := uc.repository.Get(ctx, params)
	if err != nil {
		return nil, err
	}

	if users == nil {
		return make([]entity.User, 0, 0), nil
	}

	if params["email"] != nil {
		for i := range users {
			if uc.publicKey != nil {
				msg := []byte(*users[i].EncryptedPassword)
				c, err := crypto.EncryptWithPublicKey(msg, uc.publicKey)
				if err != nil {
					return nil, err
				}
				users[i].BPassword = c
			} else {
				users[i].Password = users[i].EncryptedPassword
			}
		}
	}

	return users, nil
}

func (uc *UseCaseUser) CreateUser(ctx context.Context, user *entity.User) error {
	entries, err := uc.repository.Get(ctx, map[string]interface{}{"email": user.Email})
	if err != nil {
		return err
	}

	if len(entries) > 0 {
		return ErrUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	h := string(hash)
	user.EncryptedPassword = &h
	user.Password = nil

	err = uc.repository.Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UseCaseUser) UpdateUser(ctx context.Context, id string, user *entity.User) error {
	//block if try to update password
	user.Password = nil

	err := uc.repository.Update(ctx, id, user)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UseCaseUser) DeleteUser(ctx context.Context, id string) error {
	err := uc.repository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
