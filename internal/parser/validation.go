package parser

import (
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/mtgroupit/mt-back-generator/internal/models"
	"github.com/mtgroupit/mt-back-generator/internal/utilities"
	"github.com/pkg/errors"
)

func validate(cfg *models.Config) error {
	if cfg.Name == "" {
		return errors.New("name is empty")
	}
	if cfg.Module == "" {
		return errors.New("module is empty")
	}
	if cfg.AuthSrv == "" {
		return errors.New("auth-srv is empty")
	}
	if err := validateAccessAttributes(cfg.AccessAttributes); err != nil {
		return err
	}
	if err := validateRules(cfg); err != nil {
		return err
	}
	if err := validateModels(cfg); err != nil {
		return err
	}
	if err := validateCustomTypes(cfg.CustomTypes); err != nil {
		return err
	}
	return nil
}

func validateAccessAttributes(attributes []string) error {
	for _, attr := range attributes {
		if !isCorrectName(attr) {
			return errors.Errorf(`"%s" is invalid name for access attribute. %s`, attr, correctNameDescription)
		}
	}
	return nil
}

func validateRules(cfg *models.Config) error {
	for name, rule := range cfg.Rules {
		if !isCorrectName(name) {
			return errors.Errorf(`"%s" is invalid name for rule. %s`, name, correctNameDescription)
		}
		if len(rule.Attributes) == 0 {
			return errors.Errorf(`Rule "%s" has no any access attributes`, name)
		}
		if len(rule.Roles) == 0 {
			return errors.Errorf(`Rule "%s" has no any roles`, name)
		}

		for _, attr := range rule.Attributes {
			if !utilities.ContainsStr(cfg.AccessAttributes, attr) {
				return errors.Errorf(`Rule "%s" has access attribute "%s" that not exist`, name, attr)
			}
		}

		for _, role := range rule.Roles {
			if !utilities.ContainsStr(roles, role) {
				return errors.Errorf(`Rule "%s" has role "%s" that not exist. Available roles: %s`, name, role, strings.Join(roles, ", "))
			}
		}
	}

	return nil
}

