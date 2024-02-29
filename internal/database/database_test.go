package database // import "github.com/catwashere/microservice/internal/database"

import (
	"strconv"
	"testing"

	"github.com/creasty/defaults"
)

func TestDefaultVariables(t *testing.T) {
	cfg := &Config{}
	defaults.Set(cfg)

	expected := "mongodb"
	value := cfg.Protocol
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}

	expected = "localhost"
	value = cfg.Hosts[0].Hostname
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}

	expected = "27017"
	value = cfg.Hosts[0].Port
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}

	expected = "test"
	value = cfg.Base
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}

	expected = ""
	value = cfg.User
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}

	expected = ""
	value = cfg.Pass
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}

	expected = "0"
	value = strconv.Itoa(len(cfg.Options))
	if value != expected {
		t.Errorf(`expected "%s" params, got "%v" params`, expected, value)
	}
}

func TestBuildHosts(t *testing.T) {
	expected := "hostname:80"
	value := buildHosts([]*Host{{Hostname: "hostname", Port: "80"}})
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}

	expected = "hostname1:80,hostname2:81"
	value = buildHosts([]*Host{{Hostname: "hostname1", Port: "80"}, {Hostname: "hostname2", Port: "81"}})
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}
}

func TestBuildConnectionURI(t *testing.T) {
	cfg := &Config{
		Protocol: "mongodb",
		Base:     "local",
		Hosts:    []*Host{{Hostname: "localhost", Port: "27017"}},
	}

	expected := "mongodb://localhost:27017/local"
	value := buildConnectionURI(cfg)
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}

	cfg.Hosts = []*Host{{Hostname: "localhost", Port: "27017"}, {Hostname: "localhost", Port: "27018"}}

	expected = "mongodb://localhost:27017,localhost:27018/local"
	value = buildConnectionURI(cfg)
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}

	cfg.User = "user"
	cfg.Pass = "pass"

	expected = "mongodb://user:pass@localhost:27017,localhost:27018/local"
	value = buildConnectionURI(cfg)
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}

	cfg.Options = map[string]string{"directConnection": "true"}

	expected = "mongodb://user:pass@localhost:27017,localhost:27018/local?directConnection=true"
	value = buildConnectionURI(cfg)
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}

	cfg.Hosts = []*Host{{Hostname: "localhost", Port: "27017"}}
	cfg.User = ""
	cfg.Pass = ""

	expected = "mongodb://localhost:27017/local?directConnection=true"
	value = buildConnectionURI(cfg)
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}

	cfg.Base = ""

	expected = "mongodb://localhost:27017/?directConnection=true"
	value = buildConnectionURI(cfg)
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}
}
