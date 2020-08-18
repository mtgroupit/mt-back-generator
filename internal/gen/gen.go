package gen

import (
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

	"github.com/mtgroupit/mt-back-generator/models"
)

var goTmplFuncs = template.FuncMap{
	"Iterate": func(count int) []int {
		var i int
		var Items []int
		for i = 0; i <= count; i++ {
			Items = append(Items, i)
		}
		return Items
	},
	"IterateReverse": func(count int) []int {
		var i int
		var Items []int
		for i = count; i >= 0; i-- {
			Items = append(Items, i)
		}
		return Items
	},
	"TruncateParams": func(in interface{}) string {
		var keys []string
		iter := reflect.ValueOf(in).MapRange()
		for iter.Next() {
			keys = append(keys, iter.Key().String()+"s")
		}
		return strings.Join(keys, ", ")
	},
	"ToLower": func(in string) string {
		return strings.ToLower(in)
	},
	"LowerTitle": parser.LowerTitle,
	"NameSQL":    parser.NameSQL,
	"IsCustomList": func(method string) bool {
		return regexp.MustCompile(`^(L|l)ist.+`).Match([]byte(method))
	},
	"HaveField": func(method, modelName string) bool {
		return strings.Contains(method, modelName)
	},
}

// Srv - generate dir with service
func Srv(dir string, cfg models.Config) error {
	if err := treeDirs(dir, "service"); err != nil {
		return err
	}

	if err := swagger(path.Join(dir, "service"), cfg); err != nil {
		return err
	}

	parser.Titleize(&cfg)

	if err := gen("./templates/srv", path.Join(dir, "service"), cfg); err != nil {
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

// exec generate "name".go file or "model name".go files (if dirTMPL is range_models.go.gotmpl) in "dirTarget" directory
func exec(name, dirTMPL, dirTarget string, cfg models.Config) error {
	tmp, err := template.New(name).Funcs(goTmplFuncs).ParseFiles(path.Clean(path.Join(dirTMPL, name)))
	if err != nil {
		return err
	}

	if strings.HasPrefix(name, "range_models") {
		switch {
		case strings.Contains(name, ".go."):
			for modelName := range cfg.Models {
				var fileName string
				switch {
				case strings.Contains(name, "_integration_test."):
					fileName = parser.NameSQL(modelName) + "_integration_test.go"
				case strings.Contains(name, "_test."):
					fileName = parser.NameSQL(modelName) + "_test.go"
				default:
					fileName = parser.NameSQL(modelName) + ".go"
				}
				cfg.CurModel = modelName
				if err := createFile(fileName, dirTMPL, dirTarget, cfg, tmp); err != nil {
					return err
				}
			}
		case strings.Contains(name, ".sql."):
			counter := 0
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
		}
	} else {
		if err := createFile(name[:len(name)-len(".gotmpl")], dirTMPL, dirTarget, cfg, tmp); err != nil {
			return err
		}
	}
	return nil
}

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

// swagger generate swagger.yaml file in dir directory
func swagger(dir string, cfg models.Config) error {
	tmp, err := template.New("swagger.gotmpl").Funcs(goTmplFuncs).ParseFiles("./templates/swagger/swagger.gotmpl")
	if err != nil {
		return err
	}

	f, err := os.Create(path.Clean(path.Join(dir, "swagger.yaml")))
	if err != nil {
		return err
	}
	defer f.Close()

	if err = tmp.Execute(f, cfg); err != nil {
		return err
	}

	return nil
}

func treeDirs(p, srvName string) error {
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
	} else {
		return err
	}
}
