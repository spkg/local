package local

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	errInvalidDateFormat     = errors.New("invalid date format")
	errInvalidDateTimeFormat = errors.New("invalid date-time format")
)

var parseFormats = struct {
	calendarDates  []string
	ordinalDates   []string
	times          []string
	throwAwayTimes []string
}{
	calendarDates: []string{
		`(-?\d{4})-(\d{1,2})-(\d{1,2})`,
		`^(-?\d{4})(\d{2})(\d{2})`,
		// Not ISO 8601, but still unambiguous
		`(-?\d{4})\.(\d{1,2})\.(\d{1,2})`,
		`(-?\d{4})/(\d{1,2})/(\d{1,2})`,
	},
	ordinalDates: []string{
		`(-?\d{4})-(\d{3})`,
		`(-?\d{4})(\d{3})`,
	},
	times: []string{
		`(\d{1,2}):(\d{1,2}):(\d{1,2})(\.\d*)?`,
		`(\d{1,2}):(\d{1,2})`,
		`(\d{2})(\d{2})(\d{2})(\.\d*)?`,
		`(\d{2})(\d{2})`,
	},
	throwAwayTimes: []string{
		`(T[0-9:.zZ+-]*)?`,
	},
}

var parseRegexp = struct {
	calendarDates     []*regexp.Regexp
	ordinalDates      []*regexp.Regexp
	calendarDateTimes []*regexp.Regexp
	ordinalDateTimes  []*regexp.Regexp
}{}

const (
	startRE = `^\s*`
	endRE   = `\s*$`
)

func init() {
	for _, cd := range parseFormats.calendarDates {
		for _, tat := range parseFormats.throwAwayTimes {
			text := startRE + cd + tat + endRE
			parseRegexp.calendarDates = append(parseRegexp.calendarDates, regexp.MustCompile(text))
		}

		text := startRE + cd + endRE
		parseRegexp.calendarDateTimes = append(parseRegexp.calendarDateTimes, regexp.MustCompile(text))

		for _, tod := range parseFormats.times {
			text = startRE + cd + "T" + tod + endRE
			parseRegexp.calendarDateTimes = append(parseRegexp.calendarDateTimes, regexp.MustCompile(text))
			text = startRE + cd + `\s+` + tod + endRE
			parseRegexp.calendarDateTimes = append(parseRegexp.calendarDateTimes, regexp.MustCompile(text))
		}
	}

	for _, od := range parseFormats.ordinalDates {
		for _, tat := range parseFormats.throwAwayTimes {
			text := startRE + od + tat + endRE
			parseRegexp.ordinalDates = append(parseRegexp.ordinalDates, regexp.MustCompile(text))
		}

		text := startRE + od + endRE
		parseRegexp.ordinalDateTimes = append(parseRegexp.ordinalDateTimes, regexp.MustCompile(text))

		for _, tod := range parseFormats.times {
			text = startRE + od + "T" + tod + endRE
			parseRegexp.ordinalDateTimes = append(parseRegexp.ordinalDateTimes, regexp.MustCompile(text))
		}
	}
}

// DateParseLayout parses a formatted string and returns the date value it represents.
// The layout is based on the standard library time package and for local dates the reference is
//  Mon Jan 2 2006
// If the layout contains time or timezone fields, they are parsed and discarded.
func DateParseLayout(layout, value string) (Date, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return Date{}, err
	}
	return DateFromTime(t), nil
}

