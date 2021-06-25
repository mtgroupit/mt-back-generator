package models

import (
	"regexp"
	"strings"
)

// IsStandardMethod return true if method is custom
func (model Model) IsStandardMethod(method string) bool {
	method = CleanMethodsOptions(method)
	if method == "get" || method == "add" || method == "delete" || method == "edit" || method == "list" || IsAdjustMethod(method) || IsMyMethod(method) {
		return true
	}
	// will be used for the next task
	//
	// if strings.HasPrefix(method, "addTo") || strings.HasPrefix(method, "deleteFrom") {
	// 	method = strings.TrimPrefix(method, "addTo")
	// 	method = strings.TrimPrefix(method, "deleteFrom")
	// 	if _, ok := model.Columns[utilities.LowerTitle(method)]; ok {
	// 		return true
	// 	}
	// }
	return false
}

// IsMyMethod return true if method is standard my method
func IsMyMethod(method string) bool {
	method = CleanMethodsOptions(method)
	method = strings.ToLower(method)
	if method == "getmy" || method == "addmy" || method == "deletemy" || method == "editmy" || method == "editoraddmy" || regexp.MustCompile(`^getmy.+`).Match([]byte(method)) || regexp.MustCompile(`^editmy.+`).Match([]byte(method)) {
		return true
	}
	return false
}

// IsAdjustMethod return true if method is adjusted get, edit or list
func IsAdjustMethod(method string) bool {
	return IsAdjustGet(method) || IsAdjustEdit(method) || IsAdjustList(method)
}

// IsAdjustGet return true if method is adjusted get
func IsAdjustGet(method string) bool {
	method = CleanMethodsOptions(method)
	return regexp.MustCompile(`^get(My|my)?\(.+\)(\[[a-zA-Z0-9]+\])?$`).Match([]byte(method))
}

// IsAdjustEdit return true if method is adjusted edit
func IsAdjustEdit(method string) bool {
	method = CleanMethodsOptions(method)
	return regexp.MustCompile(`^edit(My|my)?\(.+\)(\[[a-zA-Z0-9]+\])?$`).Match([]byte(method))
}

// IsAdjustList return true if method is adjusted list
func IsAdjustList(method string) bool {
	method = CleanMethodsOptions(method)
	return regexp.MustCompile(`^list\(.+\)(\[[a-zA-Z0-9]+\])?$`).Match([]byte(method))
}

// ToAppMethodName - converts method from config to method name which using in generated service
func ToAppMethodName(method string) string {
	method = CleanMethodsOptions(method)
	if IsAdjustMethod(method) {
		method = getNameForAdjustMethods(method)
	}
	method = strings.Title(method)
	return method
}

func getNameForAdjustMethods(method string) (result string) {
	methodName := ExtractName(method)
	methodNamePostfix := extractNamePostfixForAdjustMethods(method)

	if methodNamePostfix == "" {
		fieldsStr := ExtractStrNestedFields(method)
		fieldsFull := SplitFields(fieldsStr)
		fields := TrimFieldsSuffix(fieldsFull)
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

// CleanMethodsOptions returns method without options ("{noSecure}", "{validate}").
func CleanMethodsOptions(method string) string {
	if strings.Contains(method, "{noSecure}") {
		method = strings.Replace(method, "{noSecure}", "", -1)
	}
	if strings.Contains(method, "{validate}") {
		method = strings.Replace(method, "{validate}", "", -1)
	}
	return method
}

// IsNoSecureMethod returns true if method is no secure (has "{noSecure}" suffix).
func IsNoSecureMethod(method string) bool {
	return strings.Contains(method, "{noSecure}")
}

// IsValidateMethod returns true if method need the validation token (has "{validate}" suffix).
func IsValidateMethod(method string) bool {
	return strings.Contains(method, "{validate}")
}

// ExtractName - returns only method of adjusted method
func ExtractName(method string) string {
	return strings.TrimSuffix(regexp.MustCompile("[^a-zA-Z0-9*]").Split(method, 2)[0], "*")
}

func extractNamePostfixForAdjustMethods(method string) string {
	method = CleanMethodsOptions(method)
	pattern := regexp.MustCompile(`^[a-zA-Z0-9]+\*{0,1}\(.+\)(\[(?P<value>[a-zA-Z0-9]+)\])?$`)
	result := []byte{}
	template := "$value"
	result = pattern.ExpandString(result, template, method, pattern.FindSubmatchIndex([]byte(method)))

	return string(result)
}

// ExtractStrNestedFields - returns the contents of the round brackets for the adjusted method
func ExtractStrNestedFields(method string) string {
	method = CleanMethodsOptions(method)
	pattern := regexp.MustCompile(`^[a-zA-Z0-9]+\*{0,1}\((?P<value>.+)\)(\[[a-zA-Z0-9]+\])?$`)
	result := []byte{}
	template := "$value"
	result = pattern.ExpandString(result, template, method, pattern.FindSubmatchIndex([]byte(method)))

	return string(result)
}

// SplitFields - returns splited fields from round brackets in the adjusted method
func SplitFields(fields string) []string {
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

// SplitFields - returns trimed fields from round brackets in the adjusted method
func TrimFieldsSuffix(fields []string) (out []string) {
	for i := range fields {
		out = append(out, regexp.MustCompile("[^a-zA-Z0-9]").Split(fields[i], 2)[0])
	}
	return
}
