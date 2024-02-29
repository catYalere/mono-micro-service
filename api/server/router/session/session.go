package session // import "github.com/catwashere/microservice/api/server/router/session"

import (
	"context"
	"encoding/json"
	"github.com/catwashere/microservice/api/server/httputils"
	"github.com/catwashere/microservice/internal/entity"
	"net/http"
)

func (ur *sessionRouter) getSession(ctx context.Context, w http.ResponseWriter, _ *http.Request, vars map[string]string) error {
	session, err := ur.useCase.GetSession(ctx, vars["id"])
	if err != nil {
		return err
	}

	_ = httputils.WriteJSON(w, http.StatusOK, session)
	return nil
}

func (ur *sessionRouter) getSessions(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	values := r.URL.Query()
	params := make(map[string]interface{}, len(values))
	for k, v := range values {
		params[k] = v[0]
	}

	sessions, err := ur.useCase.GetSessions(ctx, params)
	if err != nil {
		return err
	}

	_ = httputils.WriteJSON(w, http.StatusOK, sessions)
	return nil
}

func (ur *sessionRouter) createSession(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var session entity.Session

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&session)
	if err != nil {
		httputils.WriteJSON(w, http.StatusBadRequest, err.Error())
		return nil
	}

	err = ur.useCase.CreateSession(ctx, &session)
	if err != nil {
		return err
	}

	_ = httputils.WriteJSON(w, http.StatusOK, session)

	return nil
}

func (ur *sessionRouter) updateSession(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var session entity.Session

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&session)
	if err != nil {
		httputils.WriteJSON(w, http.StatusBadRequest, err.Error())
		return nil
	}

	err = ur.useCase.UpdateSession(ctx, vars["id"], &session)
	if err != nil {
		return err
	}

	_ = httputils.WriteJSON(w, http.StatusOK, session)

	return nil
}

func (ur *sessionRouter) deleteSession(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	err := ur.useCase.DeleteSession(ctx, vars["id"])
	if err != nil {
		return err
	}

	_ = httputils.WriteJSON(w, http.StatusNoContent, nil)
	return nil
}