func validateModels(cfg *models.Config) error {
	for name, model := range cfg.Models {
		if model.BoundToIsolatedEntity && model.Shared {
			return errors.Errorf(`Model: "%s". Id from isolated entity available only for not shared models`, name)
		}
		if !isCorrectName(name) {
			return errors.Errorf(`"%s" is invalid name for model. %s`, name, correctNameDescription)
		}
		if len(model.Columns) == 0 {
			return errors.Errorf(`Model "%s" has no any columns`, name)
		}

		for _, method := range model.Methods {
			if strings.Contains(method, "{noSecure}") {
				method = strings.Replace(method, "{noSecure}", "", -1)
				if !model.Shared {
					return errors.Errorf(`Model: "%s". Methods without authorization are allowed only for shared models`, name)
				}
				if models.IsMyMethod(method) {
					return errors.Errorf(`Model: "%s". "%s" Methods "My" must be with authorization`, name, method)
				}
			}
			if !model.IsStandardMethod(method) {
				switch {
				case strings.HasPrefix(method, "list") && strings.Contains(method, "("):
					return errors.Errorf(`Model: "%s". "%s"  is invalid as a adjust list. A valid adjust list shouldn't contain spaces before brackets. Correct method pattern: "list(column1, column3*, model1*(column1, model1(column1, column2))), where * means the field can be sorted by"`, name, method)
				case strings.HasPrefix(method, "get") && strings.Contains(method, "("):
					return errors.Errorf(`Model: "%s". "%s"  is invalid as a adjust get. A valid adjust get shouldn't contain spaces before brackets. Correct method pattern: "get(column1, column2)"`, name, method)
				case strings.HasPrefix(method, "edit") && strings.Contains(method, "("):
					return errors.Errorf(`Model: "%s". "%s"  is invalid as a adjust edit. A valid adjust edit shouldn't contain spaces before brackets. Correct method pattern: "edit(column1, column2)"`, name, method)
				default:
					if !isCorrectName(method) {
						return errors.Errorf(`Model: "%s". "%s"  is invalid name for method. %s`, name, method, correctNameDescription)
					}
				}
			}
			if model.BoundToIsolatedEntity && !models.IsMyMethod(method) {
				return errors.Errorf(`Model: "%s". "%s"  is invalid method for model with id from isolated entity. For model with id from isolated entity available only methods with "My" postfix`, name, method)
			}
		}

		haveDefaultSort := false
		for column, options := range model.Columns {
			if !isCorrectName(column) {
				return errors.Errorf(`Model: "%s". "%s"  is invalid name for column. %s`, name, column, correctNameDescription)
			}

			if !IsStandardColumn(column) {
				if column == "id" {
					if !(options.Type == "uuid" || options.Type == "int64") {
						return errors.Errorf(`Model: "%s". "%s"  is invalid type for id. Valid types is 'int64' or 'uuid'`, name, options.Type)
					}
				} else {
					if strings.HasPrefix(options.Type, arrayTypePrefix) {
						options.Type = options.Type[len(arrayTypePrefix):]
					}
					BusinessType := convertTypeToBusinessType(options.Type, options.Format)
					lowerTitleBusinessType := utilities.LowerTitle(BusinessType)
					if strings.HasPrefix(options.Type, structTypePrefix) {
						if _, ok := cfg.Models[lowerTitleBusinessType]; !ok {
							return errors.Errorf(`Model: "%s". Column "%s" refers to "%s" model which is not described anywhere`, name, column, lowerTitleBusinessType)
						}
					} else if strings.HasPrefix(options.Type, customTypePrefix) {
						if _, ok := cfg.CustomTypes[lowerTitleBusinessType]; !ok {
							return errors.Errorf(`Model: "%s". Column "%s" refers to "%s" custom type which is not described anywhere`, name, column, lowerTitleBusinessType)
						}
					} else {
						if !IsStandardType(options.Type) {
							_, ok := cfg.CustomTypes[options.Type]
							if !ok {
								return errors.Errorf(`Model: "%s". Column: "%s". "%s" is not correct type. You can use only one of standarad types %s, custom types or types refers to any other model`, name, column, BusinessType, strings.Join(standardTypes, ", "))
							}
						}
					}
				}
			}

			if len(options.Enum) > 0 {
				if column == "id" {
					return errors.Errorf(`Model: "%s". Column: "%s". Enum available only for not id columns`, name, column)
				}
			}

			if options.IsStruct {
				modelNameForBind := utilities.LowerTitle(options.BusinessType)
				if modelForBind, ok := cfg.Models[modelNameForBind]; ok && !modelForBind.Shared && model.Shared {
					return errors.Errorf(`Model: "%s". Column: "%s". "%s" is invalid type for column. Shared models can not use non-shared models as column type`, name, column, options.Type)
				}
			}

			if options.SortDefault {
				if options.IsStruct {
					return errors.Errorf(`Model: "%s". Column: "%s". Structure can not be as default column for sorting`, name, column)
				}
				if options.IsArray {
					return errors.Errorf(`Model: "%s". Column: "%s". Array can not be as default column for sorting`, name, column)
				}
				if !options.SortOn {
					return errors.Errorf(`Model: "%s". Column "%s" can not be as default column for sorting because sorting is not enabled for this column`, name, column)
				}
				if haveDefaultSort {
					return errors.Errorf(`Model "%s" has multiple columns as default for sorting, model should has one column as default for sorting`, name)
				}
				if options.SortOrderDefault != "" {
					orderDefault := strings.ToTitle(options.SortOrderDefault)
					if !(orderDefault == "ASC" || orderDefault == "DESC") {
						return errors.Errorf(`Model: "%s". Column: "%s". "%s" can not be as default order for sorting. Order for sorting can be only "ASC" or "DESC"`, name, column, options.SortOrderDefault)
					}
				}
				haveDefaultSort = true
			} else {
				if options.SortOrderDefault != "" {
					return errors.Errorf(`Model: "%s". Column: "%s". Default order for sorting allow only for fields which set as default for sorting`, name, column)
				}
			}

			if options.StrictFilter && options.Type != "string" {
				return errors.Errorf(`Model: "%s". Column: "%s". "strict-sorting" option not available for non "string" columns`, name, column)
			}

			if err := validateOptions(options); err != nil {
				return errors.Wrapf(err, `Model: "%s". Column: "%s"`, name, column)
			}
		}

		if err := validateRulesSet(model.RulesSet, cfg.Rules, model.Methods); err != nil {
			return errors.Wrapf(err, `Model: "%s"`, name)
		}
	}
	return nil
}

func validateRulesSet(rulesSet map[string][]string, rules map[string]models.Rule, modelMethods []string) error {
	for rule, methods := range rulesSet {
		if _, ok := rules[rule]; !ok {
			return errors.Errorf(`Rule "%s" has not exist`, rule)
		}
		for _, method := range methods {
			if !utilities.ContainsStr(modelMethods, method) {
				return errors.Errorf(`Method "%s" from rule "%s" has not exist`, method, rule)
			}
		}
	}
	return nil
}

