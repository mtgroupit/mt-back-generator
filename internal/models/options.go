package models

import (
	"github.com/mtgroupit/mt-back-generator/internal/utilities"
)

// Options contain properties of column
type Options struct {
	TitleName        string
	Type             string
	BusinessType     string
	Format           string
	Enum             []string
	Unique           bool
	Length           int64
	Default          string
	Required         bool
	Pattern          string
	SortOn           bool     `yaml:"sort-on"`
	SortDefault      bool     `yaml:"sort-default"`
	SortOrderDefault string   `yaml:"sort-order-default"`
	SortBy           []string `yaml:"sort-by"`
	StrictFilter     bool     `yaml:"strict-filter"`

	NestedSorts []string
	IsStruct    bool
	IsCustom    bool
	IsArray     bool
	Pk          string
}

// AppType - returns type which using in generated service
func (o Options) AppType() string {
	switch o.BusinessType {
	case "date":
		return "*time.Time"
	case "date-time":
		return "*time.Time"
	default:
		return o.BusinessType
	}
}

// FilterType - returns type which using in generated service for filtering in dal
func (o Options) FilterType() string {
	switch o.Type {
	case "uuid":
		return "uuid"
	case "string":
		return o.BusinessType
	case "integer":
		return o.Format
	case "number":
		if o.Format == "float" {
			return "float"
		}
		return "decimal"
	case "boolean":
		return "bool"
	default:
		if o.IsStruct {
			return "uuid"
		}
		return o.BusinessType
	}
}

// SQLName - returns column name in data base
func (o Options) SQLName(column string) (sqlName string) {
	sqlName = utilities.NameSQL(column)
	if !o.IsStruct {
		if o.IsArray || o.IsCustom {
			sqlName += "_json"
		}
	} else {
		if !o.IsArray {
			sqlName += "_id"
		}
	}
	return
}
