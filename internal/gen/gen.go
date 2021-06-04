package gen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/mtgroupit/mt-back-generator/internal/parser"
	"github.com/mtgroupit/mt-back-generator/internal/utilities"
	"github.com/rhysd/abspath"

	"github.com/mtgroupit/mt-back-generator/internal/models"
)

var goTmplFuncs = template.FuncMap{
	"Iterate": func(count int) []int {
		var items []int
		for i := 0; i <= count; i++ {
			items = append(items, i)
		}
		return items
	},
	"TruncateParams": func(in interface{}) string {
		var keys []string
		iter := reflect.ValueOf(in).MapRange()
		for iter.Next() {
			keys = append(keys, utilities.Pluralize(utilities.NameSQL(iter.Key().String())))
		}
		return strings.Join(keys, ", ")
	},
	"ToLower":    strings.ToLower,
	"LowerTitle": utilities.LowerTitle,
	"Title":      nameToTitle,

	"NameSQL":   utilities.NameSQL,
	"Pluralize": utilities.Pluralize,

	"IsTypesAdditionalType": parser.IsTypesAdditionalType,

	"ToAppMethodName": models.ToAppMethodName,
	"IsAdjustList":    models.IsAdjustList,
	"IsAdjustEdit":    models.IsAdjustEdit,
	"IsAdjustGet":     models.IsAdjustGet,
	"IsTimeFormat":    parser.IsTimeFormat,
	"HaveField": func(method, modelName string) bool {
		return strings.Contains(method, modelName)
	},
	"ContainsStr": utilities.ContainsStr,
	"IsMyMethod":  models.IsMyMethod,
	"HaveMyMethod": func(methods []string) bool {
		for _, method := range methods {
			if models.IsMyMethod(method) {
				return true
			}
		}
		return false
	},
	"FormatName":       formatName,
	"IsStandardColumn": parser.IsStandardColumn,
	"IsNotStrictFilter": func(model models.Model, column string) bool {
		if model.Columns[column].Type != "string" || model.Columns[column].StrictFilter {
			return false
		}
		return true
	},
	// HaveColumnWithModelThatIsStructAndIsArray is function for finding model in columns with column that struct or/and array
	// For withStruct=true, withArray=true "HaveColumnWithModelThatIsStructAndIsArray" will return true only if columns have column that is  model with column that is struct AND array
	"HaveColumnWithModelThatIsStructAndIsArray": func(columns map[string]models.Options, models map[string]models.Model, isStruct, isArray bool) bool {
		for _, options := range columns {
			for modelName2, model2 := range models {
				if options.BusinessType == modelName2 {
					for _, options2 := range model2.Columns {
						if options2.IsStruct == isStruct && options2.IsArray == isArray {
							return true
						}
					}
				}
			}
		}
		return false
	},
	"NeedJSONInsideColumns": func(columns map[string]models.Options, models map[string]models.Model) bool {
		for _, options := range columns {
			for modelName2, model2 := range models {
				if options.BusinessType == modelName2 {
					for _, options2 := range model2.Columns {
						if (!options2.IsStruct && options2.IsArray) || options2.IsCustom {
							return true
						}
					}
				}
			}
		}
		return false
	},
	"IsGet": func(method string) bool {
		method = models.CleanMethodsOptions(method)
		method = strings.ToLower(method)
		if method == "get" || method == "getmy" || models.IsAdjustGet(method) {
			return true
		}
		return false
	},
	"IsAdd": isAdd,
	"IsDelete": func(method string) bool {
		method = models.CleanMethodsOptions(method)
		method = strings.ToLower(method)
		if method == "delete" || method == "deletemy" {
			return true
		}
		return false
	},
	"IsEdit": func(method string) bool {
		method = models.CleanMethodsOptions(method)
		method = strings.ToLower(method)
		if method == "edit" || method == "editmy" || method == "editoraddmy" || models.IsAdjustEdit(method) {
			return true
		}
		return false
	},
	"IsList":           isList,
	"IsNoSecureMethod": models.IsNoSecureMethod,
	"IsValidateMethod": models.IsValidateMethod,
	"HaveListWithWarn": func(model models.Model) bool {
		for i, method := range model.Methods {
			if isList(method) {
				if model.HaveLazyLoading {
					if models.IsAdjustList(method) {
						if model.MethodsProps[i].NeedLazyLoading {
							return true
						}
					} else {
						return true
					}
				}
			}
		}
		return false
	},
	"NeedCustomFilter": func(model models.Model, method string) bool {
		if isList(method) {
			for i, method2 := range model.Methods {
				if method == method2 {
					if model.HaveLazyLoading {
						if models.IsAdjustList(method) {
							if len(model.MethodsProps[i].FilteredFields) != 0 || model.MethodsProps[i].HaveJSON {
								return true
							}
							if !model.MethodsProps[i].NeedLazyLoading {
								return false
							}
						}
					}
					return true
				}
			}
		}
		return false
	},
	"NeedErrorsInApi": func(model models.Model) bool {
		for _, method := range model.Methods {
			if !isList(method) && !(isAdd(method) && !models.IsMyMethod(method)) {
				return true
			}
		}
		for _, options := range model.Columns {
			if options.IsStruct && options.IsArray {
				return true
			}
		}
		return false
	},
	"ConvertApiToAppColumn": func(sourceStructName, column string, columnOptions models.Options) string {
		appValue := sourceStructName + "." + nameToTitle(column)

		if columnOptions.IsCustom {
			var s string
			if columnOptions.IsArray {
				s = "s"
			}
			return fmt.Sprintf("app%s%s(%s)", columnOptions.BusinessType, s, appValue)
		}

		switch columnOptions.Format {
		case "date-time":
			if columnOptions.IsArray {
				appValue = fmt.Sprintf("fromDateTimesArray(%s)", appValue)
			} else {
				appValue = fmt.Sprintf("(*time.Time)(%s)", appValue)
			}
		case "date":
			if columnOptions.IsArray {
				appValue = fmt.Sprintf("fromDatesArray(%s)", appValue)
			} else {
				appValue = fmt.Sprintf("(*time.Time)(%s)", appValue)
			}
		case "email":
			if columnOptions.IsArray {
				appValue = fmt.Sprintf("fromEmailsArray(%s)", appValue)
			} else {
				if columnOptions.Default != "" || columnOptions.Required {
					appValue = fmt.Sprintf("conv.EmailValue(%s)", appValue)
				}
				appValue = fmt.Sprintf("%s.String()", appValue)
			}
		default:
			if (columnOptions.Default != "" || columnOptions.Required) && !columnOptions.IsArray {
				switch columnOptions.BusinessType {
				case "float64":
					appValue = fmt.Sprintf("swag.Float32Value(%s)", appValue)
				case parser.TypesPrefix + "Decimal":
					appValue = fmt.Sprintf("swag.Float64Value(%s)", appValue)
				default:
					appValue = fmt.Sprintf("swag.%sValue(%s)", strings.Title(columnOptions.BusinessType), appValue)
				}
			}
		}

		if columnOptions.BusinessType == "float64" {
			if columnOptions.IsArray {
				appValue = fmt.Sprintf("float32to64Array(%s)", appValue)
			} else {
				appValue = fmt.Sprintf("float64(%s)", appValue)
			}
		}

		if columnOptions.BusinessType == parser.TypesPrefix+"Decimal" {
			if columnOptions.IsArray {
				appValue = fmt.Sprintf("%sFloat64ToDecimalsArray(%s)", parser.TypesPrefix, appValue)
			} else {
				appValue = fmt.Sprintf("%sNewDecimal(%s)", parser.TypesPrefix, appValue)
			}
		}

		return appValue
	},
	"ConvertAppToApiColumn": func(sourceStructName, column string, columnOptions models.Options) string {
		apiValue := sourceStructName + "." + nameToTitle(column)

		if columnOptions.IsStruct || columnOptions.IsCustom {
			var s string
			if columnOptions.IsArray {
				s = "s"
			}
			return fmt.Sprintf("api%s%s(%s)", columnOptions.BusinessType, s, apiValue)
		}

		switch columnOptions.BusinessType {
		case "float64":
			if columnOptions.IsArray {
				apiValue = fmt.Sprintf("float64to32Array(%s)", apiValue)
			} else {
				apiValue = fmt.Sprintf("float32(%s)", apiValue)
			}
		case parser.TypesPrefix + "Decimal":
			if columnOptions.IsArray {
				apiValue = fmt.Sprintf("%sDecimalsToFloat64Array(%s)", parser.TypesPrefix, apiValue)
			} else {
				apiValue = fmt.Sprintf("%s.Float64()", apiValue)
			}
		}

		switch columnOptions.Format {
		case "date-time":
			if columnOptions.IsArray {
				apiValue = fmt.Sprintf("toDateTimesArray(%s)", apiValue)
			} else {
				apiValue = fmt.Sprintf("(*strfmt.DateTime)(%s)", apiValue)
			}
		case "date":
			if columnOptions.IsArray {
				apiValue = fmt.Sprintf("toDatesArray(%s)", apiValue)
			} else {
				apiValue = fmt.Sprintf("(*strfmt.Date)(%s)", apiValue)
			}
		case "email":
			if columnOptions.IsArray {
				apiValue = fmt.Sprintf("toEmailsArray(%s)", apiValue)
			} else {
				apiValue = fmt.Sprintf("strfmt.Email(%s)", apiValue)
				if columnOptions.Default != "" || columnOptions.Required {
					apiValue = fmt.Sprintf("conv.Email(%s)", apiValue)
				}
			}
		default:
			if (columnOptions.Default != "" || columnOptions.Required) && !columnOptions.IsArray {
				switch columnOptions.BusinessType {
				case "float64":
					apiValue = fmt.Sprintf("swag.Float32(%s)", apiValue)
				case parser.TypesPrefix + "Decimal":
					apiValue = fmt.Sprintf("swag.Float64(%s)", apiValue)
				default:
					apiValue = fmt.Sprintf("swag.%s(%s)", strings.Title(columnOptions.BusinessType), apiValue)
				}
			}
		}

		return apiValue
	},
	"ConvertDalToAppColumn": func(columnOptions models.Options) string {
		appValue := "m." + columnOptions.TitleName

		if columnOptions.IsCustom {
			if columnOptions.IsArray {
				return fmt.Sprintf("app.%ssToPointers(%s)", columnOptions.BusinessType, appValue)
			}
			return fmt.Sprintf("&%s", appValue)
		}

		if columnOptions.IsStruct {
			var s string
			if columnOptions.IsArray {
				s = "s"
			}
			return fmt.Sprintf("app%s%s(%s)", columnOptions.BusinessType, s, appValue)
		} else {
			if columnOptions.IsArray {
				return appValue
			}
		}

		switch {
		case columnOptions.TitleName == "ID" && columnOptions.Type == "uuid":
			appValue = fmt.Sprintf("%s.String()", appValue)
		case parser.IsTypesAdditionalType(columnOptions.BusinessType):
			appValue = fmt.Sprintf("%s.Decimal", appValue)
		case parser.IsTimeFormat(columnOptions.Format):
			appValue = fmt.Sprintf("%s", appValue)
		default:
			appValue = fmt.Sprintf("%s.%s", appValue, strings.Title(columnOptions.BusinessType))
		}

		return appValue
	},
	"DalType": func(psqlType string) string {
		switch psqlType {
		case parser.TypesPrefix + "Decimal":
			return fmt.Sprintf("%sNullDecimal", parser.TypesPrefix)
		case "date", "date-time", "*time.Time":
			return "*time.Time"
		default:
			return fmt.Sprintf("sql.Null%s", strings.Title(psqlType))
		}
	},
	"GenApiTestValue": genApiTestValue,
	"GenAppTestValue": genAppTestValue,
	"EnumPrint": func(enum []string) string {
		return fmt.Sprintf(`[%s]`, strings.Join(enum, ", "))
	},
	"GenRule": func(rule models.Rule) string {
		var checkRole, checkNotRole, checkAttr []string
		for _, role := range rule.Roles {
			roleChecker := "r.attributes.Is" + strings.Title(role) + "(prof)"
			checkRole = append(checkRole, roleChecker)
			checkNotRole = append(checkNotRole, "!"+roleChecker)
		}
		for _, attr := range rule.Attributes {
			checkAttr = append(checkAttr, "r.attributes."+strings.Title(attr)+"(prof)")
		}
		return fmt.Sprintf("(%s) || ((%s) && %s)", strings.Join(checkNotRole, " && "), strings.Join(checkRole, " || "), strings.Join(checkAttr, " && "))
	},
	"GenRulesSet": func(rules []string) string {
		if len(rules) == 0 {
			return "true"
		}
		ruleMethods := []string{}
		for _, rule := range rules {
			ruleMethods = append(ruleMethods, "rs.rules."+nameToTitle(rule)+"(prof)")
		}
		return strings.Join(ruleMethods, " && ")
	},
	"SortColumns": parser.SortColumns,
	"AvailableKeys": func(props models.MethodProps) (result []string) {
		var availableKeys []string
		for key := range props.AvailableFilterKeys {
			availableKeys = append(availableKeys, key)
		}
		for _, nestProps := range props.NestedObjs {
			for key := range nestProps.AvailableFilterKeys {
				availableKeys = append(availableKeys, key)
			}
		}

		sort.Slice(availableKeys, func(i, j int) bool {
			return strings.Count(availableKeys[i], ".") < strings.Count(availableKeys[j], ".")
		})
		for i, key := range availableKeys {
			if strings.Contains(key, ".") {
				haveKeyPath := false
				keyArr := strings.Split(key, ".")
				for _, key2 := range availableKeys[:i] {
					if strings.TrimSuffix(key, "."+keyArr[len(keyArr)-1]) == key2 {
						haveKeyPath = true
						break
					}
				}
				if haveKeyPath {
					result = append(result, key)
				}
			} else {
				result = append(result, key)
			}
		}
		return
	},
}

