package extauthapi

import (
	"database/sql/driver"
	"errors"

	uuid "github.com/satori/go.uuid"
)

// ID contains user ID (UUIDv4) or NoID.
// On Scan/Unmarshal values nil and empty string are considered NoID.
// On String/Value/Marshal NoID is 00000000-0000-0000-0000-000000000000.
type ID uuid.UUID

// NoID means absent user ID.
var NoID = ID(uuid.Nil) //nolint:gochecknoglobals

var (
	errWrongUUIDVersion = errors.New("wrong UUID version")
)

// NewID returns new user ID.
func NewID() ID {
	id := uuid.NewV4()
	return ID(id)
}

// ParseID returns ID parsed from string input.
// Input is expected in a form accepted by UnmarshalText.
func ParseID(input string) (id ID, err error) {
	err = id.UnmarshalText([]byte(input))
	return
}

// MustParseID returns ID parsed from string input.
// Same behavior as ParseID, but panics on error.
func MustParseID(input string) ID {
	id, err := ParseID(input)
	if err != nil {
		panic(err)
	}
	return id
}

// Returns canonical string representation of UUID:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (id ID) String() string {
	return uuid.UUID(id).String()
}

// Value implements the driver.Valuer interface.
func (id ID) Value() (driver.Value, error) {
	return uuid.UUID(id).Value()
}

// Scan implements the sql.Scanner interface.
func (id *ID) Scan(src interface{}) error {
	if src == nil {
		*id = NoID
		return nil
	}
	err := (*uuid.UUID)(id).Scan(src)
	if err == nil && *id != NoID && (*uuid.UUID)(id).Version() != uuid.V4 {
		return errWrongUUIDVersion
	}
	return err
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (id *ID) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*id = NoID
		return nil
	}
	err := (*uuid.UUID)(id).UnmarshalText(text)
	if err == nil && *id != NoID && (*uuid.UUID)(id).Version() != uuid.V4 {
		return errWrongUUIDVersion
	}
	return err
}
