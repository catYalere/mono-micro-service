package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/catwashere/microservice/internal/entity"
	"github.com/catwashere/microservice/internal/repositories"
	"io"
	"net/http"
)

const UserPath = "users"

func NewUser(ctx context.Context, baseUrl string) (repositories.UserRepository, error) {
	uc, err := newGeneric[entity.User](ctx, baseUrl, UserPath)
	if err != nil {
		return nil, err
	}

	return uc.(Repository[entity.User]), nil
}

func (r Repository[T]) GetUserByEmail(ctx context.Context, email *string) (entity.User, error) {
	var entry entity.User
	resp, err := http.Get(fmt.Sprintf("%s/%s?email=%s", r.BaseUrl, r.Path, *email))
	if err != nil {
		return entry, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return entry, err
	}

	var entries []entity.User
	err = json.Unmarshal(body, &entries)
	if err != nil {
		return entry, err
	}

	if len(entries) == 0 {
		return entry, fmt.Errorf("no user found with email %s", *email)
	}

	entry = entries[0]
	entry.EncryptedPassword = entry.Password
	entry.Password = nil

	return entry, nil
}
