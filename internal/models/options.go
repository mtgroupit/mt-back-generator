package models

// Options contain properties of column
type Options struct {
	TitleName        string
	Type             string
	BusinessType     string
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

// AppType - returns type which using in generated service
func (o Options) AppType() string {
	switch o.BusinessType {
	case "date":
		return "*time.Time"
	case "date-time":
		return "*time.Time"
	default:
		return o.BusinessType
	}
}
