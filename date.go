package local

import (
	"encoding/xml"
	"fmt"
	"time"
)

// Date represents a date without a time or a timezone.
// Useful for representing date of birth, for example.
//
// Calculations on Date are performed using the standard
// library's time.Time type. For these calculations the time is
// midnight and the timezone is UTC.
type Date struct {
	t time.Time
}

// After reports whether the local date d is after e.
func (d Date) After(e Date) bool {
	return d.t.After(e.t)
}

// Before reports whether the local date d is before e.
func (d Date) Before(e Date) bool {
	return d.t.Before(e.t)
}

// Equal reports whether d and e represent the same local date.
func (d Date) Equal(e Date) bool {
	return d.t.Equal(e.t)
}

// IsZero reports whether d represents the zero local date,
// January 1, year 1.
func (d Date) IsZero() bool {
	return d.t.IsZero()
}

// Date returns the year, month and day on which d occurs.
func (d Date) Date() (year int, month time.Month, day int) {
	return d.t.Date()
}

// Unix returns d as a Unix time, the number of seconds elapsed
// since January 1, 1970 UTC to midnight of the date UTC.
func (d Date) Unix() int64 {
	return d.t.Unix()
}

// Year returns the year in which d occurs.
func (d Date) Year() int {
	return d.t.Year()
}

// Month returns the month of the year specified by d.
func (d Date) Month() time.Month {
	return d.t.Month()
}

// Day returns the day of the month specified by d.
func (d Date) Day() int {
	return d.t.Day()
}

// Weekday returns the day of the week specified by d.
func (d Date) Weekday() time.Weekday {
	return d.t.Weekday()
}

// ISOWeek returns the ISO 8601 year and week number in which d occurs.
// Week ranges from 1 to 53. Jan 01 to Jan 03 of year n might belong to
// week 52 or 53 of year n-1, and Dec 29 to Dec 31 might belong to week 1
// of year n+1.
func (d Date) ISOWeek() (year, week int) {
	year, week = d.t.ISOWeek()
	return
}

// YearDay returns the day of the year specified by D, in the range [1,365] for non-leap years,
// and [1,366] in leap years.
func (d Date) YearDay() int {
	return d.t.YearDay()
}

// Add returns the local date d + duration.
func (d Date) Add(duration time.Duration) Date {
	t := d.t.Add(toDays(duration))
	return Date{t: t}
}

// Sub returns the duration d-e, which will be an integral number of days.
// If the result exceeds the maximum (or minimum) value that can be stored
// in a Duration, the maximum (or minimum) duration will be returned.
// To compute d-duration, use d.Add(-duration).
func (d Date) Sub(e Date) time.Duration {
	return d.t.Sub(e.t)
}

// AddDate returns the local date corresponding to adding the given number of years,
// months, and days to t. For example, AddDate(-1, 2, 3) applied to January 1, 2011
// returns March 4, 2010.
//
// AddDate normalizes its result in the same way that Date does, so, for example,
// adding one month to October 31 yields December 1, the normalized form for November 31.
func (d Date) AddDate(years int, months int, days int) Date {
	t := d.t.AddDate(years, months, days)
	return Date{t: t}
}

// toDate converts the time.Time value into a Date.,
func toLocalDate(t time.Time) Date {
	y, m, d := t.Date()
	return DateFor(y, m, d)
}

// Today returns the current local date.
func Today() Date {
	return toLocalDate(time.Now())
}

// DateFor returns the Date corresponding to year, month and date.
//
// The month and day values may be outside their usual ranges
// and will be normalized during the conversion.
// For example, October 32 converts to November 1.
func DateFor(year int, month time.Month, day int) Date {
	return Date{
		t: time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
	}
}

// String returns a string representation of d. The date
// format returned is compatible with ISO 8601: yyyy-mm-dd.
func (d Date) String() string {
	return toDateString(d)
}

// toDateString returns the string representation of the date.
func toDateString(d Date) string {
	year, month, day := d.Date()
	sign := ""
	if year < 0 {
		year = -year
		sign = "-"
	}
	return fmt.Sprintf("%s%04d-%02d-%02d", sign, year, int(month), day)
}

// toQuotedDateString returns the string representation of the date in quotation marks.
func toQuotedDateString(d Date) string {
	return fmt.Sprintf(`"%s"`, toDateString(d))
}

// MarshalJSON implements the json.Marshaler interface.
// The date is a quoted string in an ISO 8601 format (yyyy-mm-dd).
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(toQuotedDateString(d)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The date is expected to be a quoted string in an ISO 8601
// format (calendar or ordinal).
func (d *Date) UnmarshalJSON(data []byte) (err error) {
	s := string(data)
	*d, err = ParseDate(s)
	return
}

// MarshalText implements the encoding.TextMarshaller interface.
// The date format is yyyy-mm-dd.
func (d Date) MarshalText() ([]byte, error) {
	return []byte(toDateString(d)), nil
}

// UnmarshalText implements the encoding.TextUnmarshaller interface.
// The date is expected to an ISO 8601 format (calendar or ordinal).
func (d *Date) UnmarshalText(data []byte) (err error) {
	s := string(data)
	*d, err = ParseDate(s)
	return
}

func (d *Date) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeElement(toDateString(*d), start)
	return nil
}

func (d *Date) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var s string

	if err := decoder.DecodeElement(&s, &start); err != nil {
		return err
	}

	if ld, err := ParseDate(s); err != nil {
		return err
	} else {
		*d = ld
	}
	return nil
}

func (d *Date) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name:  name,
		Value: d.String(),
	}, nil
}

func (d *Date) UnmarshalXMLAttr(attr xml.Attr) error {
	if ld, err := ParseDate(attr.Value); err != nil {
		return err
	} else {
		*d = ld
	}
	return nil
}
