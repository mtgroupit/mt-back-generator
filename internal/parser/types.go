package parser

import (
	"strings"

	"github.com/mtgroupit/mt-back-generator/internal/models"
	"github.com/mtgroupit/mt-back-generator/internal/utilities"
	"github.com/pkg/errors"
)

var (
	intNumbericTypes      = []string{"int", "int32", "int64"}
	fractionNumbericTypes = []string{"float", "decimal"}
	standardNumbericTypes = append(intNumbericTypes, fractionNumbericTypes...)
	standardTypes         = append([]string{"string", "bool"}, standardNumbericTypes...)

	timeFormats = []string{"date", "date-time"}
	formats     = map[string][]string{"string": append([]string{"email", "phone", "url"}, timeFormats...)}
)

const (
	structTypePrefix        = "model."
	customTypePrefix        = "custom."
	arrayTypePrefix         = "[]"
	arrayOfStructTypePrefix = arrayTypePrefix + structTypePrefix
	arrayOfCustomTypePrefix = arrayTypePrefix + customTypePrefix

	TypesPrefix = "types."
)

// IsStandardType checks standard type or not
func IsStandardType(t string) bool {
	for i := range standardTypes {
		if t == standardTypes[i] {
			return true
		}
	}
	return false
}

func isIntNumbericType(t string) bool {
	for i := range intNumbericTypes {
		if t == intNumbericTypes[i] {
			return true
		}
	}
	return false
}

func isFractionNumbericType(t string) bool {
	for i := range fractionNumbericTypes {
		if t == fractionNumbericTypes[i] {
			return true
		}
	}
	return false
}

// IsTypesAdditionalType return true if go type is from types package
func IsTypesAdditionalType(BusinessType string) bool {
	return strings.HasPrefix(BusinessType, TypesPrefix)
}

// IsTimeFormat return true if format is one of timeFormats value
func IsTimeFormat(format string) bool {
	for i := range timeFormats {
		if format == timeFormats[i] {
			return true
		}
	}
	return false
}

func convertTypeToBusinessType(columnType, format string) string {
	switch {
	case strings.HasPrefix(columnType, structTypePrefix):
		return strings.Title(columnType[len(structTypePrefix):])
	case strings.HasPrefix(columnType, arrayOfStructTypePrefix):
		return strings.Title(columnType[len(arrayOfStructTypePrefix):])
	case strings.HasPrefix(columnType, customTypePrefix):
		return strings.Title(columnType[len(customTypePrefix):])
	case strings.HasPrefix(columnType, arrayOfCustomTypePrefix):
		return strings.Title(columnType[len(arrayOfCustomTypePrefix):])
	case strings.HasPrefix(columnType, arrayTypePrefix) && !strings.HasPrefix(columnType, arrayOfStructTypePrefix) && !strings.HasPrefix(columnType, arrayOfCustomTypePrefix):
		return convertStandardTypeToBusinessType(columnType[len(arrayTypePrefix):], format)
	default:
		return convertStandardTypeToBusinessType(columnType, format)
	}
}

func convertStandardTypeToBusinessType(columnType, format string) string {
	switch {
	case columnType == "decimal":
		return TypesPrefix + "Decimal"
	// case columnType == "uuid":
	// 	return TypesPrefix + "UUID"
	case columnType == "float":
		return "float64"
	case columnType == "string" && format == "date":
		return "date"
	case columnType == "string" && format == "date-time":
		return "date-time"
	default:
		return columnType
	}
}

func parseFieldType(options models.Options, customTypes map[string]models.CustomType) (bool, bool, string, error) {
	fieldType := options.Type
	BusinessType := convertTypeToBusinessType(fieldType, options.Format)
	lowerTitleBusinessType := utilities.LowerTitle(BusinessType)
	switch {
	case strings.HasPrefix(fieldType, customTypePrefix):
		if _, ok := customTypes[lowerTitleBusinessType]; !ok {
			return false, false, "", errors.Errorf(`Field refers to "%s" custom type which is not described anywhere`, lowerTitleBusinessType)
		}
		return true, false, BusinessType, nil
	case strings.HasPrefix(fieldType, arrayOfCustomTypePrefix):
		if _, ok := customTypes[lowerTitleBusinessType]; !ok {
			return false, false, "", errors.Errorf(`Fields refers to "%s" custom type which is not described anywhere`, lowerTitleBusinessType)
		}
		return true, true, BusinessType, nil
	case strings.HasPrefix(fieldType, arrayTypePrefix) && !strings.HasPrefix(fieldType, arrayOfCustomTypePrefix):
		if IsStandardType(fieldType[len(arrayTypePrefix):]) {
			return false, true, BusinessType, nil
		}

		return false, false, "", errors.Errorf(`"%s" is not correct type. You can use only one of standarad types %s or custom types`, BusinessType, strings.Join(standardTypes, ", "))
	default:
		return false, false, BusinessType, nil
	}
}

func parseColumnType(options models.Options, cfg *models.Config) (bool, bool, bool, string, error) {
	columnType := options.Type
	BusinessType := convertTypeToBusinessType(columnType, options.Format)
	lowerTitleBusinessType := utilities.LowerTitle(BusinessType)
	switch {
	case strings.HasPrefix(columnType, structTypePrefix):
		if _, ok := cfg.Models[lowerTitleBusinessType]; !ok {
			return false, false, false, "", errors.Errorf(`Field refers to "%s" model which is not described anywhere`, lowerTitleBusinessType)
		}
		return true, false, false, BusinessType, nil
	case strings.HasPrefix(columnType, arrayOfStructTypePrefix):
		if _, ok := cfg.Models[lowerTitleBusinessType]; !ok {
			return false, false, false, "", errors.Errorf(`Fields refers to "%s" model which is not described anywhere`, lowerTitleBusinessType)
		}
		return true, false, true, BusinessType, nil
	case strings.HasPrefix(columnType, customTypePrefix):
		if _, ok := cfg.CustomTypes[lowerTitleBusinessType]; !ok {
			return false, false, false, "", errors.Errorf(`Field refers to "%s" custom type which is not described anywhere`, lowerTitleBusinessType)
		}
		return false, true, false, BusinessType, nil
	case strings.HasPrefix(columnType, arrayOfCustomTypePrefix):
		if _, ok := cfg.CustomTypes[lowerTitleBusinessType]; !ok {
			return false, false, false, "", errors.Errorf(`Fields refers to "%s" custom type which is not described anywhere`, lowerTitleBusinessType)
		}
		return false, true, true, BusinessType, nil
	case strings.HasPrefix(columnType, arrayTypePrefix) && !strings.HasPrefix(columnType, arrayOfStructTypePrefix) && !strings.HasPrefix(columnType, arrayOfCustomTypePrefix):
		if IsStandardType(columnType[len(arrayTypePrefix):]) {
			return false, false, true, BusinessType, nil
		}
		return false, false, false, "", errors.Errorf(`"%s" is not correct type. You can use only one of standarad types %s or refers to any other model`, BusinessType, strings.Join(standardTypes, ", "))
	default:
		return false, false, false, BusinessType, nil
	}
}
