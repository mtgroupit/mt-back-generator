package models

import (
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/mtgroupit/mt-back-generator/internal/shared"
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

	PrevColName string `yaml:"prev-col-name"`
}

// PsqlParams - contain properties for postgres generate
type PsqlParams struct {
	Name     string
	SQLName  string
	Unique   bool
	Type     string
	TypeSQL  string
	Model    *Model
	IsArray  bool
	IsCustom bool
	IsStruct bool
	// TODO this looks like Model.Binds, contains pointer to the model, if IsStruct == true
	FKModel *Model
	FK      string

	PrevColName string
}

// NewPsqlParams initializes new PsqlParams
func NewPsqlParams(model *Model, options Options) *PsqlParams {
	return &PsqlParams{
		Model:       model,
		IsArray:     options.IsArray,
		IsStruct:    options.IsStruct,
		PrevColName: options.PrevColName,
	}
}

// Equals says if current params have same values as in input params
func (p PsqlParams) Equals(_p PsqlParams) bool {
	return p.Name == _p.Name &&
		p.Type == _p.Type &&
		p.IsArray == _p.IsArray &&
		p.Unique == _p.Unique &&
		p.FK == _p.FK
}

// UniqueIdxName returns name of unique index for current model and column
func (p PsqlParams) UniqueIdxName() func() string {
	return func() (name string) {
		return p.Model.SQLTableName() + "_" + p.SQLName + "_key"
	}
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
	HaveJSON                bool
	JSONColumns             map[string]bool
	CustomListSQLSelect     string
	CustomListSQLWhereProps string
	CustomSQLEditStr        string
	CustomSQLExecParams     string
	FilteredFields          []string
	EditableFields          []string
	NestedObjs              []NestedObjProps
}

// ModelDifference represents sets of PsqlParams which shows model difference against other model (Create - new column, Delete - removed column, Update - changed column)
type ModelDifference struct {
	Create        []*PsqlParams
	Delete        []*PsqlParams
	Update        []*PsqlParams
	SharedChanged bool
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

	ExtraTables []*ExtraTable

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
	PsqlMap           map[string]*PsqlParams
	SQLSelectStr      string
	SQLWhereParams    string
	SQLAddStr         string
	SQLEditStr        string
	SQLAddExecParams  string
	SQLEditExecParams string

	// TODO this is used once in range_models.go.gotmpl, consider using FKModel
	Binds        []Bind
	Methods      []string
	MethodsProps []MethodProps
}

// SQLTableName returns name of SQL table for current model
func (m *Model) SQLTableName() string {
	return shared.NameSQL(m.TitleName) + "s"
}

// SQLAccessTableName returns name of SQL access table for current model
func (m *Model) SQLAccessTableName() string {
	return shared.NameSQL(m.TitleName) + "s_access"
}

// Equals says if current model has same table columns configuration as in input model
func (m *Model) Equals(_m Model) bool {
	diff := m.Difference(_m)
	return len(diff.Create) == 0 &&
		len(diff.Delete) == 0 &&
		len(diff.Update) == 0 &&
		!diff.SharedChanged
}

// Difference returns current model PsqlParams that differs with params of input model
func (m *Model) Difference(_m Model) ModelDifference {
	paramsA := m.PsqlMap
	paramsB := _m.PsqlMap
	differenceCreate := make([]*PsqlParams, 0)
	differenceUpdate := make([]*PsqlParams, 0)
	differenceDelete := make([]*PsqlParams, 0)
	deleted := make(map[string]struct{})
	renamed := make(map[string]struct{})
	readyA := make(map[string]struct{})

	for _, a := range paramsA {
		if a.PrevColName != "" {
			renamed[strings.Title(a.PrevColName)] = struct{}{}
		}
	}

	for nameA, a := range paramsA {
		if _, ok := paramsB[nameA]; !ok && a.PrevColName == "" {
			differenceCreate = append(differenceCreate, a)
			continue
		}
		for nameB, b := range paramsB {
			_, isRenamed := renamed[nameB]
			isChanged := nameA == nameB && (a.PrevColName != "" || !a.Equals(*b))
			_, isDoneBefore := readyA[nameA]
			if !isDoneBefore && (isRenamed || isChanged) {
				differenceUpdate = append(differenceUpdate, a)
				readyA[nameA] = struct{}{}
				continue
			}
			if _, ok := paramsA[nameB]; !ok && !isRenamed {
				if _, ok := deleted[nameB]; !ok {
					differenceDelete = append(differenceDelete, b)
					deleted[nameB] = struct{}{}
				}
				continue
			}
		}
	}

	return ModelDifference{
		Create:        differenceCreate,
		Update:        differenceUpdate,
		Delete:        differenceDelete,
		SharedChanged: m.Shared != _m.Shared,
	}
}

// ColsWithRefs returns parameters of columns that references to some struct
func (m *Model) ColsWithRefs() []*PsqlParams {
	ps := make([]*PsqlParams, 0)
	for _, p := range m.Psql {
		if p.IsStruct {
			ps = append(ps, &p)
		}
	}
	return ps
}

// SameOrPrevious - if input "p" has PrevColName, method returns previous column from current model, overwise it returns same column from current model
func (m *Model) SameOrPrevious(p *PsqlParams) *PsqlParams {
	if p.PrevColName != "" {
		return m.PsqlMap[strings.Title(p.PrevColName)]
	}
	return m.PsqlMap[p.Name]
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

	Model1      *Model
	Model1Col   *PsqlParams
	RefTableOne string
	RefIDOne    string
	FieldIDOne  string
	TypeIDOne   string

	Model2      *Model
	Model2Col   *PsqlParams
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
	Models      map[string]*Model
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

	Tags map[string]struct{}

	// CurModel - field for define needed model for template
	CurModel             string
	NewModelObj          *Model
	OldModelObj          *Model
	LastMigrationVersion int
}

// ForEachModel performs function "f" for each model, passing model name and model object in "f".
// Model nesting is also considered.
func (cfg *Config) ForEachModel(f func(name string, model *Model) error) error {

	// sort map keys to iterate through the map in order
	keys := make([]string, 0)
	for k := range cfg.Models {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i := 0; i <= cfg.MaxDeepNesting; i++ {
		for _, modelName := range keys {
			model := cfg.Models[modelName]
			if model.DeepNesting == i {
				err := f(modelName, model)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
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
