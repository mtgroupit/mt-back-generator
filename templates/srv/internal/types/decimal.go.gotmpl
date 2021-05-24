package types

import (
	"database/sql/driver"

	"github.com/shopspring/decimal"
)

// Decimal is implementation of github.com/shopspring/decimal.
type Decimal decimal.Decimal

// NewDecimal converts value to Decimal.
func NewDecimal(value interface{}) Decimal {
	switch value.(type) {
	case string:
		return Decimal(decimal.RequireFromString(value.(string)))
	case int32:
		return Decimal(decimal.NewFromInt32(value.(int32)))
	case int64:
		return Decimal(decimal.NewFromInt(value.(int64)))
	case float32:
		return Decimal(decimal.NewFromFloat32(value.(float32)))
	case float64:
		return Decimal(decimal.NewFromFloat(value.(float64)))
	default:
		return Decimal{}
	}
}

// NewDecimalFromString converts string to Decimal.
func NewDecimalFromString(value string) (Decimal, error) {
	dec, err := decimal.NewFromString(value)
	if err != nil {
		return Decimal{}, err
	}
	return Decimal(dec), nil
}

// DecimalsToFloat64Array converts slice of Decimal to slice of float64.
func DecimalsToFloat64Array(ds []Decimal) (fs []float64) {
	for _, d := range ds {
		fs = append(fs, d.Float64())
	}
	return
}

// Float64ToDecimalsArray converts slice of float64 to slice of Decimal.
func Float64ToDecimalsArray(fs []float64) (ds []Decimal) {
	for _, f := range fs {
		ds = append(ds, NewDecimal(f))
	}
	return
}

// Float64 returns the nearest float64 value for d and a bool indicating
// whether f represents d exactly.
func (d *Decimal) Float64() (f float64) {
	f, _ = (*decimal.Decimal)(d).Float64()
	return
}

// Equal returns whether the numbers represented by d and d2 are equal.
func (d Decimal) Equal(d2 Decimal) bool {
	return (decimal.Decimal)(d).Equal((decimal.Decimal)(d2))
}

// LessThan returns true when d is less than d2.
func (d Decimal) LessThan(d2 Decimal) bool {
	return (decimal.Decimal)(d).LessThan((decimal.Decimal)(d2))
}

// LessThanOrEqual returns true when d is less than or equal to d2.
func (d Decimal) LessThanOrEqual(d2 Decimal) bool {
	return (decimal.Decimal)(d).LessThanOrEqual((decimal.Decimal)(d2))
}

// Value implements the driver.Valuer interface for database serialization.
func (d Decimal) Value() (driver.Value, error) {
	return (decimal.Decimal)(d).Value()
}

// Scan implements the sql.Scanner interface for database deserialization.
func (d *Decimal) Scan(value interface{}) error {
	return (*decimal.Decimal)(d).Scan(value)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Decimal) UnmarshalJSON(decimalBytes []byte) error {
	return (*decimal.Decimal)(d).UnmarshalJSON(decimalBytes)
}

// MarshalJSON implements the json.Marshaler interface.
func (d Decimal) MarshalJSON() ([]byte, error) {
	return (decimal.Decimal)(d).MarshalJSON()
}

// NullDecimal represents a nullable decimal with compatibility for
// scanning null values from the database.
type NullDecimal struct {
	Decimal Decimal
	Valid   bool
}

// Scan implements the sql.Scanner interface for database deserialization.
func (d *NullDecimal) Scan(value interface{}) error {
	if value == nil {
		d.Valid = false
		return nil
	}
	d.Valid = true
	return d.Decimal.Scan(value)
}
