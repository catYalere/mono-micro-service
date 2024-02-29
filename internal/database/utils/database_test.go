package utils // import "github.com/catwashere/microservice/internal/database/utils"

import (
	"testing"
)

func TestParamsToString(t *testing.T) {
	expected := "?var=value"
	value := ParamsToString(map[string]string{"var": "value"})
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}

	expected = "?var1=value1&var2=value2"
	value = ParamsToString(map[string]string{"var1": "value1", "var2": "value2"})
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}

	expected = "hostname:80"
	value = GetHost("hostname", "80")
	if value != expected {
		t.Errorf(`expected "%s", got "%v"`, expected, value)
	}
}
