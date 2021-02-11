package models

import (
	"github.com/pkg/errors"
)

// Options contain properties of column
type Options struct {
	TitleName string
	Type      string
	GoType    string
	Format    string
	Enum      string

	SortOn           bool     `yaml:"sort-on"`
	SortDefault      bool     `yaml:"sort-default"`
	SortOrderDefault string   `yaml:"sort-order-default"`
	SortBy           []string `yaml:"sort-by"`

	StrictFilter bool `yaml:"strict-filter"`

	NestedSorts []string
	IsStruct    bool
	IsArray     bool
	Pk          string
	Unique      bool
	Length      int64
	Default     string
}

// PsqlParams - contain properties for postgres generate
type PsqlParams struct {
	Name     string
	SQLName  string
	Unique   bool
	Type     string
	TypeSQL  string
	IsArray  bool
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
	IsCustomList            bool
	NeedLazyLoading         bool
	HaveArrayOfStandardType bool
	ArrayColumns            map[string]bool
	CustomListSQLSelect     string
	CustomListSQLWhereProps string
	CustomSQLEditStr        string
	CustomSQLExecParams     string
	FilteredFields          []string
	EditableFields          []string
	NestedObjs              []NestedObjProps
}

// Model - description one component of models
type Model struct {
	Description          string
	Shared               bool
	DetailedPagination   bool `yaml:"detailed-pagination"`
	IDFromIsolatedEntity bool `yaml:"id-from-isolated-entity"`
	ReturnWhenEdit       bool `yaml:"return-when-edit"`
	Tags                 []string
	Columns              map[string]Options
	Fields               map[string]string
	Psql                 []PsqlParams
	DeepNesting          int

	HaveLazyLoading         bool
	IDIsUUID                bool
	HaveEmail               bool
	NeedConv                bool
	NeedTypes               bool
	HaveListMethod          bool
	HaveCustomMethod        bool
	HaveArrayOfStandardType bool

	SQLSelectStr      string
	SQLWhereParams    string
	SQLAddStr         string
	SQLEditStr        string
	SQLAddExecParams  string
	SQLEditExecParams string

	TitleName string

	HaveCreatedAt  bool
	HaveCreatedBy  bool
	HaveModifiedAt bool
	HaveModifiedBy bool

	Binds        []Bind
	Methods      []string
	MethodsProps []MethodProps
}

// Function contain input and output params
type Function struct {
	In          map[string]string
	Out         map[string]string
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
	Name             string
	Module           string
	AuthSrv          string `yaml:"auth-srv"`
	Description      string
	Debug            bool
	Models           map[string]Model
	Functions        map[string]Function
	HaveListMethod   bool
	HaveCustomMethod bool
	HaveDateTime     bool
	HaveTypes        bool
	HaveConv         bool
	HaveSwag         bool
	MaxDeepNesting   int

	ExtraTables []ExtraTable

	Tags map[string]struct{}
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
