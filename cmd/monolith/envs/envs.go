package envs // import "github.com/catwashere/microservice/cmd/monolith/envs"

type Service struct {
	Hostname string `required:"true" split_words:"true" default:"0.0.0.0"`
	Port     string `required:"true" split_words:"true" default:"8080"`
}

type DB struct {
	Hostname string `required:"true" split_words:"true" default:"localhost"`
	Port     string `required:"true" split_words:"true" default:"27017"`
	Base     string `required:"true" split_words:"true" default:"monolith"`
}

// Envs represents the list of well known env vars used by the Microservice API
type Envs struct {
	Service Service `required:"true" split_words:"true"`
	DB      DB      `required:"true" split_words:"true"`
}
