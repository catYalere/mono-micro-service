package user

import (
	"github.com/catwashere/microservice/api/server/router"
	"github.com/catwashere/microservice/internal/usecase/interfaces"
)

// userRouter is a router to talk with the user controller
type userRouter struct {
	routes  []router.Route
	useCase interfaces.IUseCaseUser
}

// NewRouter initializes a new image router
func NewRouter(uc interfaces.IUseCaseUser) router.Router {
	r := &userRouter{
		useCase: uc,
	}
	r.initRoutes()
	return r
}

// Routes returns the available routes to the image controller
func (r *userRouter) Routes() []router.Route {
	return r.routes
}

// initRoutes initializes the routes in the image router
func (r *userRouter) initRoutes() {
	r.routes = []router.Route{
		// GET
		router.NewPostRoute("/users", r.createUser),
		router.NewGetRoute("/users", r.getUsers),
		router.NewGetRoute("/users/:id", r.getUser),
		router.NewPutRoute("/users/:id", r.updateUser),
		router.NewDeleteRoute("/users/:id", r.deleteUser),
	}
}
