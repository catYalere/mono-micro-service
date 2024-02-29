package user // import "github.com/catwashere/microservice/api/server/router/user"

import (
	"context"
	"encoding/json"
	"github.com/catwashere/microservice/api/server/httputils"
	"github.com/catwashere/microservice/internal/entity"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (ur *userRouter) getUser(ctx context.Context, w http.ResponseWriter, _ *http.Request, vars map[string]string) error {
	user, err := ur.useCase.GetUser(ctx, vars["id"])
	if err != nil {
		return err
	}

	_ = httputils.WriteJSON(w, http.StatusOK, user)
	return nil
}

func (ur *userRouter) getUsers(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	values := r.URL.Query()
	params := make(map[string]interface{}, len(values))
	for k, v := range values {
		params[k] = v[0]
	}

	users, err := ur.useCase.GetUsers(ctx, params)
	if err != nil {
		return err
	}

	_ = httputils.WriteJSON(w, http.StatusOK, users)
	return nil
}

func (ur *userRouter) createUser(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var user entity.User

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		httputils.WriteJSON(w, http.StatusBadRequest, err.Error())
		return nil
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(user)
	if err != nil {
		httputils.WriteJSON(w, http.StatusBadRequest, err.Error())
		return nil
	}

	err = ur.useCase.CreateUser(ctx, &user)
	if err != nil {
		return err
	}

	_ = httputils.WriteJSON(w, http.StatusOK, user)

	return nil
}

func (ur *userRouter) updateUser(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var user entity.User

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		httputils.WriteJSON(w, http.StatusBadRequest, err.Error())
		return nil
	}

	err = ur.useCase.UpdateUser(ctx, vars["id"], &user)
	if err != nil {
		return err
	}

	_ = httputils.WriteJSON(w, http.StatusOK, user)

	return nil
}

func (ur *userRouter) deleteUser(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	err := ur.useCase.DeleteUser(ctx, vars["id"])
	if err != nil {
		return err
	}

	_ = httputils.WriteJSON(w, http.StatusNoContent, nil)
	return nil
}
