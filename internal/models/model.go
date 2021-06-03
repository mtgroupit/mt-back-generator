package models

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

	Psql           []PsqlParams
	SQLSelectStr   string
	SQLWhereParams string
	SQLAddStr      string
	SQLEditStr     string

	Binds        []Bind
	Methods      []string
	MethodsProps []MethodProps

	RulesSet map[string][]string `yaml:"rules-set"`
}

// RequiredColumns -  returns slice of required columns
func (m Model) RequiredColumns() []string {
	var reqColumns []string
	for column, options := range m.Columns {
		if options.Required {
			reqColumns = append(reqColumns, column)
		}
	}
	return reqColumns
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
	NotNull  bool
	FK       string
	Last     bool
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
	FilteredFields          []string
	EditableFields          []string
	NestedObjs              []NestedObjProps

	Rules []string
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
