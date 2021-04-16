package parser

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/camelcase"
	"github.com/go-openapi/strfmt"
	"github.com/jinzhu/inflection"
	"github.com/mtgroupit/mt-back-generator/models"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var (
	isCorrectName          = regexp.MustCompile(`^[a-z][A-Za-z0-9]+$`).MatchString
	correctNameDescription = "A valid name must contain only letters and numbers in camelCase"

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

// IsTimeFormat return true if format is one of timeFormats value
func IsTimeFormat(format string) bool {
	for i := range timeFormats {
		if format == timeFormats[i] {
			return true
		}
	}
	return false
}

// IsTypesGoType return true if go type is from types package
func IsTypesGoType(goType string) bool {
	return strings.HasPrefix(goType, TypesPrefix)
}

// ReadYAMLCfg create models.Config from configFile
func ReadYAMLCfg(file string) (*models.Config, error) {
	cfg := models.Config{}
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// HandleCfg handles models.Config for fill config fields for templates
func HandleCfg(inCfg *models.Config) (cfg *models.Config, err error) {
	if err = validate(inCfg); err != nil {
		return
	}

	cfg = inCfg

	cfg.Name = formatName(cfg.Name)

	cfg.Description = strconv.Quote(cfg.Description)

	err = setDeepNesting(cfg)
	if err != nil {
		return
	}

	for customTypeName, customType := range cfg.CustomTypes {
		for field, options := range customType.Fields {
			if options.IsCustom, options.IsArray, options.GoType, err = parseFieldType(options, cfg.CustomTypes); err != nil {
				return nil, errors.Wrapf(err, `Custom type: "%s". Field: "%s"`, customTypeName, field)
			}

			if IsTypesGoType(options.GoType) {
				cfg.HaveTypesInCustomTypes = true
			}
			if options.Default != "" {
				if IsTimeFormat(options.Format) || options.Format == "email" {
					cfg.HaveConvInCustomTypes = true
				} else {
					cfg.HaveSwagInCustomTypes = true
				}

			}

			if IsTimeFormat(options.Format) {
				cfg.HaveTimeInCustomTypes = true
			}

			if options.Format == "email" {
				cfg.HaveEmailInCustomTypes = true
			}

			if options.IsCustom {
				options.Type = LowerTitle(options.GoType)
			} else {
				if options.Type == "int" {
					options.Type = "int32"
					options.GoType = "int32"
				}
				if options.IsArray {
					switch options.GoType {
					case "int":
						options.Type = "int32"
						options.GoType = "int32"
					case "float64":
						options.Type = "float"
					case TypesPrefix + "Decimal":
						options.Type = "decimal"
					case "time.Time":
						options.Type = "string"
					default:
						options.Type = options.GoType
					}
				}

				switch options.Type {
				case "int32", "int64":
					options.Format = options.Type
					options.Type = "integer"
				case "bool":
					options.Type = "boolean"
				case "float":
					options.Format = "float"
				}

				if isFractionNumbericType(options.Type) {
					options.Type = "number"
				}
			}

			customType.Fields[field] = options
		}
	}

	for name, model := range cfg.Models {
		model.TitleName = strings.Title(name)

		for i := range model.Tags {
			model.Tags[i] = strings.Title(model.Tags[i])
		}
		model.Tags = append([]string{strings.Title(name)}, model.Tags...)

		var props []models.MethodProps
		for _, method := range model.Methods {
			if isCustomMethod(method) {
				cfg.HaveCustomMethod = true
				model.HaveCustomMethod = true
			}

			var prop models.MethodProps
			if method == "delete" || method == "deleteMy" {
				prop.HTTPMethod = "delete"
			} else if method == "edit" || method == "editMy" || isAdjustEdit(method) {
				prop.HTTPMethod = "put"
			} else {
				prop.HTTPMethod = "post"
			}
			props = append(props, prop)

			if method == "list" || isAdjustList(method) {
				cfg.HaveListMethod = true
				model.HaveListMethod = true
			}
		}
		model.MethodsProps = props

		psql := []models.PsqlParams{}
		var indexLastNotArrOfStruct int
		for column, options := range model.Columns {
			switch column {
			case "createdAt":
				options.Type = "string"
				options.Format = "date-time"
				model.HaveCreatedAt = true
			case "createdBy":
				options.Type = "string"
				model.HaveCreatedBy = true
			case "modifiedAt":
				options.Type = "string"
				options.Format = "date-time"
				model.HaveModifiedAt = true
			case "modifiedBy":
				options.Type = "string"
				model.HaveModifiedBy = true
			}
			if options.IsStruct, options.IsCustom, options.IsArray, options.GoType, err = parseColumnType(options, cfg); err != nil {
				return nil, errors.Wrapf(err, `Model: "%s". Column: "%s"`, name, column)
			}

			if IsTypesGoType(options.GoType) {
				cfg.HaveTypes = true
				model.NeedTypes = true
			}

			if options.SortDefault {
				if options.SortOrderDefault != "" {
					orderDefault := strings.ToTitle(options.SortOrderDefault)
					if orderDefault == "ASC" || orderDefault == "DESC" {
						options.SortOrderDefault = orderDefault
					}
				}
			}

			if IsTimeFormat(options.Format) {
				cfg.HaveTime = true
				model.NeedTime = true
			}
			if options.Format == "email" {
				cfg.HaveEmail = true
				model.HaveEmail = true
			}

			if options.Default != "" {
				if IsTimeFormat(options.Format) || options.Format == "email" {
					cfg.HaveConv = true
					model.NeedConv = true
				} else {
					cfg.HaveSwag = true
				}
			}

			if options.IsStruct {
				model.HaveLazyLoading = true

				modelNameForBind := LowerTitle(options.GoType)

				err = cfg.AddBind(modelNameForBind, models.Bind{
					ModelName: name,
					FieldName: column,
					IsArray:   options.IsArray,
				})
				if err != nil {
					return nil, err
				}

				if options.IsArray {
					et := models.ExtraTable{}

					et.Name = NameSQL(name) + "_" + NameSQL(column)

					et.RefTableOne = strings.Title(name)
					et.RefIDOne = "id"
					et.FieldIDOne = NameSQL(name) + "_id"
					if cfg.Models[name].Columns["id"].Type == "uuid" {
						et.TypeIDOne = "uuid"
					} else {
						et.TypeIDOne = "integer"
					}

					et.RefTableTwo = options.GoType
					et.RefIDTwo = "id"
					et.FieldIDTwo = NameSQL(column) + "_id"
					if cfg.Models[LowerTitle(options.GoType)].Columns["id"].Type == "uuid" {
						et.TypeIDTwo = "uuid"
					} else {
						et.TypeIDTwo = "integer"
					}

					cfg.ExtraTables = append(cfg.ExtraTables, et)
				}
			} else {
				if options.IsArray || options.IsCustom {
					model.HaveJSON = true
				}
			}

			pp := models.PsqlParams{}
			pp.IsArray = options.IsArray
			pp.IsCustom = options.IsCustom
			pp.IsStruct = options.IsStruct
			if column == "id" {
				switch options.Type {
				case "uuid":
					options.GoType = "string"
					model.IDIsUUID = true

					pp.Type = "string"
					pp.TypeSQL = "uuid"
				case "int64":
					options.Type = "integer"

					pp.Type = "int64"
					pp.TypeSQL = "SERIAL"
				}
				options.TitleName = "ID"

				pp.Name = "ID"
				pp.SQLName = "id"
			} else {
				if column == "url" {
					options.TitleName = "URL"
				} else {
					options.TitleName = strings.Title(column)
				}
				if options.IsStruct || options.IsCustom {
					options.Type = LowerTitle(options.GoType)
				} else {
					if options.Type == "int" {
						options.Type = "int32"
						options.GoType = "int32"
					}
					if options.IsArray {
						switch options.GoType {
						case "int":
							options.Type = "int32"
							options.GoType = "int32"
						case "float64":
							options.Type = "float"
							cfg.HaveFloatArr = true
						case TypesPrefix + "Decimal":
							options.Type = "decimal"
						case "time.Time":
							options.Type = "string"
						default:
							options.Type = options.GoType
						}
					}

					switch options.Type {
					case "int32", "int64":
						options.Format = options.Type
						options.Type = "integer"
					case "bool":
						options.Type = "boolean"
					case "float":
						options.Format = "float"
					}

					if isFractionNumbericType(options.Type) {
						options.Type = "number"
					}
				}

				pp.Type = options.GoType
				pp.Name = options.TitleName
				if pp.IsStruct {
					pp.SQLName = NameSQL(column) + "_id"
					pp.FK = "id"
					if cfg.Models[LowerTitle(options.GoType)].Columns["id"].Type == "uuid" {
						pp.TypeSQL = "uuid"
					} else {
						pp.TypeSQL = "integer"
					}
				} else if pp.IsCustom || pp.IsArray {
					pp.SQLName = NameSQL(column) + "_json"
					pp.TypeSQL = "jsonb"
				} else {
					pp.SQLName = NameSQL(column)
					switch options.Type {
					case "string":
						if IsTimeFormat(options.Format) {
							pp.TypeSQL = "timestamp"
						} else {
							pp.TypeSQL = "text"
						}
					case "number":
						pp.TypeSQL = "numeric"
					default:
						pp.TypeSQL = options.Type
					}
				}

			}

			model.Columns[column] = options

			pp.Unique = options.Unique

			psql = append(psql, pp)
			if !pp.IsStruct || (!pp.IsArray && pp.IsStruct) {
				indexLastNotArrOfStruct = len(psql) - 1
			}
		}
		psql[indexLastNotArrOfStruct].Last = true
		model.Psql = psql

		var SQLSelect, sqlWhereParams, sqlAdd, sqlEdit []string
		for column, options := range model.Columns {
			if !options.IsStruct {
				sqlName := NameSQL(options.TitleName)
				titleName := options.TitleName
				if options.IsArray || options.IsCustom {
					sqlName += "_json"
					titleName += "JSON"
				} else {
					titleName = "m." + titleName
				}
				SQLSelect = append(SQLSelect, sqlName)
				if !options.IsArray && !options.IsCustom {
					if options.Type != "string" || IsTimeFormat(options.Format) || options.StrictFilter {
						sqlWhereParams = append(sqlWhereParams, fmt.Sprintf(`(CAST(:%s as text) IS NULL OR %s=:%s) AND (CAST(:%s as text) IS NULL OR %s<>:%s)`, column, sqlName, column, "not_"+column, sqlName, "not_"+column))
					} else {
						sqlWhereParams = append(sqlWhereParams, fmt.Sprintf(`(CAST(:%s as text) IS NULL OR LOWER(%s) LIKE LOWER(:%s)) AND (CAST(:%s as text) IS NULL OR LOWER(%s) NOT LIKE LOWER(:%s))`, column, sqlName, column, "not_"+column, sqlName, "not_"+column))
					}
				}
				if options.TitleName != "ID" {
					sqlAdd = append(sqlAdd, sqlName)
					if !isCreatedStandardColumn(column) {
						sqlEdit = append(sqlEdit, fmt.Sprintf("%s=:%s", sqlName, sqlName))
					}
				}
			} else {
				if !options.IsArray {
					sqlName := NameSQL(options.TitleName) + "_id"
					SQLSelect = append(SQLSelect, sqlName)
					sqlWhereParams = append(sqlWhereParams, fmt.Sprintf(`(CAST(:%s as text) IS NULL OR %s=:%s) AND (CAST(:%s as text) IS NULL OR %s<>:%s)`, column, sqlName, column, "not_"+column, sqlName, "not_"+column))
					sqlAdd = append(sqlAdd, sqlName)
					sqlEdit = append(sqlEdit, fmt.Sprintf("%s=:%s", sqlName, sqlName))
				}
			}
		}
		model.SQLSelectStr = strings.Join(SQLSelect, ",\n\t\t")
		model.SQLWhereParams = strings.Join(sqlWhereParams, " AND\n\t\t")
		if model.IDIsUUID {
			sqlAdd = append(sqlAdd, "id")
		}
		if !model.HaveCreatedBy {
			sqlAdd = append(sqlAdd, "created_by")
		}
		if !model.Shared {
			sqlAdd = append(sqlAdd, "isolated_entity_id")
		}
		model.SQLAddStr = fmt.Sprintf("(\n\t\t%s\n\t) VALUES (\n\t\t:%s\n\t)", strings.Join(sqlAdd, ",\n\t\t"), strings.Join(sqlAdd, ",\n\t\t:"))
		model.SQLEditStr = strings.Join(sqlEdit, ",\n\t\t")

		cfg.Models[name] = model
	}
	for name, model := range cfg.Models {

		if err = handleSorts(cfg.Models, &model, name); err != nil {
			return
		}

		if err = handleAdjustGets(cfg.Models, &model, name); err != nil {
			return
		}

		if err = handleAdjustLists(cfg.Models, &model, name); err != nil {
			return
		}

		if err = handleAdjustEdits(cfg.Models, &model, name); err != nil {
			return
		}

		cfg.Models[name] = model
	}

	for funcName, function := range cfg.Functions {
		newFunc := models.Function{}

		var inStrs, inStrsType, inStrsFull []string
		ins := make(map[string]string)
		for nameIn, typeIn := range function.In {
			nameIn = strings.Title(nameIn)
			ins[nameIn] = typeIn
			inStrs = append(inStrs, nameIn)
			inStrsType = append(inStrsType, typeIn)
			inStrsFull = append(inStrsFull, fmt.Sprintf("%s %s", nameIn, typeIn))
		}
		newFunc.In = ins
		newFunc.InStr = strings.Join(inStrs, ", ")
		newFunc.InStrType = strings.Join(inStrsType, ", ")
		newFunc.InStrFull = strings.Join(inStrsFull, ", ")

		newFunc.InStrParams = "params.Body." + strings.Join(inStrs, ", params.Body.")

		var outStrs, outStrsType, outStrsFull []string
		outs := make(map[string]string)
		for nameOut, typeOut := range function.Out {
			nameOut = strings.Title(nameOut)
			outs[nameOut] = typeOut
			outStrs = append(outStrs, nameOut)
			outStrsType = append(outStrsType, typeOut)
			outStrsFull = append(outStrsFull, fmt.Sprintf("%s %s", nameOut, typeOut))
		}
		newFunc.Out = outs
		newFunc.OutStr = strings.Join(outStrs, ", ")
		newFunc.OutStrType = strings.Join(outStrsType, ", ")
		newFunc.OutStrFull = strings.Join(outStrsFull, ", ")

		newFunc.HaveOut = len(newFunc.OutStr) != 0

		cfg.Functions[funcName] = newFunc
	}

	titleize(cfg)

	return
}

// NameSQL converts name to "snake_case" format
func NameSQL(name string) string {
	return strings.ToLower(strings.Join(camelcase.Split(name), "_"))
}

// Pluralize - convert name to plural form
func Pluralize(name string) string {
	return inflection.Plural(name)
}

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

	if err := validateModels(cfg); err != nil {
		return err
	}
	if err := validateCustomTypes(cfg.CustomTypes); err != nil {
		return err
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
			if isCustomMethod(method) {
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
			if model.BoundToIsolatedEntity && !IsMyMethod(method) {
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
					goType := convertTypeToGoType(options.Type, options.Format)
					lowerTitleGoType := LowerTitle(goType)
					if strings.HasPrefix(options.Type, structTypePrefix) {
						if _, ok := cfg.Models[lowerTitleGoType]; !ok {
							return errors.Errorf(`Model: "%s". Column "%s" refers to "%s" model which is not described anywhere`, name, column, lowerTitleGoType)
						}
					} else if strings.HasPrefix(options.Type, customTypePrefix) {
						if _, ok := cfg.CustomTypes[lowerTitleGoType]; !ok {
							return errors.Errorf(`Model: "%s". Column "%s" refers to "%s" custom type which is not described anywhere`, name, column, lowerTitleGoType)
						}
					} else {
						if !IsStandardType(options.Type) {
							_, ok := cfg.CustomTypes[options.Type]
							if !ok {
								return errors.Errorf(`Model: "%s". Column: "%s". "%s" is not correct type. You can use only one of standarad types %s, custom types or types refers to any other model`, name, column, goType, strings.Join(standardTypes, ", "))
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
				modelNameForBind := LowerTitle(options.GoType)
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

func formatName(name string) string {
	splitedName := regexp.MustCompile("[^a-zA-Z0-9]+").Split(name, -1)
	for i := range splitedName {
		splitedName[i] = strings.ToLower(splitedName[i])
	}
	return strings.Join(splitedName, "-")
}

func setDeepNesting(cfg *models.Config) (err error) {
	for name, model := range cfg.Models {
		model.DeepNesting, err = countDeepNesting(name, cfg)
		if err != nil {
			return
		}
		if model.DeepNesting > cfg.MaxDeepNesting {
			cfg.MaxDeepNesting = model.DeepNesting
		}

		cfg.Models[name] = model
	}
	return
}

func countDeepNesting(model string, cfg *models.Config) (int, error) {
	var err error
	deepNesting := 0
	for _, options := range cfg.Models[model].Columns {
		if options.IsStruct, _, options.IsArray, options.GoType, err = parseColumnType(options, cfg); err != nil {
			return 0, err
		}
		if options.IsStruct {
			columnDeepNesting := 1
			var modelName string
			if options.IsArray {
				modelName = options.Type[8:]
			} else {
				modelName = options.Type[6:]
			}
			optTypeDeepNesting, err := countDeepNesting(modelName, cfg)
			if err != nil {
				return 0, err
			}
			if optTypeDeepNesting > 0 {
				columnDeepNesting += optTypeDeepNesting
			}
			if columnDeepNesting > deepNesting {
				deepNesting = columnDeepNesting
			}
		}

	}
	return deepNesting, nil
}

func parseFieldType(options models.Options, customTypes map[string]models.CustomType) (bool, bool, string, error) {
	fieldType := options.Type
	goType := convertTypeToGoType(fieldType, options.Format)
	lowerTitleGoType := LowerTitle(goType)
	switch {
	case strings.HasPrefix(fieldType, customTypePrefix):
		if _, ok := customTypes[lowerTitleGoType]; !ok {
			return false, false, "", errors.Errorf(`Field refers to "%s" custom type which is not described anywhere`, lowerTitleGoType)
		}
		return true, false, goType, nil
	case strings.HasPrefix(fieldType, arrayOfCustomTypePrefix):
		if _, ok := customTypes[lowerTitleGoType]; !ok {
			return false, false, "", errors.Errorf(`Fields refers to "%s" custom type which is not described anywhere`, lowerTitleGoType)
		}
		return true, true, goType, nil
	case strings.HasPrefix(fieldType, arrayTypePrefix) && !strings.HasPrefix(fieldType, arrayOfCustomTypePrefix):
		if IsStandardType(fieldType[len(arrayTypePrefix):]) {
			return false, true, goType, nil
		}

		return false, false, "", errors.Errorf(`"%s" is not correct type. You can use only one of standarad types %s or custom types`, goType, strings.Join(standardTypes, ", "))
	default:
		return false, false, goType, nil
	}
}

func parseColumnType(options models.Options, cfg *models.Config) (bool, bool, bool, string, error) {
	columnType := options.Type
	goType := convertTypeToGoType(columnType, options.Format)
	lowerTitleGoType := LowerTitle(goType)
	switch {
	case strings.HasPrefix(columnType, structTypePrefix):
		if _, ok := cfg.Models[lowerTitleGoType]; !ok {
			return false, false, false, "", errors.Errorf(`Field refers to "%s" model which is not described anywhere`, lowerTitleGoType)
		}
		return true, false, false, goType, nil
	case strings.HasPrefix(columnType, arrayOfStructTypePrefix):
		if _, ok := cfg.Models[lowerTitleGoType]; !ok {
			return false, false, false, "", errors.Errorf(`Fields refers to "%s" model which is not described anywhere`, lowerTitleGoType)
		}
		return true, false, true, goType, nil
	case strings.HasPrefix(columnType, customTypePrefix):
		if _, ok := cfg.CustomTypes[lowerTitleGoType]; !ok {
			return false, false, false, "", errors.Errorf(`Field refers to "%s" custom type which is not described anywhere`, lowerTitleGoType)
		}
		return false, true, false, goType, nil
	case strings.HasPrefix(columnType, arrayOfCustomTypePrefix):
		if _, ok := cfg.CustomTypes[lowerTitleGoType]; !ok {
			return false, false, false, "", errors.Errorf(`Fields refers to "%s" custom type which is not described anywhere`, lowerTitleGoType)
		}
		return false, true, true, goType, nil
	case strings.HasPrefix(columnType, arrayTypePrefix) && !strings.HasPrefix(columnType, arrayOfStructTypePrefix) && !strings.HasPrefix(columnType, arrayOfCustomTypePrefix):
		if IsStandardType(columnType[len(arrayTypePrefix):]) {
			return false, false, true, goType, nil
		}
		return false, false, false, "", errors.Errorf(`"%s" is not correct type. You can use only one of standarad types %s or refers to any other model`, goType, strings.Join(standardTypes, ", "))
	default:
		return false, false, false, goType, nil
	}
}

func convertTypeToGoType(columnType, format string) string {
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
		return convertStandardTypeToGoType(columnType[len(arrayTypePrefix):], format)
	default:
		return convertStandardTypeToGoType(columnType, format)
	}
}

func convertStandardTypeToGoType(columnType, format string) string {
	switch {
	case columnType == "decimal":
		return TypesPrefix + "Decimal"
	// case columnType == "uuid":
	// 	return TypesPrefix + "UUID"
	case columnType == "float":
		return "float64"
	case columnType == "string" && IsTimeFormat(format):
		return "time.Time"
	default:
		return columnType
	}
}

// isCustomMethod return true if method is custom
func isCustomMethod(method string) bool {
	method = strings.ToLower(method)
	fmt.Println(method, isAdjustGet(method))
	if method == "get" || method == "add" || method == "delete" || method == "edit" || method == "list" || isAdjustList(method) || isAdjustGet(method) || isAdjustEdit(method) || IsMyMethod(method) {
		return false
	}
	return true
}

// IsMyMethod return true if method is standard my method
func IsMyMethod(method string) bool {
	method = strings.ToLower(method)
	if method == "getmy" || method == "addmy" || method == "deletemy" || method == "editmy" || method == "editoraddmy" || regexp.MustCompile(`^getmy.+`).Match([]byte(method)) || regexp.MustCompile(`^editmy.+`).Match([]byte(method)) {
		return true
	}
	return false
}

func isAdjustGet(method string) bool {
	return regexp.MustCompile(`^get(My|my)?\(.+\)(\[[a-zA-Z0-9]+\])?$`).Match([]byte(method))
}

func isAdjustEdit(method string) bool {
	return regexp.MustCompile(`^edit(My|my)?\(.+\)(\[[a-zA-Z0-9]+\])?$`).Match([]byte(method))
}

func isAdjustList(method string) bool {
	return regexp.MustCompile(`^list\(.+\)(\[[a-zA-Z0-9]+\])?$`).Match([]byte(method))
}

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

func titleize(cfg *models.Config) {
	titleModels := make(map[string]models.Model)
	for modelName, model := range cfg.Models {
		for i := range cfg.Models[modelName].Methods {
			model.Methods[i] = strings.Title(model.Methods[i])
		}
		titleModels[strings.Title(modelName)] = model
	}
	cfg.Models = titleModels
}

func expandStrNestedFields(method string) string {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9]+\*{0,1}\((?P<value>.+)\)(\[[a-zA-Z0-9]+\])?$`)
	result := []byte{}
	template := "$value"
	result = pattern.ExpandString(result, template, method, pattern.FindSubmatchIndex([]byte(method)))

	return string(result)
}

func expandName(method string) string {
	return strings.TrimSuffix(regexp.MustCompile("[^a-zA-Z0-9*]").Split(method, 2)[0], "*")
}

func expandNamePostfixForAdjustMethods(method string) string {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9]+\*{0,1}\(.+\)(\[(?P<value>[a-zA-Z0-9]+)\])?$`)
	result := []byte{}
	template := "$value"
	result = pattern.ExpandString(result, template, method, pattern.FindSubmatchIndex([]byte(method)))

	return string(result)
}

func getNameForAdjustMethods(method string) (result string) {
	methodName := expandName(method)
	methodNamePostfix := expandNamePostfixForAdjustMethods(method)

	if methodNamePostfix == "" {
		fieldsStr := expandStrNestedFields(method)
		fieldsFull := splitFields(fieldsStr)
		fields := trimFieldsSuffix(fieldsFull)
		for i := range fields {
			fields[i] = strings.TrimSuffix(fields[i], "*")
			if strings.ToLower(fields[i]) == "id" {
				fields[i] = strings.ToUpper(fields[i])
			} else {
				fields[i] = strings.Title(fields[i])
			}
		}
		result = methodName + strings.Join(fields, "")
	} else {
		result = methodName + strings.Title(methodNamePostfix)
	}

	return
}

func splitFields(fields string) []string {
	var result []string
	for {
		fields = strings.Trim(fields, ", ")
		if strings.Index(fields, ",") >= 0 {
			if strings.Index(fields, ",") < strings.Index(fields, "(") || strings.Index(fields, "(") == -1 {
				substrs := regexp.MustCompile("[^a-zA-Z0-9*]+").Split(fields, 2)
				result = append(result, substrs[0])
				fields = substrs[1]
			} else {
				counter := 0
				var endBracket int
				for i, symb := range []rune(fields) {
					switch symb {
					case []rune("(")[0]:
						counter++
					case []rune(")")[0]:
						counter--
						if counter == 0 {
							endBracket = i
						}
					}
					if counter == 0 && i > strings.Index(fields, "(") {
						break
					}
				}
				result = append(result, fields[:endBracket+1])
				fields = fields[endBracket+1:]
			}
		} else {
			if fields != "" {
				result = append(result, fields)
			}
			break
		}
	}
	return result
}

func trimFieldsSuffix(fields []string) (out []string) {
	for i := range fields {
		out = append(out, regexp.MustCompile("[^a-zA-Z0-9*]").Split(fields[i], 2)[0])
	}
	return
}

func isStruct(method string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]+\*{0,1}\(.+\)$`).Match([]byte(method))
}

func handleNestedObjs(modelsIn map[string]models.Model, modelName, elem, nesting, parent string, isArray bool) ([]models.NestedObjProps, error) {
	objs := []models.NestedObjProps{}
	obj := models.NestedObjProps{}

	obj.Shared = modelsIn[modelName].Shared

	field := expandName(elem)
	fieldsStr := expandStrNestedFields(elem)
	fieldsFull := splitFields(fieldsStr)
	fields := trimFieldsSuffix(fieldsFull)
	SQLSelect := []string{}
	haveID := false
	for i := range fields {
		var haveFieldInColumns bool
		var structModel string
		for column, options := range modelsIn[modelName].Columns {
			if column == fields[i] {
				structModel = LowerTitle(options.GoType)
				obj.Type = strings.Title(modelName)
				haveFieldInColumns = true
				break
			}
		}
		if !haveFieldInColumns {
			return nil, errors.Errorf(`Model "%s" does not contain "%s" column for adjust list`, modelName, fields[i])
		}
		if strings.ToLower(fields[i]) == "id" {
			haveID = true
		}

		structIsArr := false
		for column, options := range modelsIn[modelName].Columns {
			if fields[i] == column {
				if !options.IsStruct {
					SQLSelect = append(SQLSelect, NameSQL(column))
				} else {
					if !options.IsArray {
						SQLSelect = append(SQLSelect, NameSQL(fields[i])+"_id")
					} else {
						structIsArr = true
					}
					obj.NeedLazyLoading = true
				}
			}
		}

		if isStruct(fieldsFull[i]) {
			objsForAdd, err := handleNestedObjs(modelsIn, structModel, fieldsFull[i], nesting+strings.Title(field), strings.Title(modelName), structIsArr)
			if err != nil {
				return nil, err
			}
			objs = append(objs, objsForAdd...)
		}

	}
	if !haveID {
		SQLSelect = append(SQLSelect, "id")
	}
	obj.SQLSelect = strings.Join(SQLSelect, ",\n\t\t")
	obj.Path = nesting
	obj.ParentStruct = parent
	obj.IsArray = isArray
	obj.Name = strings.Title(field)

	result := []models.NestedObjProps{}
	result = append(result, obj)
	if len(objs) > 0 {
		result = append(result, objs...)
	}
	return result, nil
}

func handleSorts(modelsMap map[string]models.Model, model *models.Model, modelName string) error {
	result := *model
	for column, options := range result.Columns {
		if options.SortOn {
			if column == "id" {
				return errors.Errorf(`Model: "%s". Sorting by id is not avaliable`, modelName)
			}
			if options.IsArray {
				return errors.Errorf(`Model: "%s". Column: "%s". Sorting by array is not avaliable`, modelName, column)
			}
			if !options.IsStruct && len(options.SortBy) != 0 {
				return errors.Errorf(`Model: "%s". Column: "%s". Field 'sort-by' is not avaliable for non-structure column`, modelName, column)
			}

			for i, sort := range options.SortBy {
				subModelName := options.GoType
				fields := strings.Split(sort, ".")
				options.NestedSorts = append(options.NestedSorts, "."+strings.Title(column))
				for j, field := range fields {
					if field == "id" {
						return errors.Errorf(`Model: "%s". Column: "%s". Sort-by: "%s". Sorting by id is not avaliable`, modelName, column, sort)
					}
					subModel := modelsMap[LowerTitle(subModelName)]
					var ok bool
					var typeSortByColumn string
					for subColumn, subOptions := range subModel.Columns {
						if subColumn == field {
							if j == len(fields)-1 {
								switch subOptions.GoType {
								case "int32", "int64", "string":
									typeSortByColumn = strings.Title(subOptions.GoType)
								default:
									return errors.Errorf(`Model: "%s". Column: "%s". Sort-by: "%s". Type "%s" is not avaliable for sorting`, modelName, column, sort, subOptions.GoType)
								}
							}
							ok = true
							subModelName = subOptions.GoType
						}
					}
					if !ok {
						return errors.Errorf(`Model: "%s". Column: "%s". Sort-by: "%s". Is not valid element for sort-by`, modelName, column, sort)
					}
					options.NestedSorts[i] += "." + strings.Title(field)
					if j == len(fields)-1 {
						options.NestedSorts[i] += "." + typeSortByColumn
					}
				}
			}
		} else {
			if len(options.SortBy) != 0 {
				return errors.Errorf(`Model: "%s". Field 'sort-by' is not avaliable, if sort-on is not true`, modelName)
			}
		}
		result.Columns[column] = options
	}
	model = &result
	return nil
}

func handleAdjustLists(modelsMap map[string]models.Model, model *models.Model, modelName string) error {
	result := *model
	for i, method := range result.Methods {
		if isAdjustList(method) {
			var SQLSelect, sqlWhereParams, filtredFields []string
			fieldsStr := expandStrNestedFields(method)
			fieldsFull := splitFields(fieldsStr)
			fields := trimFieldsSuffix(fieldsFull)
			haveID := false
			result.MethodsProps[i].JSONColumns = map[string]bool{}
			for j := range fields {
				var needFilter bool
				if strings.HasSuffix(fields[j], "*") {
					needFilter = true
					fields[j] = strings.TrimSuffix(fields[j], "*")
					filtredFields = append(filtredFields, fields[j])
				}
				var haveFieldInColumns bool
				var structModel string
				for column, options := range result.Columns {
					if column == fields[j] {
						structModel = options.Type
						haveFieldInColumns = true
					}
				}
				if !haveFieldInColumns {
					return errors.Errorf(`Model "%s" does not contain "%s" column for method "%s"`, modelName, fields[j], method)
				}

				if strings.ToLower(fields[j]) == "id" {
					haveID = true
					fields[j] = strings.ToUpper(fields[j])
				} else {
					fields[j] = strings.Title(fields[j])
				}

				structIsArr := false
				for column, options := range result.Columns {
					if fields[j] == options.TitleName {
						if !options.IsStruct {
							sqlName := NameSQL(options.TitleName)
							if options.IsArray || options.IsCustom {
								result.MethodsProps[i].HaveJSON = true
								sqlName += "_json"
							}
							SQLSelect = append(SQLSelect, sqlName)
							if needFilter && !options.IsArray {
								if options.Type != "string" || IsTimeFormat(options.Format) || options.StrictFilter {
									sqlWhereParams = append(sqlWhereParams, fmt.Sprintf(`(CAST(:%s as text) IS NULL OR %s=:%s) AND (CAST(:%s as text) IS NULL OR %s<>:%s)`, column, sqlName, column, "not_"+column, sqlName, "not_"+column))
								} else {
									sqlWhereParams = append(sqlWhereParams, fmt.Sprintf(`(CAST(:%s as text) IS NULL OR LOWER(%s) LIKE LOWER(:%s)) AND (CAST(:%s as text) IS NULL OR LOWER(%s) NOT LIKE LOWER(:%s))`, column, sqlName, column, "not_"+column, sqlName, "not_"+column))
								}
							}
						} else {
							if !options.IsArray {
								SQLSelect = append(SQLSelect, NameSQL(options.TitleName)+"_id")
								if needFilter {
									sqlWhereParams = append(sqlWhereParams, fmt.Sprintf(`(CAST(:%s as text) IS NULL OR %s=:%s) AND (CAST(:%s as text) IS NULL OR %s<>:%s)`, column, NameSQL(options.TitleName)+"_id", column, "not_"+column, NameSQL(options.TitleName)+"_id", "not_"+column))
								}
							} else {
								structIsArr = true
							}
						}
						result.MethodsProps[i].JSONColumns[column] = options.IsArray || options.IsCustom
					}
				}

				if isStruct(fieldsFull[j]) {
					result.MethodsProps[i].NeedLazyLoading = true

					objsForAdd, err := handleNestedObjs(modelsMap, structModel, fieldsFull[j], "", strings.Title(modelName), structIsArr)
					if err != nil {
						return err
					}
					result.MethodsProps[i].NestedObjs = append(result.MethodsProps[i].NestedObjs, objsForAdd...)
				}
			}

			result.Methods[i] = getNameForAdjustMethods(method)
			if !haveID {
				SQLSelect = append(SQLSelect, "id")
			}
			result.MethodsProps[i].AdjustSQLSelect = strings.Join(SQLSelect, ",\n\t\t")
			result.MethodsProps[i].AdjustListSQLWhereProps = strings.Join(sqlWhereParams, " AND\n\t\t")
			result.MethodsProps[i].FilteredFields = filtredFields
			result.MethodsProps[i].IsAdjustList = true

			sort.Slice(result.MethodsProps[i].NestedObjs, func(a, b int) bool {
				return result.MethodsProps[i].NestedObjs[a].Path < result.MethodsProps[i].NestedObjs[b].Path
			})

			for j := range result.MethodsProps[i].NestedObjs {
				if j == 0 {
					result.MethodsProps[i].NestedObjs[j].IsFirstForLazyLoading = true
					if len(result.MethodsProps[i].NestedObjs) == 1 {
						result.MethodsProps[i].NestedObjs[j].IsLastForLazyLoading = true
					}
				} else {
					if j == len(result.MethodsProps[i].NestedObjs)-1 {
						result.MethodsProps[i].NestedObjs[j].IsLastForLazyLoading = true
					}
					if result.MethodsProps[i].NestedObjs[j].Path != result.MethodsProps[i].NestedObjs[j-1].Path {
						result.MethodsProps[i].NestedObjs[j-1].IsLastForLazyLoading = true
						result.MethodsProps[i].NestedObjs[j].IsFirstForLazyLoading = true
					}
				}
			}
		}
	}
	model = &result
	return nil
}

func handleAdjustGets(modelsMap map[string]models.Model, model *models.Model, modelName string) error {
	result := *model
	for i, method := range result.Methods {
		if isAdjustGet(method) {
			var SQLSelect, adjustGetJSONColumns []string
			haveID := false

			fieldsStr := expandStrNestedFields(method)
			fields := splitFields(fieldsStr)
			for j := range fields {
				var haveFieldInColumns bool
				for column := range result.Columns {
					if column == fields[j] {
						haveFieldInColumns = true
					}
				}
				if !haveFieldInColumns {
					return errors.Errorf(`Model "%s" does not contain "%s" column for method "%s"`, modelName, fields[j], method)
				}

				if strings.ToLower(fields[j]) == "id" {
					haveID = true
				}

				adjustGetJSONColumns = append(adjustGetJSONColumns, fields[j])

				for column, options := range result.Columns {
					if fields[j] == column {
						if !options.IsStruct {
							sqlName := NameSQL(options.TitleName)
							if options.IsArray || options.IsCustom {
								sqlName += "_json"
							}
							SQLSelect = append(SQLSelect, sqlName)
						} else {
							if !options.IsArray {
								SQLSelect = append(SQLSelect, NameSQL(options.TitleName)+"_id")
							}
						}
					}
				}
			}
			result.Methods[i] = getNameForAdjustMethods(method)
			if !haveID {
				SQLSelect = append(SQLSelect, "id")
			}
			result.MethodsProps[i].AdjustSQLSelect = strings.Join(SQLSelect, ",\n\t\t")
			result.MethodsProps[i].AdjustGetJSONColumns = adjustGetJSONColumns
		}
	}
	model = &result
	return nil
}

func handleAdjustEdits(modelsMap map[string]models.Model, model *models.Model, modelName string) error {
	result := *model
	for i, method := range result.Methods {
		if isAdjustEdit(method) {
			var sqlEdit, editableFields []string
			count := 1
			if !model.Shared {
				count++
			}
			fieldsStr := expandStrNestedFields(method)
			fields := splitFields(fieldsStr)
			for j := range fields {
				if IsStandardColumn(fields[j]) {
					return errors.Errorf(`Model "%s". Method: "%s". "%s" can not be used in adjust edit method, it edits automatically`, modelName, method, fields[j])
				}

				var haveFieldInColumns bool
				for column := range result.Columns {
					if column == fields[j] {
						haveFieldInColumns = true
					}
				}
				if !haveFieldInColumns {
					return errors.Errorf(`Model "%s" does not contain "%s" column for method "%s"`, modelName, fields[j], method)
				}

				editableFields = append(editableFields, fields[j])

				fields[j] = strings.Title(fields[j])

				for _, options := range result.Columns {
					if fields[j] == options.TitleName {
						sqlName := NameSQL(options.TitleName)
						if !options.IsStruct {
							if options.IsArray || options.IsCustom {
								sqlName += "_json"
							}
							sqlEdit = append(sqlEdit, fmt.Sprintf("%s=:%s", sqlName, sqlName))
						} else {
							if !options.IsArray {
								sqlName+="_id"
								sqlEdit = append(sqlEdit, fmt.Sprintf("%s=:%s", sqlName, sqlName))
							}
						}
					}
				}
			}
			result.Methods[i] = getNameForAdjustMethods(method)
			result.MethodsProps[i].CustomSQLEditStr = strings.Join(sqlEdit, ",\n\t\t")
			result.MethodsProps[i].EditableFields = editableFields
		}
	}
	model = &result
	return nil
}

// IsStandardColumn check if column is one of column which description column with auto substituted when model is created or modified
func IsStandardColumn(column string) bool {
	if isCreatedStandardColumn(column) || isModifiedStandardColumn(column) {
		return true
	}
	return false
}

func isCreatedStandardColumn(column string) bool {
	if column == "createdAt" || column == "createdBy" {
		return true
	}
	return false
}

func isModifiedStandardColumn(column string) bool {
	if column == "modifiedAt" || column == "modifiedBy" {
		return true
	}
	return false
}
