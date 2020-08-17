package models

// Options contain properties for column
type Options struct {
	TitleName string
	Type      string
	GoType    string
	Format    string
	IsStruct  bool
	IsArray   bool
	Pk        string
	Unique    bool
	Length    int64
	Default   string
}

// PsqlParams - contain properties for postgres generate
type PsqlParams struct {
	Name     string
	SqlName  string
	Unique   bool
	Type     string
	TypeSql  string
	IsArray  bool
	IsStruct bool
	FK       string
	Last     bool
}

//Bind binds tables for deleting
type Bind struct {
	ModelName string
	FieldName string
	IsArray   bool
}

// NestedObjProps contains properties of nested objects for custom list method
type NestedObjProps struct {
	Name                  string
	Type                  string
	SqlSelect             string
	Path                  string
	ParentStruct          string
	NeedLazyLoading       bool
	IsArray               bool
	IsFirstForLazyLoading bool
	IsLastForLazyLoading  bool
}

// MethodProps contains all additional informations for method
type MethodProps struct {
	HTTPMethod          string
	IsCustomList        bool
	NeedLazyLoading     bool
	CustomListSqlSelect string
	NestedObjs          []NestedObjProps
}

// Model - description one compoment of models
type Model struct {
	Description     string
	Columns         map[string]Options
	Fields          map[string]string
	Psql            []PsqlParams
	DeepNesting     int
	HaveLazyLoading bool
	IDIsUUID        bool
	HaveEmail       bool
	HaveListMethod  bool
	SqlSelectStr    string
	SqlAddStr       string
	SqlEditStr      string
	SqlExecParams   string
	TitleName       string

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

// ExtraTable is table for many-to-many realations
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
	Name           string
	Description    string
	Models         map[string]Model
	Functions      map[string]Function
	HaveListMethod bool
	HaveDateTime   bool
	MaxDeepNesting int

	ExtraTables []ExtraTable

	// CurModel - field for define need model for template
	CurModel string
}
