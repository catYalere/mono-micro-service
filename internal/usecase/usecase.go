package usecase // import "github.com/catwashere/microservice/internal/usecase"

import (
	"context"
	"crypto/rsa"
	"github.com/catwashere/microservice/internal/database"
	dbrepositories "github.com/catwashere/microservice/internal/database/repositories"
	"github.com/catwashere/microservice/internal/repositories"
	restrepositories "github.com/catwashere/microservice/internal/rest/repositories"
	"github.com/catwashere/microservice/internal/usecase/interfaces"
	"github.com/catwashere/microservice/internal/usecase/internal"
)

type Type int

const (
	None Type = iota
	Internal
	External
)

type EntityConfig struct {
	Type       Type
	BaseUrl    string
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// Config provides the configuration for the Database Connection
type Config struct {
	DB      *database.Database
	User    EntityConfig
	Session EntityConfig
	Auth    EntityConfig
}

type UseCases struct {
	user    interfaces.IUseCaseUser
	session interfaces.IUseCaseSession
	auth    interfaces.IUseCaseAuth
}

// New returns a new usecase instance
func New(ctx context.Context, cfg *Config) *UseCases {

	// Create the usecase instance
	uc := &UseCases{}

	if cfg.User.Type != None {
		var r repositories.UserRepository

		if cfg.User.Type == Internal {
			// Create the internal user usecase
			r, _ = dbrepositories.NewUser(ctx, cfg.DB)
		} else if cfg.User.Type == External {
			// Create the external user usecase
			r, _ = restrepositories.NewUser(ctx, cfg.User.BaseUrl)
		}

		uc.user = internal.NewUseCaseUser(ctx, cfg.User.PrivateKey, cfg.User.PublicKey, r)
	}

	if cfg.Session.Type != None {
		var r repositories.SessionRepository

		if cfg.Session.Type == Internal {
			// Create the internal session usecase
			r, _ = dbrepositories.NewSession(ctx, cfg.DB)
		} else if cfg.Session.Type == External {
			// Create the external session usecase
			r, _ = restrepositories.NewSession(ctx, cfg.Session.BaseUrl)
		}

		uc.session = internal.NewUseCaseSession(ctx, cfg.Session.PrivateKey, cfg.Session.PublicKey, r)

		if cfg.User.Type != None {
			uc.session.(interfaces.IUseCaseSessionSetter).SetUseCaseUser(uc.user)
		}
	}

	if cfg.Auth.Type != None {
		uc.auth = internal.NewUseCaseAuth(ctx, cfg.Auth.PrivateKey, cfg.Auth.PublicKey)

		if cfg.User.Type != None {
			uc.auth.(interfaces.IUseCaseAuthSetter).SetUseCaseUser(uc.user)
		}

		if cfg.Session.Type != None {
			uc.auth.(interfaces.IUseCaseAuthSetter).SetUseCaseSession(uc.session)
		}
	}

	// Return the usecase
	return uc
}

func (u *UseCases) User() interfaces.IUseCaseUser {
	return u.user
}

func (u *UseCases) Session() interfaces.IUseCaseSession {
	return u.session
}

func (u *UseCases) Auth() interfaces.IUseCaseAuth {
	return u.auth
}
