package main

import (
	"context"
	sessionrouter "github.com/catwashere/microservice/api/server/router/session"
	userrouter "github.com/catwashere/microservice/api/server/router/user"
	"github.com/catwashere/microservice/internal/usecase"
	"github.com/catwashere/microservice/pkg/crypto"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	"github.com/catwashere/microservice/api/server"
	"github.com/catwashere/microservice/api/server/router"
	authrouter "github.com/catwashere/microservice/api/server/router/auth"
	"github.com/catwashere/microservice/cmd/three_microservices_crypto/auth/envs"
	"github.com/catwashere/microservice/internal/database"
	"github.com/catwashere/microservice/internal/repositories"
)

type routerOptions struct {
	api *server.Server
}

type databaseOptions struct {
	db           *database.Database
	repositories []*repositories.IRepository[any]
}

func init() {
	// logrus custom configurations
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "severity",
		},
	})

	logrus.SetReportCaller(false)

	logrus.Info("Logrus successfully configured")
}

func initRouter(_ context.Context, opts routerOptions, routers []router.Router) {
	opts.api.InitRouter(routers...)
}

func initDatabase(ctx context.Context, opts databaseOptions) {
	opts.db.Initialize(ctx)
}

func main() {
	logrus.Info("Initializing Auth microservice")

	ctx := context.Background()

	var envars envs.Envs
	err := envconfig.Process("", &envars)

	if err != nil {
		logrusFields := logrus.Fields{"error": err.Error()}
		logrusMessage := "Error with the environment configuration"
		logrus.WithFields(logrusFields).Fatal(logrusMessage)
	}

	dcfg := &database.Config{
		Hosts: []*database.Host{
			{
				Hostname: envars.DB.Hostname,
				Port:     envars.DB.Port,
			},
		},
		Base: envars.DB.Base,
	}

	databaseOptions := databaseOptions{}
	databaseOptions.db = database.New(dcfg)
	initDatabase(ctx, databaseOptions)
	defer databaseOptions.db.Destroy(ctx)

	f1, err := os.ReadFile("./id_rsa")
	if err != nil {
		logrusFields := logrus.Fields{"error": err.Error()}
		logrusMessage := "Error loading private key"
		logrus.WithFields(logrusFields).Fatal(logrusMessage)
	}
	privateKey, err := crypto.FileBytesToPrivateKey(f1)
	if err != nil {
		logrusFields := logrus.Fields{"error": err.Error()}
		logrusMessage := "Error loading private key"
		logrus.WithFields(logrusFields).Fatal(logrusMessage)
	}

	uccfg := &usecase.Config{
		DB: databaseOptions.db,
		User: usecase.EntityConfig{
			Type:       usecase.External,
			BaseUrl:    envars.User.BaseUrl,
			PrivateKey: privateKey,
		},
		Session: usecase.EntityConfig{
			Type:    usecase.External,
			BaseUrl: envars.Session.BaseUrl,
		},
		Auth: usecase.EntityConfig{
			Type: usecase.Internal,
		},
	}

	uc := usecase.New(ctx, uccfg)
	routers := []router.Router{
		userrouter.NewRouter(uc.User()),
		sessionrouter.NewRouter(uc.Session()),
		authrouter.NewRouter(uc.Auth()),
	}

	rcfg := &server.Config{
		Hostname: envars.Service.Hostname,
		Port:     envars.Service.Port,
	}

	routerOptions := routerOptions{}
	routerOptions.api = server.New(rcfg)
	initRouter(ctx, routerOptions, routers)

	go func() {
		err = routerOptions.api.ServeAPI()
		if err != nil {
			logrusFields := logrus.Fields{"error": err.Error()}
			logrusMessage := "Error with the environment configuration"
			logrus.WithFields(logrusFields).Fatal(logrusMessage)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
