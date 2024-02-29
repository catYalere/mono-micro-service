package envs // import "github.com/catwashere/microservice/cmd/monolith/envs"

// Envs represents the list of well known env vars used by the Microservice API
type Service struct {
	Hostname string `required:"true" split_words:"true" default:"0.0.0.0"`
	Port     string `required:"true" split_words:"true" default:"8080"`
}

type Session struct {
	BaseUrl string `required:"true" split_words:"true" default:"http://localhost:8081"`
}

type User struct {
	BaseUrl string `required:"true" split_words:"true" default:"http://localhost:8082"`
}

type DB struct {
	Hostname string `required:"true" split_words:"true" default:"localhost"`
	Port     string `required:"true" split_words:"true" default:"27017"`
	Base     string `required:"true" split_words:"true" default:"auth"`
}

type Envs struct {
	User    User    `required:"true" split_words:"true"`
	Session Session `required:"true" split_words:"true"`
	Service Service `required:"true" split_words:"true"`
	DB      DB      `required:"true" split_words:"true"`
}
