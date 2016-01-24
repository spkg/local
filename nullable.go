package local

import (
	"database/sql/driver"
	"errors"
)

var errNilPtr = errors.New("destination pointer is nil")

// NullDate represents a Date that may be null.
// NullDate implements the sql Scanner interface so
// it can be used as a scan destination, similar to
// sql.NullString.
type NullDate struct {
	Date  Date
	Valid bool // Valid is true if Date is not NULL
}

// Scan implements the sql Scanner interface
func (n *NullDate) Scan(value interface{}) error {
	if n == nil {
		return errNilPtr
	}

	if value == nil {
		n.Date, n.Valid = Date{}, false
		return nil
	}

	err := n.Date.Scan(value)
	if err == nil {
		n.Valid = true
	}

	return err
}

// Value implements the driver Valuer interface.
func (n NullDate) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Date.Value()
}
