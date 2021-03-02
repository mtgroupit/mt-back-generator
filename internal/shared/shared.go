package shared

import (
	"strings"

	"github.com/fatih/camelcase"
)

// NameSQL converts name to "snake_case" format
func NameSQL(name string) string {
	return strings.ToLower(strings.Join(camelcase.Split(name), "_"))
}
