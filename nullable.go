package local

import (
	"bytes"
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

var (
	nullText = []byte("null")
)

// MarshalJSON implements the json.Marshaler interface.
func (n NullDate) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Date.MarshalJSON()
	}
	return nullText, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *NullDate) UnmarshalJSON(p []byte) error {
	if bytes.Equal(p, nullText) {
		n.Valid = false
		n.Date = Date{}
		return nil
	}
	var d Date
	if err := d.UnmarshalJSON(p); err != nil {
		return err
	}
	n.Valid = true
	n.Date = d
	return nil
}

// NullDateTime represents a DateTime that may be null.
// NullDateTime implements the sql Scanner interface so
// it can be used as a scan destination, similar to
// sql.NullString.
type NullDateTime struct {
	DateTime DateTime
	Valid    bool // Valid is true if Date is not NULL
}

// Scan implements the sql Scanner interface
func (n *NullDateTime) Scan(value interface{}) error {
	if n == nil {
		return errNilPtr
	}

	if value == nil {
		n.DateTime, n.Valid = DateTime{}, false
		return nil
	}

	err := n.DateTime.Scan(value)
	if err == nil {
		n.Valid = true
	}

	return err
}

// Value implements the driver Valuer interface.
func (n NullDateTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.DateTime.Value()
}

// MarshalJSON implements the json.Marshaler interface.
func (n NullDateTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.DateTime.MarshalJSON()
	}
	return nullText, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *NullDateTime) UnmarshalJSON(p []byte) error {
	if bytes.Equal(p, nullText) {
		n.Valid = false
		n.DateTime = DateTime{}
		return nil
	}
	var d DateTime
	if err := d.UnmarshalJSON(p); err != nil {
		return err
	}
	n.Valid = true
	n.DateTime = d
	return nil
}