// Srv - generate dir with service
func Srv(dir string, cfg *models.Config) error {
	abs, err := abspath.ExpandFrom(dir)
	if err != nil {
		return err
	}
	dir = abs.String()
	if err := ensureDir(dir, ""); err != nil {
		return err
	}

	if err := buildTreeDirs(dir, cfg.Name); err != nil {
		return err
	}

	abs, err = abspath.ExpandFrom("~/mt-gen/templates/srv")
	if err != nil {
		return err
	}
	if err := gen(abs.String(), path.Join(dir, cfg.Name), *cfg); err != nil {
		return err
	}

	return nil
}

// gen recursively browses folder with templates and run exec function for them
func gen(dirTMPL, dirTarget string, cfg models.Config) error {
	files, err := ioutil.ReadDir(dirTMPL)
	if err != nil {
		return err
	}

	for _, info := range files {
		if info.IsDir() {
			if err := gen(path.Join(dirTMPL, info.Name()), path.Join(dirTarget, info.Name()), cfg); err != nil {
				return err
			}
		} else {
			if err := exec(info.Name(), dirTMPL, dirTarget, cfg); err != nil {
				return err
			}
		}
	}
	return nil
}

// exec generate "name" template  in "dirTarget" directory
func exec(name, dirTMPL, dirTarget string, cfg models.Config) error {
	tmp, err := template.New(name).Funcs(goTmplFuncs).ParseFiles(path.Clean(path.Join(dirTMPL, name)))
	if err != nil {
		return err
	}

	if strings.HasPrefix(name, "range_models") {
		switch {
		case strings.HasSuffix(name, ".sql.gotmpl"):
			counter := 0
			if cfg.Debug {
				counter = 1
			}
			for i := 0; i <= cfg.MaxDeepNesting; i++ {
				for modelName, model := range cfg.Models {
					if model.DeepNesting == i {
						counter++
						fileName := fmt.Sprintf("%05d_%s.sql", counter, utilities.NameSQL(modelName))
						cfg.CurModel = modelName
						if err := createFile(fileName, dirTMPL, dirTarget, cfg, tmp); err != nil {
							return err
						}
					}
				}
			}
		default:
			for modelName, model := range cfg.Models {
				if !model.HaveCustomMethod && strings.HasSuffix(name, "custom.go.gotmpl") {
					continue
				}
				if len(model.Methods) == 0 && strings.HasSuffix(name, "_test.go.gotmpl") {
					continue
				}
				fileName := utilities.NameSQL(modelName) + name[len("range_models"):len(name)-len(".gotmpl")]
				if strings.HasSuffix(name, "custom.go.gotmpl") && checkExistenseFile(path.Join(dirTarget, fileName)) {
					file, err := ioutil.ReadFile(path.Join(dirTarget, fileName))
					if err != nil {
						return err
					}
					for _, method := range model.Methods {
						if !model.IsStandardMethod(method) && !regexp.MustCompile(`func \(.+\) `+models.ToAppMethodName(method)+modelName).Match(file) {
							var pattern, tag string
							switch {
							case strings.HasSuffix(dirTMPL, "api"):
								pattern = apiPattern
								if len(model.Tags) != 0 {
									tag = utilities.LowerTitle(model.Tags[0])
								}
							case strings.HasSuffix(dirTMPL, "app"):
								pattern = appPattern
							case strings.HasSuffix(dirTMPL, "dal"):
								pattern = dalPattern
							}
							t := template.Must(template.New("func").Parse(pattern))
							var buf bytes.Buffer
							if err := t.Execute(&buf, struct {
								Method    string
								ModelName string
								Tag       string
							}{
								models.ToAppMethodName(method),
								modelName,
								tag,
							}); err != nil {
								return err
							}
							file = append(file, buf.Bytes()...)
							if err := ioutil.WriteFile(path.Join(dirTarget, fileName), file, 0644); err != nil {
								return nil
							}
						}
					}

				} else {
					cfg.CurModel = modelName
					if err := createFile(fileName, dirTMPL, dirTarget, cfg, tmp); err != nil {
						return err
					}
				}
			}
		}
	} else {
		if !cfg.HaveCustomMethod && strings.HasSuffix(name, "custom.go.gotmpl") && !strings.HasSuffix(dirTarget, "authorization") {
			return nil
		}
		if !cfg.Debug && name == "00001_reset_all.sql.gotmpl" {
			return nil
		}
		fileName := name[:len(name)-len(".gotmpl")]
		if strings.HasSuffix(name, "custom.go.gotmpl") && checkExistenseFile(path.Join(dirTarget, fileName)) {
			file, err := ioutil.ReadFile(path.Join(dirTarget, fileName))
			if err != nil {
				return err
			}
			if strings.HasSuffix(dirTarget, "authorization") {
				for _, attr := range cfg.AccessAttributes {
					attrName := nameToTitle(attr)
					if !regexp.MustCompile(`func \(attr \*attributes\) ` + attrName + `\(.*\) bool`).Match(file) {
						t := template.Must(template.New("func").Parse(fmt.Sprintf(attrPattern, attrName)))
						var buf bytes.Buffer
						if err := t.Execute(&buf, nil); err != nil {
							return err
						}
						file = append(file, buf.Bytes()...)
					}
				}
			} else {
				for modelName, model := range cfg.Models {
					for _, method := range model.Methods {
						if !model.IsStandardMethod(method) && !regexp.MustCompile(`\s`+models.ToAppMethodName(method)+modelName).Match(file) {
							parts := strings.SplitN(string(file), "\n}\n", 3)
							file = []byte(parts[0] + "\n\t" + models.ToAppMethodName(method) + modelName + `(prof Profile, m *` + modelName + ") error\n}\n" +
								parts[1] + "\n\t" + models.ToAppMethodName(method) + modelName + `(m *` + modelName + ") error\n}\n" +
								parts[2])
						}
					}
				}
			}
			if err := ioutil.WriteFile(path.Join(dirTarget, fileName), file, 0644); err != nil {
				return nil
			}
		} else {
			if err := createFile(fileName, dirTMPL, dirTarget, cfg, tmp); err != nil {
				return err
			}
		}
	}
	return nil
}

