package parser

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/mtgroupit/mt-back-generator/french-back-template/generator/models"

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

		model.FirstLetter = string(name[0])
		model.TitleName = strings.Title(name)

		var httpMethods []string
		for _, method := range model.Methods {
			var httpMethod string
			switch method {
			case "edit":
				httpMethod = "put"
			case "delete":
				httpMethod = "delete"
			default:
				httpMethod = "post"
			}
			httpMethods = append(httpMethods, httpMethod)

			if method == "list" || IsCustomList(method) {
				cfg.HaveListMethod = true
			}
			if method == "filter" {
				cfg.HaveFilterMethod = true
			}
		}
		model.HTTPMethods = httpMethods

		psql := []models.PsqlParams{}
		for column, options := range model.Columns {

			if options.IsStruct, options.IsArray, options.GoType, err = checkColumn(options.Type, cfg); err != nil {
				return
			}
			if options.IsArray {
				model.LenParams--
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

				et.Name = name + "_" + column

				et.RefTableOne = strings.Title(name)
				et.RefIDOne = string(name[0]) + "_id"
				et.FieldIDOne = strings.ToLower(name) + "_id"

				et.RefTableTwo = options.GoType
				et.RefIDTwo = string(strings.ToLower(options.GoType)[0]) + "_id"
				et.FieldIDTwo = strings.ToLower(column) + "_id"

				cfg.ExtraTables = append(cfg.ExtraTables, et)
			}

			pp := models.PsqlParams{}
			pp.IsArray = options.IsArray
			pp.IsStruct = options.IsStruct
			if column == "id" {
				options.TitleName = "ID"
				options.Type = "integer"
				options.GoType = "int64"

				pp.Name = "ID"
				pp.SqlName = string(name[0]) + "_id"
				pp.Type = "int64"
				pp.TypeSql = "SERIAL"
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
					pp.SqlName = strings.ToLower(column) + "_id"
					pp.FK = string(strings.ToLower(options.Type)[0]) + "_id"
					pp.TypeSql = "integer"
				} else {
					pp.SqlName = string(name[0]) + "_" + strings.ToLower(column)
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
		}
		model.Psql = psql
		model.LenParams = len(psql) - 1

		var sqlSelect, sqlAdd, sqlEdit, sqlExexParams, countFields []string
		count := 1
		for _, options := range model.Columns {
			if !options.IsStruct {
				sqlSelect = append(sqlSelect, string(name[0])+"_"+strings.ToLower(options.TitleName))
				if options.TitleName != "ID" {
					sqlAdd = append(sqlAdd, string(name[0])+"_"+strings.ToLower(options.TitleName))
					sqlExexParams = append(sqlExexParams, "m."+options.TitleName)
					countFields = append(countFields, fmt.Sprintf("$%d", count))
					count++
					sqlEdit = append(sqlEdit, fmt.Sprintf("%s_%s=$%d", string(name[0]), strings.ToLower(options.TitleName), count))
				}
			} else {
				if !options.IsArray {
					sqlSelect = append(sqlSelect, "COALESCE("+strings.ToLower(options.TitleName)+"_id, 0) AS "+strings.ToLower(options.TitleName)+"_id")
					sqlAdd = append(sqlAdd, strings.ToLower(options.TitleName)+"_id")
					sqlExexParams = append(sqlExexParams, "m."+options.TitleName+".ID")
					countFields = append(countFields, fmt.Sprintf("$%d", count))
					count++
					sqlEdit = append(sqlEdit, fmt.Sprintf("%s_id=$%d", strings.ToLower(options.TitleName), count))
				}
			}
		}
		model.SqlSelectStr = strings.Join(sqlSelect, ", ")
		model.SqlAddStr = fmt.Sprintf("(%s) VALUES (%s)", strings.Join(sqlAdd, ", "), strings.Join(countFields, ", "))
		model.SqlEditStr = strings.Join(sqlEdit, ", ")
		model.SqlExecParams = strings.Join(sqlExexParams, ", ")

		var selectStrs []string
		var needListLazyLoadingSlice []bool
		for i, method := range model.Methods {
			var sqlSelect []string
			var needListLazyLoading bool
			if IsCustomList(method) {
				var fields []string
				fields, err = extractFields(method)
				if err != nil {
					return
				}
				for i := range fields {
					var haveFieldInColumns bool
					for column := range model.Columns {
						if column == fields[i] {
							haveFieldInColumns = true
						}
					}
					if !haveFieldInColumns {
						err = errors.Errorf(`model "%s" not contain "%s" column for method "%s"`, name, fields[i], method)
						return
					}

					if fields[i] == "id" {
						fields[i] = strings.ToUpper(fields[i])
					} else {
						fields[i] = strings.Title(fields[i])
					}

					for column, options := range model.Columns {
						if fields[i] == options.TitleName {
							if !options.IsStruct {
								sqlSelect = append(sqlSelect, string(name[0])+"_"+column)
							} else {
								needListLazyLoading = true
								if !options.IsArray {
									sqlSelect = append(sqlSelect, "COALESCE("+strings.ToLower(options.TitleName)+"_id, 0) AS "+strings.ToLower(options.TitleName)+"_id")
								}
							}
						}
					}
				}
				model.Methods[i] = "list" + strings.Join(fields, "")
			}
			selectStrs = append(selectStrs, strings.Join(sqlSelect, ", "))
			needListLazyLoadingSlice = append(needListLazyLoadingSlice, needListLazyLoading)
		}
		model.SqlSelectListStrs = selectStrs
		model.NeedListLazyLoading = needListLazyLoadingSlice

		cfg.Models[name] = model
	}
	for name, model := range cfg.Models {
		for bindModel, bind := range binds {
			if strings.Title(name) == bindModel {
				model.Binds = append(model.Binds, bind)
				break
			}
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
	if strings.HasPrefix(strings.ToLower(method), "list") && len(method) > 4 {
		return true
	}
	return false
}
func extractFields(method string) ([]string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return nil, err
	}
	var fields []string
	for _, field := range reg.Split(method, -1)[1:] {
		if field != "" {
			fields = append(fields, field)
		}
	}
	return fields, nil
}