// DateParse attempts to parse a string into a local date. Leading
// and trailing space and quotation marks are ignored. The following
// date formates are recognized: yyyy-mm-dd, yyyymmdd, yyyy.mm.dd,
// yyyy/mm/dd, yyyy-ddd, yyyyddd.
//
// DateParse is used to parse dates where no layout is provided, for example
// when marshaling and unmarshaling JSON and XML.
func DateParse(s string) (Date, error) {
	s = strings.Trim(s, " \t\"'")
	for _, regexp := range parseRegexp.calendarDates {
		match := regexp.FindStringSubmatch(s)
		if match != nil {
			// no error checking here because matching the regexp
			// guarantees that parsing the strings will succeed.
			year, _ := strconv.ParseInt(match[1], 10, 0)
			month, _ := strconv.ParseInt(match[2], 10, 0)
			day, _ := strconv.ParseInt(match[3], 10, 0)
			return DateFor(int(year), time.Month(month), int(day)), nil
		}
	}

	for _, regexp := range parseRegexp.ordinalDates {
		match := regexp.FindStringSubmatch(s)
		if match != nil {
			// no error checking here because matching the regexp
			// guarantees that parsing the strings will succeed.
			year, _ := strconv.ParseInt(match[1], 10, 0)
			dayOfYear, _ := strconv.ParseInt(match[2], 10, 0)
			duration := time.Duration((dayOfYear - 1) * nanosecondsPerDay)
			return DateFor(int(year), 1, 1).Add(duration), nil
		}
	}

	return Date{}, errInvalidDateFormat
}

// DateTimeParseLayout parses a formatted string and returns the date value it represents.
// The layout is based on the standard library time package and for local date-times the reference is
//  Mon Jan 2 2006 15:04:05
// If the layout contains a timezone field, it is parsed and discarded.
func DateTimeParseLayout(layout, value string) (DateTime, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return DateTime{}, err
	}
	return DateTimeFromTime(t), nil
}

// DateTimeParse attempts to parse a string into a local date-time. Leading
// and trailing space and quotation marks are ignored. The following
// date formates are recognized: yyyy-mm-dd, yyyymmdd, yyyy.mm.dd,
// yyyy/mm/dd, yyyy-ddd, yyyyddd. The following time formats are recognized:
// HH:MM:SS, HH:MM, HHMMSS, HHMM.
func DateTimeParse(s string) (DateTime, error) {
	s = strings.Trim(s, " \t\"'")
	for _, regexp := range parseRegexp.calendarDateTimes {
		match := regexp.FindStringSubmatch(s)
		if match != nil {
			// no error checking here because matching the regexp
			// guarantees that parsing the strings will succeed.
			year, _ := strconv.ParseInt(match[1], 10, 0)
			month, _ := strconv.ParseInt(match[2], 10, 0)
			day, _ := strconv.ParseInt(match[3], 10, 0)

			var hour, minute, second int64
			if len(match) > 4 {
				hour, _ = strconv.ParseInt(match[4], 10, 0)
			}
			if len(match) > 5 {
				minute, _ = strconv.ParseInt(match[5], 10, 0)
			}
			if len(match) > 6 {
				second, _ = strconv.ParseInt(match[6], 10, 0)
			}

			return DateTimeFor(int(year), time.Month(month), int(day), int(hour), int(minute), int(second)), nil
		}
	}

	for _, regexp := range parseRegexp.ordinalDateTimes {
		match := regexp.FindStringSubmatch(s)
		if match != nil {
			// no error checking here because matching the regexp
			// guarantees that parsing the strings will succeed.
			year, _ := strconv.ParseInt(match[1], 10, 0)
			dayOfYear, _ := strconv.ParseInt(match[2], 10, 0)

			var hour, minute, second int64
			if len(match) > 3 {
				hour, _ = strconv.ParseInt(match[3], 10, 0)
			}
			if len(match) > 4 {
				minute, _ = strconv.ParseInt(match[4], 10, 0)
			}
			if len(match) > 5 {
				second, _ = strconv.ParseInt(match[5], 10, 0)
			}

			duration := time.Duration((dayOfYear - 1) * nanosecondsPerDay)
			return DateTimeFor(int(year), 1, 1, int(hour), int(minute), int(second)).Add(duration), nil
		}
	}

	return DateTime{}, errInvalidDateTimeFormat
}
