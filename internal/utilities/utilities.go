package utilities

import (
	"strings"

	"github.com/fatih/camelcase"
	"github.com/jinzhu/inflection"
)

// LowerTitle cancels strings.Title
func LowerTitle(in string) string {
	switch len(in) {
	case 0:
		return ""
	case 1:
		return strings.ToLower(string(in))
	default:
		return strings.ToLower(string(in[0])) + string(in[1:])
	}
}

// NameSQL converts name to "snake_case" format
func NameSQL(name string) string {
	return strings.ToLower(strings.Join(camelcase.Split(name), "_"))
}

// Pluralize - convert name to plural form
func Pluralize(name string) string {
	return inflection.Plural(name)
}

// ContainsStr returns "true" if 'slice' have element which equal to 'str'
func ContainsStr(slice []string, str string) bool {
	for i := range slice {
		if slice[i] == str {
			return true
		}
	}
	return false
}
