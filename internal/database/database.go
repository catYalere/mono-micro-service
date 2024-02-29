package database // import "github.com/catwashere/microservice/internal/database"

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/catwashere/microservice/internal/database/utils"
	"github.com/creasty/defaults"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IDatabase interface {
	Initialize(ctx context.Context) error
	Destroy(ctx context.Context) error
	GetDatabase() *mongo.Database
}

// Config provides the configuration for the Database Connection
type Config struct {
	Protocol string `default:"mongodb"`

	Hosts []*Host `default:"[{\"Hostname\":\"localhost\"}]"`

	Base    string            `default:"test"`
	User    string            `default:""`
	Pass    string            `default:""`
	Options map[string]string `default:"{}"`
}

type Host struct {
	Hostname string `default:"localhost"`
	Port     string `default:"27017"`
}

// Database contains instance details for the database
type Database struct {
	cfg    *Config
	client *mongo.Client
	base   *mongo.Database
}

var _ IDatabase = (*Database)(nil)

// New returns a new instance of the database based on the specified configuration.
func New(cfg *Config) *Database {
	defaults.Set(cfg)
	return &Database{
		cfg: cfg,
	}
}

// Initialize database connection.
func (d *Database) Initialize(ctx context.Context) error {
	err := d.create(ctx)

	if err != nil {
		//change for common error not mongodb error
		fmt.Print(err)
		return err
	}

	return nil
}

// Destroy database connection.
func (d *Database) Destroy(ctx context.Context) error {
	err := d.delete(ctx)

	if err != nil {
		//change for common error not mongodb error
		fmt.Print(err)
		return err
	}

	return nil
}

// buildHosts convert arrays to string.
func buildHosts(hosts []*Host) string {
	keys := make([]string, 0, len(hosts))
	for k := range hosts {
		keys = append(keys, utils.GetHost(hosts[k].Hostname, hosts[k].Port))
	}
	return strings.Join(keys, ",")
}

// buildConnectionURI transform Config into connection uri.
func buildConnectionURI(cfg *Config) string {
	var uri string
	if cfg.User != "" && cfg.Pass != "" {
		uri = fmt.Sprintf("%s://%s:%s@%s/%s", cfg.Protocol, cfg.User, cfg.Pass, buildHosts(cfg.Hosts), cfg.Base)
	} else {
		uri = fmt.Sprintf("%s://%s/%s", cfg.Protocol, buildHosts(cfg.Hosts), cfg.Base)
	}

	if len(cfg.Options) > 0 {
		uri = uri + utils.ParamsToString(cfg.Options)
	}

	return uri
}

// create initializes the client that the database uses.
func (d *Database) create(ctx context.Context) error {
	clientOptions := options.Client().ApplyURI(buildConnectionURI(d.cfg))

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return err
	}

	ctxdb, _ := context.WithTimeout(ctx, 10*time.Second)
	err = client.Connect(ctxdb)
	if err != nil {
		return err
	}

	d.client = client
	return nil
}

// delete disconnect the client that the database uses.
func (d *Database) delete(ctx context.Context) error {
	ctxdb, _ := context.WithTimeout(ctx, 10*time.Second)
	err := d.client.Disconnect(ctxdb)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetDatabase() *mongo.Database {
	if d.base == nil {
		d.base = d.client.Database(d.cfg.Base)
	}

	return d.base
}
