package parser

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/mtgroupit/mt-back-generator/internal/models"
	"github.com/mtgroupit/mt-back-generator/internal/utilities"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var (
	isCorrectName          = regexp.MustCompile(`^[a-z][A-Za-z0-9]+$`).MatchString
	correctNameDescription = "A valid name must contain only letters and numbers in camelCase"

	roles = []string{"admin", "manager", "user", "guest"}
)

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
	for name, profileField := range cfg.AddProfileFields {
		field := models.ProfileField{
			Name: strings.Title(name),
			Type: profileField.Type,
		}
		err = field.SetBusinessType()
		if err != nil {
			return
		}
		cfg.AddProfileFields[name] = field
	}
	for customTypeName, customType := range cfg.CustomTypes {
		for field, options := range customType.Fields {
			if options.IsCustom, options.IsArray, options.BusinessType, err = parseFieldType(options, cfg.CustomTypes); err != nil {
				return nil, errors.Wrapf(err, `Custom type: "%s". Field: "%s"`, customTypeName, field)
			}

			if IsTypesAdditionalType(options.BusinessType) {
				cfg.HaveTypesInCustomTypes = true
			}
			if options.Default != "" {
				if IsTimeFormat(options.Format) || options.Format == "email" {
					cfg.HaveConvInCustomTypes = true
				} else {
					cfg.HaveSwagInCustomTypes = true
				}
			}
			if options.Required {
				if !IsTimeFormat(options.Format) && options.Format != "email" {
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
				options.Type = utilities.LowerTitle(options.BusinessType)
			} else {
				if options.Type == "int" {
					options.Type = "int32"
					options.BusinessType = "int32"
				}
				if options.IsArray {
					switch options.BusinessType {
					case "int":
						options.Type = "int32"
						options.BusinessType = "int32"
					case "float64":
						options.Type = "float"
					case TypesPrefix + "Decimal":
						options.Type = "decimal"
					case "date", "date-time":
						options.Type = "string"
					default:
						options.Type = options.BusinessType
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
		columns := SortColumns(model.Columns)

		for i := range model.Tags {
			model.Tags[i] = strings.Title(model.Tags[i])
		}
		model.Tags = append([]string{strings.Title(name)}, model.Tags...)

		var props []models.MethodProps
		for _, method := range model.Methods {
			noSecure := false
			if strings.Contains(method, "{noSecure}") {
				noSecure = true
			}
			if !model.IsStandardMethod(method) {
				cfg.HaveCustomMethod = true
				model.HaveCustomMethod = true
			}

			var prop models.MethodProps
			if method == "delete" || method == "deleteMy" {
				prop.HTTPMethod = "delete"
			} else if method == "edit" || method == "editMy" || models.IsAdjustEdit(method) {
				prop.HTTPMethod = "put"
			} else {
				prop.HTTPMethod = "post"
			}
			prop.NoSecure = noSecure

			if method == "list" || models.IsAdjustList(method) {
				cfg.HaveListMethod = true
				model.HaveListMethod = true
			}

			rules := []string{}
			for ruleName, ruleMethods := range model.RulesSet {
				if utilities.ContainsStr(ruleMethods, method) {
					rules = append(rules, ruleName)
				}
			}
			prop.Rules = rules

			props = append(props, prop)
		}
		model.MethodsProps = props

		psql := []models.PsqlParams{}
		var indexLastNotArrOfStruct int
		for _, column := range columns {
			options := model.Columns[column]
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
			if options.IsStruct, options.IsCustom, options.IsArray, options.BusinessType, err = parseColumnType(options, cfg); err != nil {
				return nil, errors.Wrapf(err, `Model: "%s". Column: "%s"`, name, column)
			}

			if IsTypesAdditionalType(options.BusinessType) {
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
			if options.Required {
				if options.Format == "email" {
					model.NeedConv = true
				}
			}

			if options.IsStruct {
				model.HaveLazyLoading = true

				modelNameForBind := utilities.LowerTitle(options.BusinessType)

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

					et.Name = utilities.NameSQL(name) + "_" + utilities.NameSQL(column)

					et.RefTableOne = strings.Title(name)
					et.RefIDOne = "id"
					et.FieldIDOne = utilities.NameSQL(name) + "_id"
					if cfg.Models[name].Columns["id"].Type == "uuid" {
						et.TypeIDOne = "uuid"
					} else {
						et.TypeIDOne = "integer"
					}

					et.RefTableTwo = options.BusinessType
					et.RefIDTwo = "id"
					et.FieldIDTwo = utilities.NameSQL(column) + "_id"
					if cfg.Models[utilities.LowerTitle(options.BusinessType)].Columns["id"].Type == "uuid" {
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
			pp.NotNull = options.Required
			if column == "id" {
				switch options.Type {
				case "uuid":
					options.BusinessType = "string"
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
					options.Type = utilities.LowerTitle(options.BusinessType)
				} else {
					if options.Type == "int" {
						options.Type = "int32"
						options.BusinessType = "int32"
					}
					if options.IsArray {
						switch options.BusinessType {
						case "int":
							options.Type = "int32"
							options.BusinessType = "int32"
						case "float64":
							options.Type = "float"
							cfg.HaveFloatArr = true
						case TypesPrefix + "Decimal":
							options.Type = "decimal"
						case "date", "date-time":
							options.Type = "string"
						default:
							options.Type = options.BusinessType
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

				pp.Type = options.BusinessType
				pp.Name = options.TitleName
				if pp.IsStruct {
					pp.SQLName = utilities.NameSQL(column) + "_id"
					pp.FK = "id"
					if cfg.Models[utilities.LowerTitle(options.BusinessType)].Columns["id"].Type == "uuid" {
						pp.TypeSQL = "uuid"
					} else {
						pp.TypeSQL = "integer"
					}
				} else if pp.IsCustom || pp.IsArray {
					pp.SQLName = utilities.NameSQL(column) + "_json"
					pp.TypeSQL = "jsonb"
				} else {
					pp.SQLName = utilities.NameSQL(column)
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
		for _, column := range columns {
			options := model.Columns[column]
			if !options.IsStruct {
				sqlName := utilities.NameSQL(options.TitleName)
				titleName := options.TitleName
				if options.IsArray || options.IsCustom {
					sqlName += "_json"
					titleName += "JSON"
				} else {
					titleName = "m." + titleName
				}
				SQLSelect = append(SQLSelect, sqlName)
				if !options.IsArray && !options.IsCustom {
					sqlColumn := utilities.NameSQL(column)
					if options.Type != "string" || IsTimeFormat(options.Format) || options.StrictFilter {
						sqlWhereParams = append(sqlWhereParams, fmt.Sprintf("(CAST(:%s as text) IS NULL OR %s=:%s) AND\n\t\t(CAST(:%s as text) IS NULL OR %s<>:%s)", sqlColumn, sqlName, sqlColumn, "not_"+sqlColumn, sqlName, "not_"+sqlColumn))
					} else {
						sqlWhereParams = append(sqlWhereParams, fmt.Sprintf("(CAST(:%s as text) IS NULL OR LOWER(%s) LIKE LOWER(:%s)) AND\n\t\t(CAST(:%s as text) IS NULL OR LOWER(%s) NOT LIKE LOWER(:%s))", sqlColumn, sqlName, sqlColumn, "not_"+sqlColumn, sqlName, "not_"+sqlColumn))
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
					sqlColumn := utilities.NameSQL(column)
					sqlName := utilities.NameSQL(options.TitleName) + "_id"
					SQLSelect = append(SQLSelect, sqlName)
					sqlWhereParams = append(sqlWhereParams, fmt.Sprintf("(CAST(:%s as text) IS NULL OR %s=:%s) AND\n\t\t(CAST(:%s as text) IS NULL OR %s<>:%s)", sqlColumn, sqlName, sqlColumn, "not_"+sqlColumn, sqlName, "not_"+sqlColumn))
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

	sort.Slice(cfg.ExtraTables, func(i, j int) bool { return cfg.ExtraTables[i].Name < cfg.ExtraTables[j].Name })

	titleize(cfg)

	return
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
		if options.IsStruct, _, options.IsArray, options.BusinessType, err = parseColumnType(options, cfg); err != nil {
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

// SortColumns returns sorted keys for columns map ("id" will be first, if it exist).
func SortColumns(columns map[string]models.Options) []string {
	keys := make([]string, 0, len(columns))
	for k := range columns {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		if strings.ToLower(keys[i]) == "id" {
			return true
		}
		if strings.ToLower(keys[j]) == "id" {
			return false
		}
		return keys[i] < keys[j]
	})

	return keys
}

func titleize(cfg *models.Config) {
	titleModels := make(map[string]models.Model)
	for modelName, model := range cfg.Models {
		titleModels[strings.Title(modelName)] = model
	}
	cfg.Models = titleModels
}

func fieldIsStruct(field string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]+\*{0,1}\(.+\)$`).Match([]byte(field))
}

func handleNestedObjs(modelsIn map[string]models.Model, modelName, elem, nesting, parent string, isArray bool) ([]models.NestedObjProps, error) {
	objs := []models.NestedObjProps{}
	obj := models.NestedObjProps{}

	obj.Shared = modelsIn[modelName].Shared

	field := models.ExtractName(elem)
	fieldsStr := models.ExtractStrNestedFields(elem)
	fieldsFull := models.SplitFields(fieldsStr)
	fields := models.TrimFieldsSuffix(fieldsFull)
	SQLSelect := []string{}
	haveID := false
	for i := range fields {
		var haveFieldInColumns bool
		var structModel string
		for column, options := range modelsIn[modelName].Columns {
			if column == fields[i] {
				structModel = utilities.LowerTitle(options.BusinessType)
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
					SQLSelect = append(SQLSelect, utilities.NameSQL(column))
				} else {
					if !options.IsArray {
						SQLSelect = append(SQLSelect, utilities.NameSQL(fields[i])+"_id")
					} else {
						structIsArr = true
					}
					obj.NeedLazyLoading = true
				}
			}
		}

		if fieldIsStruct(fieldsFull[i]) {
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
				subModelName := options.BusinessType
				fields := strings.Split(sort, ".")
				options.NestedSorts = append(options.NestedSorts, "."+strings.Title(column))
				for j, field := range fields {
					if field == "id" {
						return errors.Errorf(`Model: "%s". Column: "%s". Sort-by: "%s". Sorting by id is not avaliable`, modelName, column, sort)
					}
					subModel := modelsMap[utilities.LowerTitle(subModelName)]
					var ok bool
					var typeSortByColumn string
					for subColumn, subOptions := range subModel.Columns {
						if subColumn == field {
							if j == len(fields)-1 {
								switch subOptions.BusinessType {
								case "int32", "int64", "string":
									typeSortByColumn = strings.Title(subOptions.BusinessType)
								default:
									return errors.Errorf(`Model: "%s". Column: "%s". Sort-by: "%s". Type "%s" is not avaliable for sorting`, modelName, column, sort, subOptions.BusinessType)
								}
							}
							ok = true
							subModelName = subOptions.BusinessType
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
		if models.IsAdjustList(method) {
			var SQLSelect, sqlWhereParams, filtredFields []string
			fieldsStr := models.ExtractStrNestedFields(method)
			fieldsFull := models.SplitFields(fieldsStr)
			fields := models.TrimFieldsSuffix(fieldsFull)
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
							sqlName := utilities.NameSQL(options.TitleName)
							if options.IsArray || options.IsCustom {
								result.MethodsProps[i].HaveJSON = true
								sqlName += "_json"
							}
							SQLSelect = append(SQLSelect, sqlName)
							if needFilter && !options.IsArray {
								sqlColumn := utilities.NameSQL(column)
								if options.Type != "string" || IsTimeFormat(options.Format) || options.StrictFilter {
									sqlWhereParams = append(sqlWhereParams, fmt.Sprintf("(CAST(:%s as text) IS NULL OR %s=:%s) AND\n\t\t(CAST(:%s as text) IS NULL OR %s<>:%s)", sqlColumn, sqlName, sqlColumn, "not_"+sqlColumn, sqlName, "not_"+sqlColumn))
								} else {
									sqlWhereParams = append(sqlWhereParams, fmt.Sprintf("(CAST(:%s as text) IS NULL OR LOWER(%s) LIKE LOWER(:%s)) AND\n\t\t(CAST(:%s as text) IS NULL OR LOWER(%s) NOT LIKE LOWER(:%s))", sqlColumn, sqlName, sqlColumn, "not_"+sqlColumn, sqlName, "not_"+sqlColumn))
								}
							}
						} else {
							if !options.IsArray {
								sqlColumn := utilities.NameSQL(column)
								sqlName := utilities.NameSQL(options.TitleName) + "_id"
								SQLSelect = append(SQLSelect, sqlName)
								if needFilter {
									sqlWhereParams = append(sqlWhereParams, fmt.Sprintf("(CAST(:%s as text) IS NULL OR %s=:%s) AND\n\t\t(CAST(:%s as text) IS NULL OR %s<>:%s)", sqlColumn, sqlName, sqlColumn, "not_"+sqlColumn, sqlName, "not_"+sqlColumn))
								}
							} else {
								structIsArr = true
							}
						}
						result.MethodsProps[i].JSONColumns[column] = options.IsArray || options.IsCustom
					}
				}

				if fieldIsStruct(fieldsFull[j]) {
					result.MethodsProps[i].NeedLazyLoading = true

					objsForAdd, err := handleNestedObjs(modelsMap, structModel, fieldsFull[j], "", strings.Title(modelName), structIsArr)
					if err != nil {
						return err
					}
					result.MethodsProps[i].NestedObjs = append(result.MethodsProps[i].NestedObjs, objsForAdd...)
				}
			}

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
		if models.IsAdjustGet(method) {
			var SQLSelect, adjustGetJSONColumns []string
			haveID := false

			fieldsStr := models.ExtractStrNestedFields(method)
			fields := models.SplitFields(fieldsStr)
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
							sqlName := utilities.NameSQL(options.TitleName)
							if options.IsArray || options.IsCustom {
								sqlName += "_json"
							}
							SQLSelect = append(SQLSelect, sqlName)
						} else {
							if !options.IsArray {
								SQLSelect = append(SQLSelect, utilities.NameSQL(options.TitleName)+"_id")
							}
						}
					}
				}
			}

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
		if models.IsAdjustEdit(method) {
			var sqlEdit, editableFields []string
			count := 1
			if !model.Shared {
				count++
			}
			fieldsStr := models.ExtractStrNestedFields(method)
			fields := models.SplitFields(fieldsStr)
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
						sqlName := utilities.NameSQL(options.TitleName)
						if !options.IsStruct {
							if options.IsArray || options.IsCustom {
								sqlName += "_json"
							}
							sqlEdit = append(sqlEdit, fmt.Sprintf("%s=:%s", sqlName, sqlName))
						} else {
							if !options.IsArray {
								sqlName += "_id"
								sqlEdit = append(sqlEdit, fmt.Sprintf("%s=:%s", sqlName, sqlName))
							}
						}
					}
				}
			}
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
