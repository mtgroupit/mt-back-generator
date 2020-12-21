package extauthapi

import (
	"database/sql/driver"

	uuid "github.com/satori/go.uuid"
)

// UserID contains user ID (UUIDv4) or NoUserID.
// On Scan/Unmarshal values nil and empty string are considered NoUserID.
// On String/Value/Marshal NoUserID is 00000000-0000-0000-0000-000000000000.
type UserID uuid.UUID

// NoUserID means absent user ID.
var NoUserID = UserID(uuid.Nil) //nolint:gochecknoglobals

// NewUserID returns new user ID.
func NewUserID() UserID {
	id := uuid.NewV4()
	return UserID(id)
}

// ParseUserID returns UserID parsed from string input.
// Input is expected in a form accepted by UnmarshalText.
func ParseUserID(input string) (id UserID, err error) {
	err = id.UnmarshalText([]byte(input))
	return
}

// MustParseUserID returns UserID parsed from string input.
// Same behavior as ParseUserID, but panics on error.
func MustParseUserID(input string) UserID {
	id, err := ParseUserID(input)
	if err != nil {
		panic(err)
	}
	return id
}

// Returns canonical string representation of UUID:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (id UserID) String() string {
	return uuid.UUID(id).String()
}

// Value implements the driver.Valuer interface.
func (id UserID) Value() (driver.Value, error) {
	return uuid.UUID(id).Value()
}

// Scan implements the sql.Scanner interface.
func (id *UserID) Scan(src interface{}) error {
	if src == nil {
		*id = NoUserID
		return nil
	}
	err := (*uuid.UUID)(id).Scan(src)
	if err == nil && *id != NoUserID && (*uuid.UUID)(id).Version() != uuid.V4 {
		return errWrongUUIDVersion
	}
	return err
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (id *UserID) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*id = NoUserID
		return nil
	}
	err := (*uuid.UUID)(id).UnmarshalText(text)
	if err == nil && *id != NoUserID && (*uuid.UUID)(id).Version() != uuid.V4 {
		return errWrongUUIDVersion
	}
	return err
}
