package utils // import "github.com/catwashere/microservice/internal/database/utils"

import (
	"strings"
)

// ParamsToString convert map to query string.
func ParamsToString(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k+"="+m[k])
	}
	return "?" + strings.Join(keys, "&")
}

func GetHost(hostname string, port string) string {
	return hostname + ":" + port
}
