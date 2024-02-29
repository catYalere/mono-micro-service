package envs // import "github.com/catwashere/microservice/cmd/monolith/envs"

import (
	"testing"

	"github.com/creasty/defaults"
)

func TestDefaultVariables(t *testing.T) {
	envs := &Envs{}
	defaults.Set(envs)

	expected := "0.0.0.0"
	value := envs.Service.Hostname
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}

	expected = "8082"
	value = envs.Service.Port
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}

	expected = "localhost"
	value = envs.DB.Hostname
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}

	expected = "27017"
	value = envs.DB.Port
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}

	expected = "users"
	value = envs.DB.Base
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}
}
