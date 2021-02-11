package gen

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/mtgroupit/mt-back-generator/internal/parser"
	"github.com/mtgroupit/mt-back-generator/models"

	"github.com/brianvoe/gofakeit"
)

func genApiTestValue(columnOptions models.Options) string {
	if columnOptions.GoType == parser.TypesPrefix+"Decimal" {
		columnOptions.GoType = "float64"
	}
	if columnOptions.IsArray {
		return genApiTestArray(columnOptions)
	}
	return genApiTestValueWithFormat(columnOptions)

}

func genAppTestValue(columnOptions models.Options) string {
	if columnOptions.IsArray {
		return genAppTestArray(columnOptions)
	}
	return genTestValue(columnOptions)
}

const (
	minLenthArray = 1
	maxLenthArray = 10
)

func genApiTestArray(columnOptions models.Options) string {
	var arr []string
	for i := gofakeit.Number(minLenthArray, maxLenthArray); i <= maxLenthArray; i++ {
		arr = append(arr, genApiTestValueWithFormat(columnOptions))
	}

	switch columnOptions.Format {
	case "date-time":
		return fmt.Sprintf("[]strfmt.DateTime{%s}", strings.Join(arr, ", "))
	case "email":
		return fmt.Sprintf("[]strfmt.Email{%s}", strings.Join(arr, ", "))
	default:
		return fmt.Sprintf("[]%s{%s}", columnOptions.GoType, strings.Join(arr, ", "))
	}
}

func genApiTestValueWithFormat(columnOptions models.Options) string {
	testValue := genTestValue(columnOptions)
	switch columnOptions.Format {
	case "date-time":
		testValue = fmt.Sprintf("toDateTime(%s)", testValue)
		if columnOptions.Default != "" {
			testValue = fmt.Sprintf("conv.DateTime(%s)", testValue)
		}
	case "email":
		testValue = fmt.Sprintf("strfmt.Email(%s)", testValue)
		if columnOptions.Default != "" {
			testValue = fmt.Sprintf("conv.Email(%s)", testValue)
		}
	default:
		if columnOptions.Default != "" {
			switch columnOptions.Format {
			case "float":
				testValue = fmt.Sprintf("swag.Float32(%s)", testValue)
			default:
				testValue = fmt.Sprintf("swag.%s(%s)", strings.Title(columnOptions.GoType), testValue)
			}
		}
	}
	return testValue
}

func genAppTestArray(columnOptions models.Options) string {
	var arr []string
	for i := gofakeit.Number(minLenthArray, maxLenthArray); i <= maxLenthArray; i++ {
		arr = append(arr, genTestValue(columnOptions))
	}
	return fmt.Sprintf("[]%s{%s}", columnOptions.GoType, strings.Join(arr, ", "))
}

func genTestValue(columnOptions models.Options) (str string) {
	gofakeit.Seed(time.Now().UnixNano())
	switch columnOptions.GoType {
	case "string":
		switch columnOptions.Format {
		case "date-time":
			dateTime := strfmt.NewDateTime()
			dateTime.Scan(gofakeit.Date())
			str = dateTime.String()
		case "email":
			str = gofakeit.Email()
		case "url":
			str = gofakeit.URL()
		case "phone":
			str = gofakeit.Phone()
		default:
			str = gofakeit.Word()
		}
		str = fmt.Sprintf(`"%s"`, str)
	case "int", "int32", "int64":
		str = fmt.Sprintf("%d", gofakeit.Int32())
	case "float32", "float64":
		str = fmt.Sprintf("%.2f", gofakeit.Float32Range(1.0, 1000.0))
	case "bool":
		str = fmt.Sprintf("%t", gofakeit.Bool())
	case parser.TypesPrefix + "Decimal":
		str = fmt.Sprintf("%sNewDecimal(%.2f)", parser.TypesPrefix, gofakeit.Float64Range(1.0, 1000.0))
	default:
		str = fmt.Sprintf("interface{}.(%s)", columnOptions.GoType)
	}
	return
}
