package server // import "github.com/catwashere/microservice/api/server"

import (
	"net/http"

	"github.com/catwashere/microservice/api/server/httpstatus"
	"github.com/catwashere/microservice/api/server/httputils"
	"github.com/catwashere/microservice/api/types"
)

// makeErrorHandler makes an HTTP handler that decodes a API error and
// returns it in the response.
func makeErrorHandler(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusCode := httpstatus.FromError(err)
		response := &types.ErrorResponse{
			Message: err.Error(),
		}
		_ = httputils.WriteJSON(w, statusCode, response)
	}
}
