package parser

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/mtgroupit/mt-back-generator/models"

	"github.com/fatih/camelcase"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var isCorectName = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]+$`).MatchString

var standardTypes = []string{"string", "int", "int32", "int64"}

func isStandardTypes(t string) bool {
	for i := range standardTypes {
		if t == standardTypes[i] {
			return true
		}
	}
	return false
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
	cfg = inCfg

	if cfg.Name == "" {
		return nil, errors.New("name is empty")
	}
	if cfg.Module == "" {
		return nil, errors.New("module is empty")
	}
	if cfg.AuthSrv == "" {
		return nil, errors.New("auth-srv is empty")
	}

	cfg.Description = strconv.Quote(cfg.Description)

	cfg.Tags = make(map[string]struct{})

	err = setDeepNesting(cfg)
	if err != nil {
		return
	}

	binds := map[string]models.Bind{}
	for name, model := range cfg.Models {
		if name == strings.Title(name) {
			return nil, errors.Errorf(`Model "%s" starts with captial letter, please rename it to "%s" starting with small letter`, name, LowerTitle(name))
		}
		if !isCorectName(name) {
			return nil, errors.Errorf(`"%s" is invalid name for model. A valid name must contain only letters and numbers in camelCase`, name)
		}
		model.TitleName = strings.Title(name)

		if len(model.Columns) == 0 {
			return nil, errors.Errorf(`Model "%s" has no any columns`, name)
		}

		for i := range model.Tags {
			cfg.Tags[LowerTitle(model.Tags[i])] = struct{}{}
			model.Tags[i] = strings.Title(model.Tags[i])
		}

		var props []models.MethodProps
		for _, method := range model.Methods {
			if isCustomMethod(method) {
				switch {
				case strings.HasPrefix(method, "list") && strings.Contains(method, "("):
					return nil, errors.Errorf(`Model: "%s". "%s"  is invalid as a custom list. A valid custom list shouldn't contain spaces before brackets. Correct method pattern: "list(column1, column3*, model1*(column1, model1(column1, column2))), where * means the field can be sorted by"`, name, method)
				case strings.HasPrefix(method, "edit") && strings.Contains(method, "("):
					return nil, errors.Errorf(`Model: "%s". "%s"  is invalid as a custom edit. A valid custom edit shouldn't contain spaces before brackets. Correct method pattern: "edit(column1, column2)"`, name, method)
				default:
					if !isCorectName(method) {
						return nil, errors.Errorf(`Model: "%s". "%s"  is invalid name for method. A valid name must contain only letters and numbers in camelCase`, name, method)
					}
				}
				cfg.HaveCustomMethod = true
				model.HaveCustomMethod = true
			}

			var prop models.MethodProps
			if method == "delete" {
				prop.HTTPMethod = "delete"
			} else if method == "edit" || isCustomEdit(method) {
				prop.HTTPMethod = "put"
			} else {
				prop.HTTPMethod = "post"
			}
			props = append(props, prop)

			if method == "list" || isCustomList(method) {
				cfg.HaveListMethod = true
				model.HaveListMethod = true
			}
		}
		model.MethodsProps = props

		psql := []models.PsqlParams{}
		var indexLastNotArrOfStruct int
		var haveDefaultSort bool
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
			if !isCorectName(column) {
				return nil, errors.Errorf(`Model: "%s". "%s"  is invalid name for column. A valid name must contain only letters and numbers in camelCase`, name, column)
			}
			if options.IsStruct, options.IsArray, options.GoType, err = checkColumn(options.Type, cfg); err != nil {
				return
			}

			if options.SortDefault {
				if options.IsStruct {
					return nil, errors.Errorf(`Model: "%s". Column: "%s". Structure can not be as default column for sorting`, name, column)
				}
				if options.IsArray {
					return nil, errors.Errorf(`Model: "%s". Column: "%s". Array can not be as default column for sorting`, name, column)
				}
				if !options.SortOn {
					return nil, errors.Errorf(`Model: "%s". Column "%s" can not be as default column for sorting because sorting is not enabled for this column`, name, column)
				}
				if haveDefaultSort {
					return nil, errors.Errorf(`Model "%s" has multiple columns as default for sorting, model should has one column as default for sorting`, name)
				}
				if options.SortOrderDefault != "" {
					orderDefault := strings.ToTitle(options.SortOrderDefault)
					if orderDefault == "ASC" || orderDefault == "DESC" {
						options.SortOrderDefault = orderDefault
					} else {
						return nil, errors.Errorf(`Model: "%s". Column: "%s". "%s" can not be as default order for sorting. Order for sorting can be only "ASC" or "DESC"`, name, column, options.SortOrderDefault)
					}
				}
				haveDefaultSort = true
			} else {
				if options.SortOrderDefault != "" {
					return nil, errors.Errorf(`Model: "%s". Column: "%s". Default order for sorting allow only for fields which set as default for sorting`, name, column)
				}
			}

			if options.StrictFilter {
				if options.Type != "string" && options.GoType != "string" {
					return nil, errors.Errorf(`Model: "%s". Column: "%s". "strict-sorting" option not available for non "string" columns`, name, column)
				}
			}

			if options.Format == "date-time" {
				cfg.HaveDateTime = true
			}
			if options.Format == "email" {
				model.HaveEmail = true
			}

			if options.Enum != "" {
				if isStandardTypes(options.Type) && column != "id" {
					if options.Type == "string" {
						if !regexp.MustCompile(`^\[\s*('.+'\s*,\s*)*'.+'+\s*\]$`).Match([]byte(options.Enum)) {
							return nil, errors.Errorf(`Model: "%s". Column: "%s". Enum for strings must be in this format: ['asc', 'desc'] and inputed as string`, name, column)
						}
					} else {
						if !regexp.MustCompile(`^\[\s*([0-9]+\s*,\s*)*[0-9]+\s*\]$`).Match([]byte(options.Enum)) {
							return nil, errors.Errorf(`Model: "%s". Column: "%s". Enum for types 'int', 'int32' and 'int64' must be in this format: [1, 2, 3] and inputed as string`, name, column)
						}
					}
				} else {
					return nil, errors.Errorf(`Model: "%s". Column: "%s". Enum available only for non id and types: %s`, name, column, standardTypes)
				}
			}

			if options.IsStruct {
				model.HaveLazyLoading = true
				binds[options.GoType] = models.Bind{
					ModelName: name,
					FieldName: column,
					IsArray:   options.IsArray,
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
				if options.IsArray {
					model.HaveArrayOfStandatrType = true
				}
			}

			pp := models.PsqlParams{}
			pp.IsArray = options.IsArray
			pp.IsStruct = options.IsStruct
			if column == "id" {
				switch options.Type {
				case "uuid":
					options.GoType = "string"
					model.IDIsUUID = true

					pp.Type = "string"
					pp.TypeSQL = "uuid"
				case "int64":
					options.GoType = options.Type
					options.Type = "integer"

					pp.Type = "int64"
					pp.TypeSQL = "SERIAL"
				default:
					return nil, errors.Errorf(`Model: "%s". "%s"  is invalid type for id. Valid types is 'int64' and 'uuid'`, name, options.Type)
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
				if options.IsStruct {
					options.Type = strings.ToLower(string(options.GoType[0])) + options.GoType[1:]
				} else {
					if options.Type == "int" {
						options.Type = "int32"
					}
					if !options.IsArray {
						options.GoType = options.Type
					} else {
						if options.GoType == "int" {
							options.Type = "int32"
							options.GoType = "int32"
						} else {
							options.Type = options.GoType
						}
					}
					switch options.Type {
					case "int32", "int64":
						options.Format = options.Type
						options.Type = "integer"
					case "bool":
						options.Type = "boolean"
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
				} else {
					if pp.IsArray {
						pp.SQLName = NameSQL(column) + "_json"
						pp.TypeSQL = "jsonb"
					} else {
						pp.SQLName = NameSQL(column)
						switch options.Type {
						case "string":
							pp.TypeSQL = "text"
						default:
							pp.TypeSQL = options.Type
						}
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

		var SQLSelect, sqlWhereParams, sqlAdd, sqlEdit, sqlAddExecParams, sqlEditExecParams, countFields []string
		count := 1
		countCreatedColumns := 0
		for column, options := range model.Columns {
			if !options.IsStruct {
				sqlName := NameSQL(options.TitleName)
				titleName := options.TitleName
				if options.IsArray {
					sqlName += "_json"
					titleName += "JSON"
				} else {
					titleName = "m." + titleName
				}
				SQLSelect = append(SQLSelect, sqlName)
				if !options.IsArray {
					if options.Type != "string" || options.StrictFilter {
						sqlWhereParams = append(sqlWhereParams, fmt.Sprintf(`((COALESCE(:%s, '1')='1' AND COALESCE(:%s, '2')='2') OR %s=:%s) AND ((COALESCE(:%s, '1')='1' AND COALESCE(:%s, '2')='2') OR %s<>:%s)`, column, column, sqlName, column, "not_"+column, "not_"+column, sqlName, "not_"+column))
					} else {
						sqlWhereParams = append(sqlWhereParams, fmt.Sprintf(`((COALESCE(:%s, '1')='1' AND COALESCE(:%s, '2')='2') OR LOWER(%s) LIKE LOWER(:%s)) AND ((COALESCE(:%s, '1')='1' AND COALESCE(:%s, '2')='2') OR LOWER(%s) NOT LIKE LOWER(:%s))`, column, column, sqlName, column, "not_"+column, "not_"+column, sqlName, "not_"+column))
					}
				}
				if options.TitleName != "ID" {
					sqlAdd = append(sqlAdd, sqlName)
					sqlAddExecParams = append(sqlAddExecParams, titleName)
					countFields = append(countFields, fmt.Sprintf("$%d", count))
					count++
					if isCreatedStandardColumn(column) {
						countCreatedColumns++
					} else {
						if model.Shared {
							sqlEdit = append(sqlEdit, fmt.Sprintf("%s=$%d", sqlName, count-countCreatedColumns))
						} else {
							sqlEdit = append(sqlEdit, fmt.Sprintf("%s=$%d", sqlName, (count-countCreatedColumns)+1))
						}
						sqlEditExecParams = append(sqlEditExecParams, titleName)
					}
				}
			} else {
				if !options.IsArray {
					if cfg.Models[LowerTitle(options.GoType)].Columns["id"].Type == "uuid" {
						SQLSelect = append(SQLSelect, "COALESCE("+NameSQL(options.TitleName)+"_id, '00000000-0000-0000-0000-000000000000') AS "+NameSQL(options.TitleName)+"_id")
					} else {
						SQLSelect = append(SQLSelect, "COALESCE("+NameSQL(options.TitleName)+"_id, 0) AS "+NameSQL(options.TitleName)+"_id")
					}
					sqlWhereParams = append(sqlWhereParams, fmt.Sprintf(`((COALESCE(:%s, '1')='1' AND COALESCE(:%s, '2')='2') OR %s=:%s) AND ((COALESCE(:%s, '1')='1' AND COALESCE(:%s, '2')='2') OR %s<>:%s)`, column, column, NameSQL(options.TitleName)+"_id", column, "not_"+column, "not_"+column, NameSQL(options.TitleName)+"_id", "not_"+column))
					sqlAdd = append(sqlAdd, NameSQL(options.TitleName)+"_id")
					sqlAddExecParams = append(sqlAddExecParams, column+"ID")
					sqlEditExecParams = append(sqlEditExecParams, column+"ID")
					countFields = append(countFields, fmt.Sprintf("$%d", count))
					count++
					sqlEdit = append(sqlEdit, fmt.Sprintf("%s_id=$%d", NameSQL(options.TitleName), count-countCreatedColumns))
				}
			}
		}
		model.SQLSelectStr = strings.Join(SQLSelect, ", ")
		model.SQLWhereParams = strings.Join(sqlWhereParams, " AND ")
		if model.IDIsUUID {
			sqlAdd = append(sqlAdd, "id")
			countFields = append(countFields, "$"+strconv.Itoa(len(countFields)+1))
		}
		if !model.Shared {
			sqlAdd = append(sqlAdd, "isolated_entity_id")
			countFields = append(countFields, "$"+strconv.Itoa(len(countFields)+1))
		}
		model.SQLAddStr = fmt.Sprintf("(%s) VALUES (%s)", strings.Join(sqlAdd, ", "), strings.Join(countFields, ", "))
		model.SQLEditStr = strings.Join(sqlEdit, ", ")
		model.SQLAddExecParams = strings.Join(sqlAddExecParams, ", ")
		model.SQLEditExecParams = strings.Join(sqlEditExecParams, ", ")

		cfg.Models[name] = model
	}
	for name, model := range cfg.Models {
		for bindModel, bind := range binds {
			if strings.Title(name) == bindModel {
				model.Binds = append(model.Binds, bind)
				break
			}
		}

		if err = handleSorts(cfg.Models, &model, name); err != nil {
			return
		}

		if err = handleCustomLists(cfg.Models, &model, name); err != nil {
			return
		}

		if err = handleCustomEdits(cfg.Models, &model, name); err != nil {
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
		if options.IsStruct, options.IsArray, options.GoType, err = checkColumn(options.Type, cfg); err != nil {
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

func checkColumn(columnType string, cfg *models.Config) (bool, bool, string, error) {
	switch {
	case strings.HasPrefix(columnType, "model."):
		if _, ok := cfg.Models[columnType[6:]]; !ok {
			return false, false, "", errors.Errorf(`One of the fields refers to "%s" model which is not described anywhere`, columnType[6:])
		}
		return true, false, strings.Title(columnType[6:]), nil
	case strings.HasPrefix(columnType, "[]model."):
		if _, ok := cfg.Models[columnType[8:]]; !ok {
			return false, false, "", errors.Errorf(`One of the fields refers to "%s" model which is not described anywhere`, columnType[8:])
		}
		return true, true, strings.Title(columnType[8:]), nil
	case strings.HasPrefix(columnType, "[]") && !strings.HasPrefix(columnType, "[]model."):
		t := columnType[2:]
		switch {
		case isStandardTypes(t):
			return false, true, t, nil
		default:
			return false, false, "", errors.Errorf(`"%s" is not correct type. You can use only one of standarat types %s or refers to any other model`, t, standardTypes)
		}

	}
	return false, false, "", nil
}

// NameSQL converts name to "snake_case" format
func NameSQL(name string) string {
	return strings.ToLower(strings.Join(camelcase.Split(name), "_"))
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

func isCustomList(method string) bool {
	return regexp.MustCompile(`^list\(.+\)$`).Match([]byte(method))
}

func isCustomEdit(method string) bool {
	return regexp.MustCompile(`^edit(My)?\(.+\)$`).Match([]byte(method))
}

func expandStrNestedFields(method string) (string, string) {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9]+\*{0,1}\((?P<value>.+)\)$`)
	result := []byte{}
	template := "$value"
	result = pattern.ExpandString(result, template, method, pattern.FindSubmatchIndex([]byte(method)))

	return strings.TrimSuffix(regexp.MustCompile("[^a-zA-Z0-9*]").Split(method, 2)[0], "*"), string(result)
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

	field, fieldsStr := expandStrNestedFields(elem)
	fieldsFull := splitFields(fieldsStr)
	fields := trimFieldsSuffix(fieldsFull)
	SQLSelect := []string{}
	haveID := false
	haveArr := false
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
			return nil, errors.Errorf(`Model "%s" does not contain "%s" column for custom list`, modelName, fields[i])
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
						if modelsIn[LowerTitle(options.GoType)].Columns["id"].Type == "uuid" {
							SQLSelect = append(SQLSelect, "COALESCE("+NameSQL(fields[i])+"_id, '00000000-0000-0000-0000-000000000000') AS "+NameSQL(fields[i])+"_id")
						} else {
							SQLSelect = append(SQLSelect, "COALESCE("+NameSQL(fields[i])+"_id, 0) AS "+NameSQL(fields[i])+"_id")
						}
					} else {
						haveArr = true
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
	if !haveID && haveArr {
		SQLSelect = append(SQLSelect, "id")
	}
	obj.SQLSelect = strings.Join(SQLSelect, ", ")
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

func handleCustomLists(modelsMap map[string]models.Model, model *models.Model, modelName string) error {
	result := *model
	for i, method := range result.Methods {
		if isCustomList(method) {
			var SQLSelect, sqlWhereParams, filtredFields []string
			_, fieldsStr := expandStrNestedFields(method)
			fieldsFull := splitFields(fieldsStr)
			fields := trimFieldsSuffix(fieldsFull)
			haveID := false
			haveArr := false
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
							SQLSelect = append(SQLSelect, NameSQL(column))
							if needFilter {
								if options.Type != "string" || options.StrictFilter {
									sqlWhereParams = append(sqlWhereParams, fmt.Sprintf(`((COALESCE(:%s, '1')='1' AND COALESCE(:%s, '2')='2') OR %s=:%s) AND ((COALESCE(:%s, '1')='1' AND COALESCE(:%s, '2')='2') OR %s<>:%s)`, column, column, NameSQL(options.TitleName), column, "not_"+column, "not_"+column, NameSQL(options.TitleName), "not_"+column))
								} else {
									sqlWhereParams = append(sqlWhereParams, fmt.Sprintf(`((COALESCE(:%s, '1')='1' AND COALESCE(:%s, '2')='2') OR LOWER(%s) LIKE LOWER(:%s)) AND ((COALESCE(:%s, '1')='1' AND COALESCE(:%s, '2')='2') OR LOWER(%s) NOT LIKE LOWER(:%s))`, column, column, NameSQL(options.TitleName), column, "not_"+column, "not_"+column, NameSQL(options.TitleName), "not_"+column))
								}
							}
						} else {
							if !options.IsArray {
								if modelsMap[LowerTitle(options.GoType)].Columns["id"].Type == "uuid" {
									SQLSelect = append(SQLSelect, "COALESCE("+NameSQL(options.TitleName)+"_id, '00000000-0000-0000-0000-000000000000') AS "+NameSQL(options.TitleName)+"_id")
								} else {
									SQLSelect = append(SQLSelect, "COALESCE("+NameSQL(options.TitleName)+"_id, 0) AS "+NameSQL(options.TitleName)+"_id")
								}
								if needFilter {
									sqlWhereParams = append(sqlWhereParams, fmt.Sprintf(`((COALESCE(:%s, '1')='1' AND COALESCE(:%s, '2')='2') OR %s=:%s) AND ((COALESCE(:%s, '1')='1' AND COALESCE(:%s, '2')='2') OR %s<>:%s)`, column, column, NameSQL(options.TitleName)+"_id", column, "not_"+column, "not_"+column, NameSQL(options.TitleName)+"_id", "not_"+column))
								}
							} else {
								haveArr = true
								structIsArr = true
							}
						}
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
			result.Methods[i] = "list" + strings.Join(fields, "")
			if !haveID && haveArr {
				SQLSelect = append(SQLSelect, "id")
			}
			result.MethodsProps[i].CustomListSQLSelect = strings.Join(SQLSelect, ", ")
			result.MethodsProps[i].CustomListSQLWhereProps = strings.Join(sqlWhereParams, " AND ")
			result.MethodsProps[i].FilteredFields = filtredFields
			result.MethodsProps[i].IsCustomList = true

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

func handleCustomEdits(modelsMap map[string]models.Model, model *models.Model, modelName string) error {
	result := *model
	for i, method := range result.Methods {
		if isCustomEdit(method) {
			var sqlEdit, sqlAddExecParams, editableFields []string
			count := 1
			methodName, fieldsStr := expandStrNestedFields(method)
			fields := splitFields(fieldsStr)
			for j := range fields {
				if IsStandardColumn(fields[j]) {
					return errors.Errorf(`Model "%s". Method: "%s". "%s" can not be used in custom edit method, it edits automatically`, modelName, method, fields[j])
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

				for column, options := range result.Columns {
					if fields[j] == options.TitleName {
						if !options.IsStruct {
							sqlAddExecParams = append(sqlAddExecParams, "m."+options.TitleName)
							count++
							sqlEdit = append(sqlEdit, fmt.Sprintf("%s=$%d", NameSQL(options.TitleName), count))
						} else {
							if !options.IsArray {
								sqlAddExecParams = append(sqlAddExecParams, column+"ID")
								count++
								sqlEdit = append(sqlEdit, fmt.Sprintf("%s_id=$%d", NameSQL(options.TitleName), count))
							}
						}
					}
				}
			}
			result.Methods[i] = methodName + strings.Join(fields, "")
			result.MethodsProps[i].CustomSQLEditStr = strings.Join(sqlEdit, ", ")
			result.MethodsProps[i].CustomSQLExecParams = strings.Join(sqlAddExecParams, ", ")
			result.MethodsProps[i].EditableFields = editableFields
		}
	}
	model = &result
	return nil
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

// isCustomMethod return true if method is custom
func isCustomMethod(method string) bool {
	method = strings.ToLower(method)
	if method == "get" || method == "add" || method == "delete" || method == "edit" || method == "list" || isCustomList(method) || isCustomEdit(method) || IsMyMethod(method) {
		return false
	}
	return true
}

// IsMyMethod return true if method is standard my method
func IsMyMethod(method string) bool {
	method = strings.ToLower(method)
	if method == "getmy" || method == "addmy" || method == "deletemy" || method == "editmy" || regexp.MustCompile(`^editmy.+`).Match([]byte(method)) {
		return true
	}
	return false
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