func validateCustomTypes(customTypes map[string]models.CustomType) error {
	for customTypeName, customType := range customTypes {
		if !isCorrectName(customTypeName) {
			return errors.Errorf(`"%s" is invalid name for custom type. %s`, customTypeName, correctNameDescription)
		}
		if len(customType.Fields) == 0 {
			return errors.Errorf(`Custom type "%s" has no any fields`, customTypeName)
		}
		for fieldName, options := range customType.Fields {
			if !isCorrectName(fieldName) {
				return errors.Errorf(`Custom type: "%s". "%s" is invalid name for field. %s`, customTypeName, fieldName, correctNameDescription)
			}
			fieldType := options.Type
			if strings.HasPrefix(fieldType, arrayTypePrefix) {
				fieldType = fieldType[len(arrayTypePrefix):]
			}
			if strings.HasPrefix(fieldType, customTypePrefix) {
				fieldType = fieldType[len(customTypePrefix):]
			}
			if !IsStandardType(fieldType) {
				_, ok := customTypes[fieldType]
				if !ok {
					return errors.Errorf(`Custom type: "%s". Field: "%s". Custom type fields must have standard or other custom types. "%s" is not valid type`, customTypeName, fieldName, options.Type)
				}
			}

			if err := validateOptions(options); err != nil {
				return errors.Wrapf(err, `Custom type: "%s". Field: "%s"`, customTypeName, fieldName)
			}
		}
	}
	return nil
}

func validateOptions(options models.Options) error {
	if err := validateFormats(options.Type, options.Format); err != nil {
		return err
	}
	if err := validateEnum(options.Enum, options.Type); err != nil {
		return err
	}
	if err := validateDefault(options); err != nil {
		return err
	}
	return nil
}

func validateFormats(typeName, format string) error {
	if format == "" {
		return nil
	}
	if strings.HasPrefix(typeName, arrayTypePrefix) {
		typeName = typeName[len(arrayTypePrefix):]
	}

	typeFormats, ok := formats[typeName]
	if !ok {
		return errors.Errorf(`Type "%s" do not support formats`, typeName)
	}
	validFormat := false
	for i := range typeFormats {
		if format == typeFormats[i] {
			validFormat = true
		}
	}
	if !validFormat {
		return errors.Errorf(`Type "%s" do not support format: "%s"`, typeName, format)
	}
	return nil
}

func validateEnum(enum []string, columnType string) error {
	if len(enum) == 0 {
		return nil
	}
	if IsStandardType(columnType) {
		if columnType == "string" {
			return nil
		}
		return validateNumberEnum(enum, columnType)
	}
	return errors.Errorf(`Enum available only for standard types: %s`, strings.Join(standardTypes, ", "))
}

func validateNumberEnum(enum []string, columnType string) error {
	switch {
	case isIntNumbericType(columnType):
		return intNumbericEnumValidate(enum)
	case isFractionNumbericType(columnType):
		return fractionNumbericEnumValidate(enum)
	default:
		return errors.Errorf(`Enum of numbers available only for standard numberic types: %s`, strings.Join(standardNumbericTypes, ", "))
	}
}

func intNumbericEnumValidate(enum []string) error {
	for _, e := range enum {
		_, err := strconv.Atoi(e)
		if err != nil {
			return errors.Wrapf(err, `Incorrect enum. Enum for types %s must be in this format: [1, 2, 3]`, strings.Join(intNumbericTypes, ", "))
		}
	}
	return nil
}

func fractionNumbericEnumValidate(enum []string) error {
	for _, e := range enum {
		_, err := strconv.ParseFloat(e, 64)
		if err != nil {
			return errors.Wrapf(err, `Incorrect enum. Enum for types %s must be in this format: [1.1, 2, 0.3, .44]`, strings.Join(fractionNumbericTypes, ", "))
		}
	}
	return nil
}

func validateDefault(options models.Options) error {
	if options.Default == "" {
		return nil
	}

	if len(options.Enum) > 0 {
		found := false
		for _, e := range options.Enum {
			if e == options.Default {
				found = true
				break
			}
		}
		if !found {
			return errors.Errorf(`Default ("%s") should be one from enum: %s`, options.Default, strings.Join(options.Enum, ", "))
		}
	}

	if options.Format != "" {
		valid := false
		switch options.Format {
		case "date-time":
			valid = strfmt.IsDateTime(options.Default)
		case "date":
			valid = strfmt.IsDate(options.Default)
		case "email":
			valid = strfmt.IsEmail(options.Default)
		default:
			valid = true
		}
		if !valid {
			return errors.Errorf(`Default ("%s") should match the %s format`, options.Default, options.Format)
		}
	}

	return nil
}