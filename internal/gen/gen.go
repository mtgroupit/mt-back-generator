package gen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/mtgroupit/mt-back-generator/internal/parser"
	"github.com/mtgroupit/mt-back-generator/internal/shared"
	"github.com/rhysd/abspath"

	"github.com/mtgroupit/mt-back-generator/models"
)

var goTmplFuncs = template.FuncMap{
	"Iterate": func(count int) []int {
		var i int
		var items []int
		for i = 0; i <= count; i++ {
			items = append(items, i)
		}
		return items
	},
	"TruncateParams": func(in interface{}) string {
		var keys []string
		iter := reflect.ValueOf(in).MapRange()
		for iter.Next() {
			keys = append(keys, shared.NameSQL(iter.Key().String()+"s"))
		}
		return strings.Join(keys, ", ")
	},
	"ToLower": func(in string) string {
		return strings.ToLower(in)
	},
	"LowerTitle":   parser.LowerTitle,
	"Title":        strings.Title,
	"NameSQL":      shared.NameSQL,
	"IsCustomList": isCustomList,
	"IsCustomEdit": isCustomEdit,
	"HaveField": func(method, modelName string) bool {
		return strings.Contains(method, modelName)
	},
	"IsCustomMethod": isCustomMethod,
	"ContainsStr": func(slice []string, str string) bool {
		for i := range slice {
			if slice[i] == str {
				return true
			}
		}
		return false
	},
	"IsMyMethod": parser.IsMyMethod,
	"HaveMyMethod": func(methods []string) bool {
		for _, method := range methods {
			if parser.IsMyMethod(method) {
				return true
			}
		}
		return false
	},
	"FormatName":       formatName,
	"IsStandardColumn": parser.IsStandardColumn,
	"IsNotStrictFilter": func(model *models.Model, column string) bool {
		if model.Columns[column].Type != "string" || model.Columns[column].StrictFilter {
			return false
		}
		return true
	},
	// HaveColumnWithModelThatIsStructAndIsArray is function for finding model in columns with column that struct or/and array
	// For withStruct=true, withArray=true "HaveColumnWithModelThatIsStructAndIsArray" will return true only if columns have column that is  model with column that is struct AND array
	"HaveColumnWithModelThatIsStructAndIsArray": func(columns map[string]models.Options, models map[string]*models.Model, isStruct, isArray bool) bool {
		for _, options := range columns {
			for modelName2, model2 := range models {
				if options.GoType == modelName2 {
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
	"IsGet": func(method string) bool {
		method = strings.ToLower(method)
		if method == "get" || method == "getmy" {
			return true
		}
		return false
	},
	"IsAdd": isAdd,
	"IsDelete": func(method string) bool {
		method = strings.ToLower(method)
		if method == "delete" || method == "deletemy" {
			return true
		}
		return false
	},
	"IsEdit": func(method string) bool {
		method = strings.ToLower(method)
		if method == "edit" || method == "editmy" || method == "editoraddmy" || isCustomEdit(method) {
			return true
		}
		return false
	},
	"IsList": isList,
	"HaveListWithWarn": func(model models.Model) bool {
		for i, method := range model.Methods {
			if isList(method) {
				if model.HaveLazyLoading {
					if isCustomList(method) {
						if model.MethodsProps[i].NeedLazyLoading {
							return true
						}
					} else {
						return true
					}
				} else {
					return true
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
						if isCustomList(method) {
							if len(model.MethodsProps[i].FilteredFields) != 0 || model.MethodsProps[i].HaveArrayOfStandardType {
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
			if !isList(method) && !(isAdd(method) && !parser.IsMyMethod(method)) {
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
	"ConvertApiToAppColumn": func(sourceStructName string, columnOptions models.Options) string {
		appValue := sourceStructName + "." + columnOptions.TitleName

		switch columnOptions.Format {
		case "date-time":
			if columnOptions.IsArray {
				appValue = fmt.Sprintf("fromDateTimesArray(%s)", appValue)
			} else {
				if columnOptions.Default != "" {
					appValue = fmt.Sprintf("conv.DateTimeValue(%s)", appValue)
				}
				appValue = fmt.Sprintf("%s.String()", appValue)
			}
		case "email":
			if columnOptions.IsArray {
				appValue = fmt.Sprintf("fromEmailsArray(%s)", appValue)
			} else {
				if columnOptions.Default != "" {
					appValue = fmt.Sprintf("conv.EmailValue(%s)", appValue)
				}
				appValue = fmt.Sprintf("%s.String()", appValue)
			}
		default:
			if columnOptions.Default != "" {
				switch columnOptions.GoType {
				case "float64":
					appValue = fmt.Sprintf("swag.Float32Value(%s)", appValue)
				case parser.TypesPrefix + "Decimal":
					appValue = fmt.Sprintf("swag.Float64Value(%s)", appValue)
				default:
					appValue = fmt.Sprintf("swag.%sValue(%s)", strings.Title(columnOptions.GoType), appValue)
				}
			}
		}

		if columnOptions.GoType == "float64" {
			if columnOptions.IsArray {
				appValue = fmt.Sprintf("float32to64Array(%s)", appValue)
			} else {
				appValue = fmt.Sprintf("float64(%s)", appValue)
			}
		}

		if columnOptions.GoType == parser.TypesPrefix+"Decimal" {
			if columnOptions.IsArray {
				appValue = fmt.Sprintf("%sFloat64ToDecimalsArray(%s)", parser.TypesPrefix, appValue)
			} else {
				appValue = fmt.Sprintf("%sNewDecimal(%s)", parser.TypesPrefix, appValue)
			}
		}

		return appValue
	},
	"ConvertAppToApiColumn": func(columnOptions models.Options) string {
		apiValue := "a." + columnOptions.TitleName

		if columnOptions.IsStruct {
			var s string
			if columnOptions.IsArray {
				s = "s"
			}
			return fmt.Sprintf("api%s%s(%s)", columnOptions.GoType, s, apiValue)
		}

		switch columnOptions.GoType {
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
				apiValue = fmt.Sprintf("toDateTime(%s)", apiValue)
				if columnOptions.Default != "" {
					apiValue = fmt.Sprintf("conv.DateTime(%s)", apiValue)
				}
			}
		case "email":
			if columnOptions.IsArray {
				apiValue = fmt.Sprintf("toEmailsArray(%s)", apiValue)
			} else {
				apiValue = fmt.Sprintf("strfmt.Email(%s)", apiValue)
				if columnOptions.Default != "" {
					apiValue = fmt.Sprintf("conv.Email(%s)", apiValue)
				}
			}
		default:
			if columnOptions.Default != "" {
				switch columnOptions.GoType {
				case "float64":
					apiValue = fmt.Sprintf("swag.Float32(%s)", apiValue)
				case parser.TypesPrefix + "Decimal":
					apiValue = fmt.Sprintf("swag.Float64(%s)", apiValue)
				default:
					apiValue = fmt.Sprintf("swag.%s(%s)", strings.Title(columnOptions.GoType), apiValue)
				}
			}
		}

		return apiValue
	},
	"ConvertDalToAppColumn": func(columnOptions models.Options) string {
		appValue := "m." + columnOptions.TitleName

		if columnOptions.IsStruct {
			var s string
			if columnOptions.IsArray {
				s = "s"
			}
			return fmt.Sprintf("app%s%s(%s)", columnOptions.GoType, s, appValue)
		} else {
			if columnOptions.IsArray {
				return appValue
			}
		}

		switch {
		case columnOptions.TitleName == "ID" && columnOptions.Type == "uuid":
			appValue = fmt.Sprintf("%s.String()", appValue)
		case strings.HasPrefix(columnOptions.GoType, parser.TypesPrefix):
			appValue = fmt.Sprintf("%s(%s.Decimal)", columnOptions.GoType, appValue)
		default:
			appValue = fmt.Sprintf("%s.%s", appValue, strings.Title(columnOptions.GoType))
		}

		return appValue
	},
	"DalType": func(psqlType string) string {
		switch psqlType {
		case parser.TypesPrefix + "Decimal":
			return fmt.Sprintf("%sNullDecimal", parser.TypesPrefix)
		default:
			return fmt.Sprintf("sql.Null%s", strings.Title(psqlType))
		}
	},
	"GenApiTestValue": genApiTestValue,
	"GenAppTestValue": genAppTestValue,
	"MathAdd": func(a, b int) int {
		return a + b
	},
	// for debug purposes
	"Log": func(in interface{}) string {
		return fmt.Sprintf("%#v", in)
	},
}

// Srv - generate dir with service, if prevCfg is specified, then additional migrations will be generated
func Srv(dir string, cfg *models.Config, prevCfg *models.Config) error {
	abs, err := abspath.ExpandFrom(dir)
	if err != nil {
		return err
	}
	dir = abs.String()
	if err := ensureDir(dir, "", false); err != nil {
		return err
	}

	if err := buildTreeDirs(dir, cfg.Name); err != nil {
		return err
	}

	migVer, err := lastMigrationVersion(dir, cfg.Name)
	if err != nil {
		return err
	}
	cfg.LastMigrationVersion = migVer

	abs, err = abspath.ExpandFrom("~/mt-gen/templates/srv")
	if err != nil {
		return err
	}
	if err := gen(abs.String(), path.Join(dir, cfg.Name), cfg, prevCfg); err != nil {
		return err
	}

	return nil
}

func createFile(name, dirTMPL, dirTarget string, cfg *models.Config, tmp *template.Template) error {
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

// exec generate "name" template  in "dirTarget" directory
func exec(name, dirTMPL, dirTarget string, cfg *models.Config, prevCfg *models.Config) error {
	tmp, err := template.New(name).Funcs(goTmplFuncs).ParseFiles(path.Clean(path.Join(dirTMPL, name)))
	if err != nil {
		return err
	}

	if strings.HasPrefix(name, "zero_migration") {
		counter := 0
		// TODO
		// if cfg.Debug {
		// 	counter = 1
		// }
		err := cfg.ForEachModel(func(modelName string, model *models.Model) error {
			counter++
			fileName := fmt.Sprintf("%05d_%s.sql", counter, shared.NameSQL(modelName))
			cfg.CurModel = modelName
			cfg.NewModelObj = model
			return createFile(fileName, dirTMPL, dirTarget, cfg, tmp)
		})
		if err != nil {
			return err
		}
	} else if strings.HasPrefix(name, "range_models") {
		switch {
		case strings.HasSuffix(name, ".sql.gotmpl"):
			if prevCfg != nil {
				checkedModels := make(map[string]struct{})

				log.Printf("Previous configuration found, last DB version is %d\n", cfg.LastMigrationVersion)

				// check new config against previous one for new and modified models
				err := cfg.ForEachModel(func(newModelName string, newModel *models.Model) error {
					oldModel, ok := prevCfg.Models[newModelName]
					checkedModels[newModelName] = struct{}{}
					if !ok {
						// CREATE migration
						cfg.LastMigrationVersion++
						cfg.CurModel = newModelName
						cfg.NewModelObj = newModel
						cfg.OldModelObj = nil
						fileName := fmt.Sprintf("%05d_create_%s.sql", cfg.LastMigrationVersion, shared.NameSQL(cfg.CurModel))
						return createFile(fileName, dirTMPL, dirTarget, cfg, tmp)
					}
					if !newModel.Equals(*oldModel) {
						// ALTER migration
						cfg.LastMigrationVersion++
						cfg.CurModel = newModelName
						cfg.NewModelObj = newModel
						cfg.OldModelObj = oldModel
						fileName := fmt.Sprintf("%05d_alter_%s.sql", cfg.LastMigrationVersion, shared.NameSQL(cfg.CurModel))
						return createFile(fileName, dirTMPL, dirTarget, cfg, tmp)
					}
					return nil
				})
				if err != nil {
					return err
				}

				// check previous config against new one for deleted models
				err = prevCfg.ForEachModel(func(oldModelName string, oldModel *models.Model) error {
					if _, ok := checkedModels[oldModelName]; !ok {
						// DROP migration
						cfg.LastMigrationVersion++
						cfg.CurModel = oldModelName
						cfg.NewModelObj = nil
						cfg.OldModelObj = oldModel
						fileName := fmt.Sprintf("%05d_drop_%s.sql", cfg.LastMigrationVersion, shared.NameSQL(cfg.CurModel))
						return createFile(fileName, dirTMPL, dirTarget, cfg, tmp)
					}
					return nil
				})
				if err != nil {
					return err
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
				fileName := shared.NameSQL(modelName) + name[len("range_models"):len(name)-len(".gotmpl")]
				if strings.HasSuffix(name, "custom.go.gotmpl") && checkExistenseFile(path.Join(dirTarget, fileName)) {
					file, err := ioutil.ReadFile(path.Join(dirTarget, fileName))
					if err != nil {
						return err
					}
					for _, method := range model.Methods {
						if isCustomMethod(method) && !regexp.MustCompile(`func \(.+\) `+method+modelName).Match(file) {
							var pattern, tag string
							switch {
							case strings.HasSuffix(dirTMPL, "api"):
								pattern = apiPattern
								if len(model.Tags) != 0 {
									tag = parser.LowerTitle(model.Tags[0])
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
								method,
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
		if !cfg.HaveCustomMethod && strings.HasSuffix(name, "custom.go.gotmpl") {
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
			for modelName, model := range cfg.Models {
				for _, method := range model.Methods {
					if isCustomMethod(method) && !regexp.MustCompile(`\s`+method+modelName).Match(file) {
						file = []byte(regexp.MustCompile(`(?s)\n}\n`).ReplaceAllString(string(file), "\n\t"+method+modelName+`(m *`+modelName+") error\n}\n"))
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

// gen recursively browses folder with templates and run exec function for them
func gen(dirTMPL, dirTarget string, cfg *models.Config, prevCfg *models.Config) error {
	files, err := ioutil.ReadDir(dirTMPL)
	if err != nil {
		return err
	}

	for _, info := range files {
		if info.IsDir() {
			if err := gen(path.Join(dirTMPL, info.Name()), path.Join(dirTarget, info.Name()), cfg, prevCfg); err != nil {
				return err
			}
		} else {
			if err := exec(info.Name(), dirTMPL, dirTarget, cfg, prevCfg); err != nil {
				return err
			}
		}
	}
	return nil
}

// TODO make abstration for dir tree with options like "clear"
func buildTreeDirs(p, srvName string) error {
	if err := ensureDir(p, srvName, false); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName), "cmd", false); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName, "cmd"), "main", false); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName), "internal", false); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName, "internal"), "app", true); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName, "internal"), "api", true); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName, "internal", "api"), "restapi", true); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName, "internal"), "dal", true); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName, "internal"), "def", false); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName), "migration", false); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName, "internal"), "types", true); err != nil {
		return err
	}
	if err := ensureDir(path.Join(p, srvName), "zero-migration", true); err != nil {
		return err
	}
	return nil
}

func ensureDir(p, dirName string, clear bool) error {
	fullPath := path.Clean(path.Join(p, dirName))
	err := os.Mkdir(fullPath, 0777)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		if clear {
			return clearDir(fullPath)
		}
		return nil
	}
	return err
}

func clearDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func checkExistenseFile(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func isCustomMethod(method string) bool {
	method = strings.ToLower(method)
	if method == "get" || method == "add" || method == "delete" || method == "edit" || method == "list" || isCustomEdit(method) || isCustomList(method) || parser.IsMyMethod(method) {
		return false
	}
	return true
}

func isCustomList(method string) bool {
	return regexp.MustCompile(`^(L|l)ist.+`).Match([]byte(method)) && !parser.IsMyMethod(method)
}

func isCustomEdit(method string) bool {
	return regexp.MustCompile(`^(E|e)dit.+`).Match([]byte(method)) && strings.ToLower(method) != "editmy" && strings.ToLower(method) != "editoraddmy"
}

func formatName(name string) string {
	splitedName := regexp.MustCompile("[^a-zA-Z0-9]+").Split(name, -1)
	for i := range splitedName {
		splitedName[i] = strings.Title(splitedName[i])
	}
	return strings.Join(splitedName, "")
}

func isAdd(method string) bool {
	method = strings.ToLower(method)
	if method == "add" || method == "addmy" {
		return true
	}
	return false
}

func isList(method string) bool {
	method = strings.ToLower(method)
	return method == "list" || isCustomList(method)
}

const (
	apiPattern = `
func (svc *service) {{.Method}}{{.ModelName}}(params {{if .Tag}}{{.Tag}}{{else}}operations{{end}}.{{.Method}}{{.ModelName}}Params, profile interface{}) middleware.Responder {
	return {{if .Tag}}{{.Tag}}{{else}}operations{{end}}.New{{.Method}}{{.ModelName}}OK()
}`
	appPattern = `
func (a *app) {{.Method}}{{.ModelName}}(m *{{.ModelName}}) error {
	return a.cust.{{.Method}}{{.ModelName}}(m)
}`
	dalPattern = `
func (a *Customs) {{.Method}}{{.ModelName}}(m *app.{{.ModelName}}) error {
	return nil
}`
)

// TODO rework this function
func lastMigrationVersion(p, srvName string) (int, error) {
	dir := path.Join(p, srvName, "migration")
	file, err := os.Open(dir)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	ver := 0
	list, err := file.Readdirnames(0)
	if err != nil {
		return 0, err
	}
	for _, name := range list {
		if strings.HasSuffix(name, ".sql") {
			v, _ := strconv.Atoi(name[:5])
			if v > ver {
				ver = v
			}
		}
	}
	return ver, nil
}
