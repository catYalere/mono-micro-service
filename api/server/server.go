package server // import "github.com/catwashere/microservice/api/server"

import (
	"context"
	"fmt"
	"net/http"

	"github.com/catwashere/microservice/api/server/httpstatus"
	"github.com/catwashere/microservice/api/server/httputils"
	"github.com/catwashere/microservice/api/server/middleware"
	"github.com/catwashere/microservice/api/server/router"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

const versionMatcher = "/v:version"

// Config provides the configuration for the API server
type Config struct {
	Hostname    string
	Port        string
	CorsHeaders string
	Version     string
}

// Server contains instance details for the server
type Server struct {
	cfg         *Config
	server      *httprouter.Router
	routers     []router.Router
	middlewares []middleware.Middleware
}

// New returns a new instance of the server based on the specified configuration.
func New(cfg *Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

// create initializes the main router the server uses.
func (s *Server) create() *httprouter.Router {
	router := httprouter.New()

	logrus.Info("Registering routers")
	for _, apiRouter := range s.routers {
		for _, r := range apiRouter.Routes() {
			f := s.makeHTTPHandler(r.Handler())

			logrus.Debugf("Registering %s, %s", r.Method(), r.Path())

			router.Handle(r.Method(), versionMatcher+r.Path(), f)
			router.Handle(r.Method(), r.Path(), f)
		}
	}

	return router
}

// UseMiddleware appends a new middleware to the request chain.
// This needs to be called before the API routes are configured.
func (s *Server) UseMiddleware(m middleware.Middleware) {
	s.middlewares = append(s.middlewares, m)
}

// InitRouter initializes the list of routers for the server.
func (s *Server) InitRouter(routers ...router.Router) {
	s.routers = append(s.routers, routers...)
}

// serveAPI spawns goroutine of a initialized server.
func (s *Server) ServeAPI() error {
	srv := s.create()
	addr := fmt.Sprintf("%s:%s", s.cfg.Hostname, s.cfg.Port)
	err := http.ListenAndServe(addr, srv)
	if err != nil {
		//change for common error not http error
		fmt.Print(err)
		return err
	}

	return nil
}

func (s *Server) convertToMap(ps httprouter.Params) map[string]string {
	vars := make(map[string]string)
	for i := range ps {
		vars[ps[i].Key] = ps[i].Value
	}
	return vars
}

func (s *Server) makeHTTPHandler(handler httputils.APIFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Define the context that we'll pass around info
		//
		// The 'context' will be used for global data that should
		// apply to all requests. Data that is specific to the
		// immediate function being called should still be passed
		// as 'args' on the function call.

		// use intermediate variable to prevent "should not use basic type
		// string as key in context.WithValue" golint errors

		ctx := context.WithValue(r.Context(), "userkey", r.Header.Get("User-Agent"))
		r = r.WithContext(ctx)
		handlerFunc := s.handlerWithGlobalMiddlewares(handler)

		vars := s.convertToMap(ps)

		if err := handlerFunc(ctx, w, r, vars); err != nil {
			statusCode := httpstatus.FromError(err)
			if statusCode >= 500 {
				logrus.Errorf("Handler for %s %s returned error: %v", r.Method, r.URL.Path, err)
			}
			makeErrorHandler(err)(w, r)
		}
	}
}
