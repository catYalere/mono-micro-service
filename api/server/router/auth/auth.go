package auth // import "github.com/catwashere/microservice/api/server/router/auth"

import (
	"context"
	"encoding/json"
	"github.com/catwashere/microservice/api/server/httputils"
	"github.com/catwashere/microservice/internal/entity"
	"net/http"
)

func (ur *authRouter) login(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var credentials entity.Credentials

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		httputils.WriteJSON(w, http.StatusBadRequest, err.Error())
		return nil
	}

	auth, err := ur.useCase.Login(ctx, &credentials)
	if err != nil {
		return err
	}

	_ = httputils.WriteJSON(w, http.StatusOK, auth)

	return nil
}
