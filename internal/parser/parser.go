package parser

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/mtgroupit/mt-back-generator/french-back-template/generator/models"

	"github.com/fatih/camelcase"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func readYAMLConfig(file string, cfg *models.Config) error {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		return err
	}
	return nil
}

// Cfg create models.Config from configFile
func Cfg(configFile string) (cfg models.Config, err error) {
	if err = readYAMLConfig(configFile, &cfg); err != nil {
		return
	}

	cfg.Description = strconv.Quote(cfg.Description)

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return
	}
	cfg.Name = reg.ReplaceAllString(cfg.Name, "")

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
	binds := map[string]models.Bind{}
	for name, model := range cfg.Models {
		model.TitleName = strings.Title(name)

		var props []models.MethodProps
		for _, method := range model.Methods {
			var prop models.MethodProps
			switch method {
			case "edit":
				prop.HTTPMethod = "put"
			case "delete":
				prop.HTTPMethod = "delete"
			default:
				prop.HTTPMethod = "post"
			}
			props = append(props, prop)

			if method == "list" || IsCustomList(method) {
				cfg.HaveListMethod = true
			}
		}
		model.MethodsProps = props

		psql := []models.PsqlParams{}
		var indexLastNotArr int
		for column, options := range model.Columns {
			if options.IsStruct, options.IsArray, options.GoType, err = checkColumn(options.Type, cfg); err != nil {
				return
			}

			if options.Format == "date-time" {
				cfg.HaveDateTime = true
			}
			if options.Format == "email" {
				cfg.HaveEmail = true
			}

			if options.IsStruct {
				model.HaveLazyLoading = true
				binds[options.GoType] = models.Bind{
					ModelName: name,
					FieldName: column,
					IsArray:   options.IsArray,
				}
			}

			if options.IsArray {
				et := models.ExtraTable{}

				et.Name = NameSQL(name) + "_" + NameSQL(column)

				et.RefTableOne = strings.Title(name)
				et.RefIDOne = "id"
				et.FieldIDOne = NameSQL(name) + "_id"
				fmt.Println("one", name, cfg.Models[name].Columns["id"].Type)
				if cfg.Models[name].Columns["id"].Type == "uuid" {
					et.TypeIDOne = "uuid"
				} else {
					et.TypeIDOne = "integer"
				}

				et.RefTableTwo = options.GoType
				et.RefIDTwo = "id"
				et.FieldIDTwo = NameSQL(column) + "_id"
				_, ok1 := cfg.Models[LowerTitle(options.GoType)]
				_, ok2 := cfg.Models[LowerTitle(options.GoType)].Columns["id"]
				fmt.Println("two", LowerTitle(options.GoType), cfg.Models[LowerTitle(options.GoType)].Columns["id"].Type, ok1, ok2)
				if cfg.Models[LowerTitle(options.GoType)].Columns["id"].Type == "uuid" {
					et.TypeIDTwo = "uuid"
				} else {
					et.TypeIDTwo = "integer"
				}

				cfg.ExtraTables = append(cfg.ExtraTables, et)
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
					pp.TypeSql = "uuid"
				default:
					options.Type = "integer"
					options.GoType = "int64"

					pp.Type = "int64"
					pp.TypeSql = "SERIAL"
				}
				options.TitleName = "ID"

				pp.Name = "ID"
				pp.SqlName = "id"
			} else {
				if column == "url" {
					options.TitleName = "URL"
				} else {
					options.TitleName = strings.Title(column)
				}
				if options.IsStruct {
					options.Type = strings.ToLower(string(options.GoType[0])) + options.GoType[1:]
				} else {
					options.GoType = options.Type
					switch options.Type {
					case "int", "int32", "int64":
						options.Format = options.Type
						options.Type = "integer"
					case "bool":
						options.Type = "boolean"
					}
				}

				pp.Type = options.GoType

				pp.Name = options.TitleName
				if pp.IsStruct {
					pp.SqlName = NameSQL(column) + "_id"
					pp.FK = "id"
					if cfg.Models[LowerTitle(options.GoType)].Columns["id"].Type == "uuid" {
						pp.TypeSql = "uuid"
					} else {
						pp.TypeSql = "integer"
					}
				} else {
					pp.SqlName = NameSQL(column)
					switch options.Type {
					case "string":
						pp.TypeSql = "text"
					default:
						pp.TypeSql = options.Type
					}
				}
			}

			model.Columns[column] = options

			pp.Unique = options.Unique

			psql = append(psql, pp)
			if !pp.IsArray {
				indexLastNotArr = len(psql) - 1
			}
		}
		psql[indexLastNotArr].Last = true
		model.Psql = psql

		var sqlSelect, sqlAdd, sqlEdit, sqlExexParams, countFields []string
		count := 1
		for _, options := range model.Columns {
			if !options.IsStruct {
				sqlSelect = append(sqlSelect, NameSQL(options.TitleName))
				if options.TitleName != "ID" {
					sqlAdd = append(sqlAdd, NameSQL(options.TitleName))
					sqlExexParams = append(sqlExexParams, "m."+options.TitleName)
					countFields = append(countFields, fmt.Sprintf("$%d", count))
					count++
					sqlEdit = append(sqlEdit, fmt.Sprintf("%s=$%d", NameSQL(options.TitleName), count))
				}
			} else {
				if !options.IsArray {
					if cfg.Models[LowerTitle(options.GoType)].Columns["id"].Type == "uuid" {
						sqlSelect = append(sqlSelect, "COALESCE("+NameSQL(options.TitleName)+"_id, '00000000-0000-0000-0000-000000000000') AS "+NameSQL(options.TitleName)+"_id")
					} else {
						sqlSelect = append(sqlSelect, "COALESCE("+NameSQL(options.TitleName)+"_id, 0) AS "+NameSQL(options.TitleName)+"_id")
					}
					sqlAdd = append(sqlAdd, NameSQL(options.TitleName)+"_id")
					sqlExexParams = append(sqlExexParams, "m."+options.TitleName+".ID")
					countFields = append(countFields, fmt.Sprintf("$%d", count))
					count++
					sqlEdit = append(sqlEdit, fmt.Sprintf("%s_id=$%d", NameSQL(options.TitleName), count))
				}
			}
		}
		model.SqlSelectStr = strings.Join(sqlSelect, ", ")
		if model.IDIsUUID{
			sqlAdd = append(sqlAdd, "id")
			countFields = append(countFields, "$"+strconv.Itoa(len(countFields)+1))
		}
		model.SqlAddStr = fmt.Sprintf("(%s) VALUES (%s)", strings.Join(sqlAdd, ", "), strings.Join(countFields, ", "))
		model.SqlEditStr = strings.Join(sqlEdit, ", ")
		model.SqlExecParams = strings.Join(sqlExexParams, ", ")

		cfg.Models[name] = model
	}
	for name, model := range cfg.Models {
		for bindModel, bind := range binds {
			if strings.Title(name) == bindModel {
				model.Binds = append(model.Binds, bind)
				break
			}
		}

		if err = handleCustomLists(cfg.Models, &model, name); err != nil {
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

	return
}

func checkColumn(columnType string, cfg models.Config) (bool, bool, string, error) {
	switch {
	case strings.HasPrefix(columnType, "model."):
		if _, ok := cfg.Models[columnType[6:]]; !ok {
			return false, false, "", errors.Errorf(`config not contain "%s" field`, columnType[6:])
		}
		return true, false, strings.Title(columnType[6:]), nil
	case strings.HasPrefix(columnType, "[]model."):
		if _, ok := cfg.Models[columnType[8:]]; !ok {
			return false, false, "", errors.Errorf(`config not contain "%s" field`, columnType[8:])
		}
		return true, true, strings.Title(columnType[8:]), nil
	}
	return false, false, "", nil
}

func countDeepNesting(model string, cfg models.Config) (int, error) {
	var err error
	deepNesting := 0
	for _, options := range cfg.Models[model].Columns {
		if options.IsStruct, options.IsArray, options.GoType, err = checkColumn(options.Type, cfg); err != nil {
			return 0, err
		}
		if options.IsStruct && !options.IsArray {
			columnDeepNesting := 1
			optTypeDeepNesting, err := countDeepNesting(options.Type[6:], cfg)
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

func NameSQL(name string) string {
	return strings.ToLower(strings.Join(camelcase.Split(name), "_"))
}

// Titleize makes models keys as titles
func Titleize(cfg *models.Config) {
	titleModels := make(map[string]models.Model)
	for modelName, model := range cfg.Models {
		for i := range cfg.Models[modelName].Methods {
			model.Methods[i] = strings.Title(model.Methods[i])
		}
		titleModels[strings.Title(modelName)] = model
	}
	cfg.Models = titleModels
}

func IsCustomList(method string) bool {
	return regexp.MustCompile(`^list\(.+\)$`).Match([]byte(method))
}

func expandStrNestedFields(method string) (string, string) {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9]+\((?P<value>.+)\)$`)
	result := []byte{}
	template := "$value"
	result = pattern.ExpandString(result, template, method, pattern.FindSubmatchIndex([]byte(method)))

	return regexp.MustCompile("[^a-zA-Z0-9]").Split(method, 2)[0], string(result)
}

func splitFields(fields string) []string {
	var result []string
	for {
		fields = strings.Trim(fields, ", ")
		if strings.Index(fields, ",") >= 0 {
			if strings.Index(fields, ",") < strings.Index(fields, "(") || strings.Index(fields, "(") == -1 {
				substrs := regexp.MustCompile("[^a-zA-Z0-9]+").Split(fields, 2)
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
		out = append(out, regexp.MustCompile("[^a-zA-Z0-9]").Split(fields[i], 2)[0])
	}
	return
}
func isStruct(method string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]+\(.+\)$`).Match([]byte(method))
}
func handleNestedObjs(modelsIn map[string]models.Model, modelName, elem, nesting, parent string, isArray bool) ([]models.NestedObjProps, error) {
	objs := []models.NestedObjProps{}
	obj := models.NestedObjProps{}

	field, fieldsStr := expandStrNestedFields(elem)
	fieldsFull := splitFields(fieldsStr)
	fields := trimFieldsSuffix(fieldsFull)
	sqlSelect := []string{}
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
			return nil, errors.Errorf(`model "%s" not contain "%s" column for custom method`, modelName, fields[i])
		}
		if strings.ToLower(fields[i]) == "id" {
			haveID = true
		}

		structIsArr := false
		for column, options := range modelsIn[modelName].Columns {
			if fields[i] == column {
				if !options.IsStruct {
					sqlSelect = append(sqlSelect, NameSQL(column))
				} else {
					if !options.IsArray {
						if modelsIn[LowerTitle(options.GoType)].Columns["id"].Type == "uuid" {
							sqlSelect = append(sqlSelect, "COALESCE("+NameSQL(fields[i])+"_id, '00000000-0000-0000-0000-000000000000') AS "+NameSQL(fields[i])+"_id")
						} else {
							sqlSelect = append(sqlSelect, "COALESCE("+NameSQL(fields[i])+"_id, 0) AS "+NameSQL(fields[i])+"_id")
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
		sqlSelect = append(sqlSelect, "id")
	}
	obj.SqlSelect = strings.Join(sqlSelect, ", ")
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
func handleCustomLists(modelsMap map[string]models.Model, model *models.Model, modelName string) error {
	result := *model
	for i, method := range result.Methods {
		if IsCustomList(method) {
			var sqlSelect []string
			_, fieldsStr := expandStrNestedFields(method)
			fieldsFull := splitFields(fieldsStr)
			fields := trimFieldsSuffix(fieldsFull)
			haveID := false
			haveArr := false
			for j := range fields {
				var haveFieldInColumns bool
				var structModel string
				for column, options := range result.Columns {
					if column == fields[j] {
						structModel = options.Type
						haveFieldInColumns = true
					}
				}
				if !haveFieldInColumns {
					return errors.Errorf(`model "%s" not contain "%s" column for method "%s"`, modelName, fields[j], method)
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
							sqlSelect = append(sqlSelect, NameSQL(column))
						} else {
							if !options.IsArray {
								if modelsMap[LowerTitle(options.GoType)].Columns["id"].Type == "uuid" {
									sqlSelect = append(sqlSelect, "COALESCE("+NameSQL(options.TitleName)+"_id, '00000000-0000-0000-0000-000000000000') AS "+NameSQL(options.TitleName)+"_id")
								} else {
									sqlSelect = append(sqlSelect, "COALESCE("+NameSQL(options.TitleName)+"_id, 0) AS "+NameSQL(options.TitleName)+"_id")
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
				sqlSelect = append(sqlSelect, "id")
			}
			result.MethodsProps[i].CustomListSqlSelect = strings.Join(sqlSelect, ", ")
			result.MethodsProps[i].IsCustomList = true

			sort.Slice(result.MethodsProps[i].NestedObjs, func(a, b int) bool {
				return result.MethodsProps[i].NestedObjs[a].Path < result.MethodsProps[i].NestedObjs[b].Path
			})

			for j := range result.MethodsProps[i].NestedObjs {
				switch j {
				case 0:
					result.MethodsProps[i].NestedObjs[j].IsFirstForLazyLoading = true
					if len(result.MethodsProps[i].NestedObjs) == 1 {
						result.MethodsProps[i].NestedObjs[j].IsLastForLazyLoading = true
					}
				case len(result.MethodsProps[i].NestedObjs) - 1:
					result.MethodsProps[i].NestedObjs[j].IsLastForLazyLoading = true
					if len(result.MethodsProps[i].NestedObjs) == 2 {
						if result.MethodsProps[i].NestedObjs[j].Path != result.MethodsProps[i].NestedObjs[j-1].Path {
							result.MethodsProps[i].NestedObjs[j-1].IsLastForLazyLoading = true
							result.MethodsProps[i].NestedObjs[j].IsFirstForLazyLoading = true
						}
					}
				default:
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
