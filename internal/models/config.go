package models

import (
	"github.com/pkg/errors"
)

// Config - description service, models and functions from yaml file
type Config struct {
	Name             string
	Module           string
	AuthSrv          string `yaml:"auth-srv"`
	AuthPath         string `yaml:"auth-path"`
	Notifier         string
	Description      string
	Debug            bool
	Models           map[string]Model
	CustomTypes      map[string]CustomType `yaml:"custom-types"`
	Functions        map[string]Function
	AddProfileFields map[string]ProfileField `yaml:"add-profile-fields"`
	AccessAttributes []string                `yaml:"access-attributes"`
	Rules            map[string]Rule

	HaveListMethod         bool
	HaveCustomMethod       bool
	HaveFloatArr           bool
	HaveTime               bool
	HaveTimeInCustomTypes  bool
	HaveEmail              bool
	HaveEmailInCustomTypes bool
	HaveTypes              bool
	HaveTypesInCustomTypes bool
	HaveConv               bool
	HaveConvInCustomTypes  bool
	HaveSwag               bool
	HaveSwagInCustomTypes  bool
	MaxDeepNesting         int

	ExtraTables []ExtraTable

	// CurModel - field for define needed model for template
	CurModel string
}

// CustomType - describes custom types
type CustomType struct {
	Description string
	Fields      map[string]Options

	NeedTime bool
}

// RequiredFields -  returns slice of required fields
func (ct CustomType) RequiredFields() []string {
	var reqFields []string
	for field, options := range ct.Fields {
		if options.Required {
			reqFields = append(reqFields, field)
		}
	}
	return reqFields
}

// Function contain input and output params
type Function struct {
	In  map[string]string
	Out map[string]string

	InStr       string
	InStrParams string
	OutStr      string
	InStrType   string
	OutStrType  string
	InStrFull   string
	OutStrFull  string
	HaveOut     bool
}

// ProfileField fields need to be added to user profile
type ProfileField struct {
	Name         string
	Type         string
	BusinessType string
}

// AppType returns the field type for app
func (p ProfileField) AppType() string {
	switch p.BusinessType {
	case "string":
		return "string"
	case "int":
		return "int"
	case "date":
		return "*time.Time"
	case "date-time":
		return "*time.Time"
	default:
		return p.BusinessType
	}
}

// SetBusinessType sets the type of the business field
func (p *ProfileField) SetBusinessType() error {
	switch p.Type {
	case "string":
		p.BusinessType = "string"
	case "int":
		p.BusinessType = "int"
	case "date":
		p.BusinessType = "date"
	case "date-time":
		p.BusinessType = "date-time"
	default:
		return errors.Errorf(`unknown type "%s" field "%s"`, p.Type, p.Name)
	}
	return nil
}

// Rule - contains rule properties
type Rule struct {
	Attributes []string
	Roles      []string
}

// ExtraTable is table for many-to-many relations
type ExtraTable struct {
	Name string

	RefTableOne string
	RefIDOne    string
	FieldIDOne  string
	TypeIDOne   string

	RefTableTwo string
	RefIDTwo    string
	FieldIDTwo  string
	TypeIDTwo   string
}

//Bind binds tables for build delete method
type Bind struct {
	ModelName string
	FieldName string
	IsArray   bool
}

// AddBind - adding external bind to model with name 'nameModelTo'.
func (c *Config) AddBind(nameModelTo string, bind Bind) error {
	modelTo, ok := c.Models[nameModelTo]
	if !ok {
		return errors.Errorf(`Config has not "%s" model`, nameModelTo)
	}
	modelTo.Binds = append(modelTo.Binds, bind)
	c.Models[nameModelTo] = modelTo
	return nil
}
