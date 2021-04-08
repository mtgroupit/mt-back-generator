package models

import (
	"github.com/pkg/errors"
)

// Options contain properties of column
type Options struct {
	TitleName        string
	Type             string
	GoType           string
	Format           string
	Enum             []string
	Unique           bool
	Length           int64
	Default          string
	SortOn           bool     `yaml:"sort-on"`
	SortDefault      bool     `yaml:"sort-default"`
	SortOrderDefault string   `yaml:"sort-order-default"`
	SortBy           []string `yaml:"sort-by"`
	StrictFilter     bool     `yaml:"strict-filter"`

	NestedSorts []string
	IsStruct    bool
	IsCustom    bool
	IsArray     bool
	Pk          string
}

// PsqlParams - contain properties for postgres generate
type PsqlParams struct {
	Name     string
	SQLName  string
	Unique   bool
	Type     string
	TypeSQL  string
	IsArray  bool
	IsCustom bool
	IsStruct bool
	FK       string
	Last     bool
}

//Bind binds tables for build delete method
type Bind struct {
	ModelName string
	FieldName string
	IsArray   bool
}

// NestedObjProps contains properties of nested objects for custom list method
type NestedObjProps struct {
	Name                  string
	Type                  string
	SQLSelect             string
	Path                  string
	ParentStruct          string
	NeedLazyLoading       bool
	IsArray               bool
	IsFirstForLazyLoading bool
	IsLastForLazyLoading  bool
	Shared                bool
}

// MethodProps contains all additional information for method
type MethodProps struct {
	HTTPMethod              string
	IsAdjustList            bool
	NeedLazyLoading         bool
	HaveJSON                bool
	JSONColumns             map[string]bool
	AdjustSQLSelect         string
	AdjustGetJSONColumns    []string
	AdjustListSQLWhereProps string
	CustomSQLEditStr        string
	CustomSQLExecParams     string
	FilteredFields          []string
	EditableFields          []string
	NestedObjs              []NestedObjProps
}

// Model - description one component of models
type Model struct {
	Description           string
	Shared                bool
	Tags                  []string
	BoundToIsolatedEntity bool `yaml:"bind-to-isolated-entity"`
	DetailedPagination    bool `yaml:"detailed-pagination"`
	ReturnWhenEdit        bool `yaml:"return-when-edit"`
	Columns               map[string]Options

	TitleName        string
	Fields           map[string]string
	DeepNesting      int
	HaveLazyLoading  bool
	IDIsUUID         bool
	HaveEmail        bool
	NeedConv         bool
	NeedTypes        bool
	NeedTime         bool
	HaveListMethod   bool
	HaveCustomMethod bool
	HaveJSON         bool

	HaveCreatedAt  bool
	HaveCreatedBy  bool
	HaveModifiedAt bool
	HaveModifiedBy bool

	Psql              []PsqlParams
	SQLSelectStr      string
	SQLWhereParams    string
	SQLAddStr         string
	SQLEditStr        string
	SQLAddExecParams  string
	SQLEditExecParams string

	Binds        []Bind
	Methods      []string
	MethodsProps []MethodProps
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

// Config - description service, models and functions from yaml file
type Config struct {
	Name        string
	Module      string
	AuthSrv     string `yaml:"auth-srv"`
	Description string
	Debug       bool
	Models      map[string]Model
	CustomTypes map[string]CustomType `yaml:"custom-types"`
	Functions   map[string]Function

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

// CustomType - describes custom types
type CustomType struct {
	Description string
	Fields      map[string]Options

	NeedTime bool
}