func createFile(name, dirTMPL, dirTarget string, cfg models.Config, tmp *template.Template) error {
	f, err := os.Create(path.Clean(path.Join(dirTarget, name)))
	if err != nil {
		return err
	}
	defer f.Close()

	if err = tmp.Execute(f, cfg); err != nil {
		return err
	}

	log.Printf("%s created", path.Clean(path.Join(dirTarget, name)))
	return nil
}

func buildTreeDirs(p, srvName string) error {
	if err := ensureDir(p, srvName); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName), "cmd"); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName, "cmd"), "main"); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName), "internal"); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName, "internal"), "app"); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName, "internal"), "api"); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName, "internal", "api"), "restapi"); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName, "internal"), "dal"); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName, "internal"), "def"); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName, "internal"), "types"); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName, "internal"), "authorization"); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName), "migration"); err != nil {
		return err
	}
	return nil
}

func ensureDir(p, dirName string) error {
	err := os.Mkdir(path.Clean(path.Join(p, dirName)), 0777)
	if err == nil || os.IsExist(err) {
		return nil
	}
	return err
}

func checkExistenseFile(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func formatName(name string) string {
	splitedName := regexp.MustCompile("[^a-zA-Z0-9]+").Split(name, -1)
	for i := range splitedName {
		splitedName[i] = strings.Title(splitedName[i])
	}
	return strings.Join(splitedName, "")
}

func isAdd(method string) bool {
	method = models.CleanMethodsOptions(method)
	method = strings.ToLower(method)
	if method == "add" || method == "addmy" {
		return true
	}
	return false
}

func isList(method string) bool {
	method = models.CleanMethodsOptions(method)
	method = strings.ToLower(method)
	return method == "list" || models.IsAdjustList(method)
}

func nameToTitle(name string) string {
	if strings.ToLower(name) == "id" || strings.ToLower(name) == "url" {
		return strings.ToTitle(name)
	}
	return strings.Title(name)
}

const (
	apiPattern = `
func (svc *service) {{.Method}}{{.ModelName}}(params {{if .Tag}}{{.Tag}}{{else}}operations{{end}}.{{.Method}}{{.ModelName}}Params, profile interface{}) middleware.Responder {
	return {{if .Tag}}{{.Tag}}{{else}}operations{{end}}.New{{.Method}}{{.ModelName}}OK()
}`
	appPattern = `
func (a *app) {{.Method}}{{.ModelName}}(prof Profile, m *{{.ModelName}}) error {
	return a.custRepo.{{.Method}}{{.ModelName}}(m)
}`
	dalPattern = `
func (a *CustomsRepo) {{.Method}}{{.ModelName}}(m *app.{{.ModelName}}) error {
	return nil
}`
	attrPattern = `
func (attr *attributes) %s() bool {
	return true
}`
)
