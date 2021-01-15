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
	"strings"
	"text/template"

	"github.com/mtgroupit/mt-back-generator/internal/parser"
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
			keys = append(keys, parser.NameSQL(iter.Key().String()+"s"))
		}
		return strings.Join(keys, ", ")
	},
	"ToLower": func(in string) string {
		return strings.ToLower(in)
	},
	"LowerTitle":   parser.LowerTitle,
	"Title":        strings.Title,
	"NameSQL":      parser.NameSQL,
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
	"IsAdd": func(method string) bool {
		method = strings.ToLower(method)
		if method == "add" || method == "addmy" {
			return true
		}
		return false
	},
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
	fmt.Println(cfg.Name)

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
						fileName := fmt.Sprintf("%05d_%s.sql", counter, parser.NameSQL(modelName))
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
				fileName := parser.NameSQL(modelName) + name[len("range_models"):len(name)-len(".gotmpl")]
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
