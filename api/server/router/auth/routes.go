package auth

import (
	"github.com/catwashere/microservice/api/server/router"
	"github.com/catwashere/microservice/internal/usecase/interfaces"
)

// authRouter is a router to talk with the auth controller
type authRouter struct {
	routes  []router.Route
	useCase interfaces.IUseCaseAuth
}

// NewRouter initializes a new image router
func NewRouter(uc interfaces.IUseCaseAuth) router.Router {
	r := &authRouter{
		useCase: uc,
	}
	r.initRoutes()
	return r
}

// Routes returns the available routes to the image controller
func (r *authRouter) Routes() []router.Route {
	return r.routes
}

// initRoutes initializes the routes in the image router
func (r *authRouter) initRoutes() {
	r.routes = []router.Route{
		// POST
		router.NewPostRoute("/login", r.login),
	}
}
