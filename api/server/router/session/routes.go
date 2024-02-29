package session

import (
	"github.com/catwashere/microservice/api/server/router"
	"github.com/catwashere/microservice/internal/usecase/interfaces"
)

// sessionRouter is a router to talk with the session controller
type sessionRouter struct {
	routes  []router.Route
	useCase interfaces.IUseCaseSession
}

// NewRouter initializes a new image router
func NewRouter(uc interfaces.IUseCaseSession) router.Router {
	r := &sessionRouter{
		useCase: uc,
	}
	r.initRoutes()
	return r
}

// Routes returns the available routes to the image controller
func (r *sessionRouter) Routes() []router.Route {
	return r.routes
}

// initRoutes initializes the routes in the image router
func (r *sessionRouter) initRoutes() {
	r.routes = []router.Route{
		// GET
		router.NewPostRoute("/sessions", r.createSession),
		router.NewGetRoute("/sessions", r.getSessions),
		router.NewGetRoute("/sessions/:id", r.getSession),
		router.NewPutRoute("/sessions/:id", r.updateSession),
		router.NewDeleteRoute("/sessions/:id", r.deleteSession),
	}
}
