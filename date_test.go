package local

import (
	"encoding/xml"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToday(t *testing.T) {
	date := Today()
	now := time.Now()

	d1, m1, y1 := now.Date()
	d2, m2, y2 := date.Date()

	assert.Equal(t, d1, d2)
	assert.Equal(t, m1, m2)
	assert.Equal(t, y1, y2)
}

func BenchmarkYears(b *testing.B) {
	for i := 0; i < b.N; i++ {
		month := 5
		day := 20
		year := 1934

		_ = DateFor(year, time.Month(month), day)
		//CheckLocalDate(t, date, year, month, day)
	}
}

func TestYears(t *testing.T) {
	for year := -9999; year <= 9999; year++ {
		month := 5
		day := 20

		date := DateFor(year, time.Month(month), day)
		CheckLocalDate(t, date, year, month, day)
	}
}

func TestMonths(t *testing.T) {
	for month := 1; month <= 12; month++ {
		year := 1969
		day := 12

		date := DateFor(year, time.Month(month), day)
		CheckLocalDate(t, date, year, month, day)
	}
}

func TestDays(t *testing.T) {
	for day := 1; day <= 31; day++ {
		year := 1969
		month := 1

		date := DateFor(year, time.Month(month), day)
		CheckLocalDate(t, date, year, month, day)
	}
}

func CheckLocalDate(t *testing.T, date Date, year, month, day int) {
	assert := assert.New(t)
	assert.Equal(year, date.Year())
	assert.Equal(month, int(date.Month()))
	assert.Equal(day, date.Day())

	// Calculate expected text representation
	var text string
	if year < 0 {
		text = fmt.Sprintf("%05d-%02d-%02d", year, month, day)

	} else {
		text = fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	}

	assert.Equal(text, date.String())

	if date2, err := DateParse(text); err != nil || !date.Equal(date2) {
		if err != nil {
			t.Errorf("DateParse: %s: unexpected error: %v", text, err)
		} else {
			t.Errorf("DateParse: expected=%v, actual=%v", date, date2)
		}
	}

	// for non-negative years, can check parsing with time package
	if year >= 0 {
		if tm, err := time.Parse("2006-01-02", text); err != nil {
			t.Errorf("time.Parse: unexpected error parsing %s: %v", text, err)
		} else {
			y := tm.Year()
			m := int(tm.Month())
			d := tm.Day()
			if y != year {
				t.Errorf("time.Parse: Year: expected %d, actual %d", year, y)
			}
			if m != month {
				t.Errorf("time.Parse: Month: expected %d, actual %d", month, m)
			}
			if d != day {
				t.Errorf("time.Parse: Day: expected %d, actual %d", day, d)
			}
		}
	}

	// check marshalling and unmarshalling JSON
	data, err := date.MarshalJSON()
	if err != nil {
		t.Errorf("MarshalJSON: %s: unexpected error: %v", text, err)
	} else {
		assert.Equal(`"`+text+`"`, string(data))
		var date2 Date
		err = date2.UnmarshalJSON(data)
		if err != nil {
			t.Errorf("UnmarshalJSON: %s: unexpected error: %v", text, err)
		} else {
			if !date.Equal(date2) {
				t.Errorf("UnmarshalJSON: expected %s, actual %s", date.String(), date2.String())
			}
		}
	}

	// check marshalling and unmarshalling text
	data, err = date.MarshalText()
	if err != nil {
		t.Errorf("MarshalText: %s: unexpected error: %v", text, err)
	} else {
		assert.Equal(text, string(data))
		var date2 Date
		err = date2.UnmarshalText(data)
		if err != nil {
			t.Errorf("UnmarshalText: %s: unexpected error: %v", text, err)
		} else {
			if !date.Equal(date2) {
				t.Errorf("UnmarshalText: expected %s, actual %s", date.String(), date2.String())
			}
		}
	}

	// marshal and unmarshal binary
	data, err = date.MarshalBinary()
	if err != nil {
		t.Errorf("MarshalBinary: %s: unexpected error: %v", text, err)
	} else {
		// binary should be the same as the equivalent time binary
		tdata, _ := date.t.MarshalBinary()
		assert.Equal(tdata, data)
		var date2 Date
		err = date2.UnmarshalBinary(data)
		assert.NoError(err, date.String())
	}
}

func TestParseDate(t *testing.T) {
	testCases := []struct {
		Text  string
		Valid bool
		Day   int
		Month time.Month
		Year  int
	}{
		{
			Text:  "2095-09-30",
			Valid: true,
			Day:   30,
			Month: time.September,
			Year:  2095,
		},
		{
			Text:  "2195-060",
			Valid: true,
			Day:   1,
			Month: time.March,
			Year:  2195,
		},
		{
			Text:  "2095.09.30",
			Valid: true,
			Day:   30,
			Month: time.September,
			Year:  2095,
		},
		{
			Text:  "2095/09/30",
			Valid: true,
			Day:   30,
			Month: time.September,
			Year:  2095,
		},
		{
			Text:  "20951030",
			Valid: true,
			Day:   30,
			Month: time.October,
			Year:  2095,
		},
		{
			Text:  "2195-060",
			Valid: true,
			Day:   1,
			Month: time.March,
			Year:  2195,
		},
		{
			Text:  "2195074",
			Valid: true,
			Day:   15,
			Month: time.March,
			Year:  2195,
		},
	}
	assert := assert.New(t)

	for _, tc := range testCases {
		for _, suffix := range []string{"", "T00:00:00Z", "T00:00:00", "T00:00:00+10:000", "T000000+0900"} {
			baseText := tc.Text + suffix
			for _, text := range []string{baseText, " \t" + baseText + "\t\t "} {
				ld, err := DateParse(text)
				if tc.Valid {
					assert.NoError(err, text)
					assert.Equal(tc.Day, ld.Day())
					assert.Equal(tc.Month, ld.Month())
					assert.Equal(tc.Year, ld.Year())
				} else {
					assert.Error(err, text)
				}
			}
		}
	}
}

func TestMarshalXML(t *testing.T) {
	assert := assert.New(t)
	type testStruct struct {
		XMLName        xml.Name `xml:"TestCase"`
		Element        Date
		AnotherElement Date `xml:"another"`
		Attribute1     Date `xml:",attr"`
		Attribute2     Date `xml:"attribute-2,attr"`
	}

	testCases := []struct {
		st  testStruct
		xml string
	}{
		{
			st: testStruct{
				Element:        mustParseDate("2021-01-02"),
				AnotherElement: mustParseDate("2021-01-03"),
				Attribute1:     mustParseDate("2021-01-04"),
				Attribute2:     mustParseDate("2021-01-05"),
			},
			xml: `<TestCase Attribute1="2021-01-04" attribute-2="2021-01-05"><Element>2021-01-02</Element><another>2021-01-03</another></TestCase>`,
		},
	}

	for _, tc := range testCases {
		b, err := xml.Marshal(&tc.st)
		assert.NoError(err)
		assert.Equal(tc.xml, string(b))
		var st testStruct
		err = xml.Unmarshal([]byte(tc.xml), &st)
		assert.NoError(err)
		assert.Equal("", st.XMLName.Space)
		assert.Equal("TestCase", st.XMLName.Local)
		st.XMLName.Local = ""
		assert.Equal(tc.st, st)
	}
}

func TestDateAfter(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		Date1, Date2 Date
	}{
		{DateFor(1999, 9, 30), DateFor(1999, 10, 1)},
		{DateFor(0, 9, 30), DateFor(0, 10, 1)},
	}

	for _, tc := range testCases {
		assert.True(tc.Date1.Before(tc.Date2))
		assert.True(tc.Date2.After(tc.Date1))
		assert.False(tc.Date2.Before(tc.Date1))
		assert.False(tc.Date1.After(tc.Date2))
	}
}

func TestDateWeekday(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		Date    Date
		Weekday time.Weekday
	}{
		{DateFor(1999, 9, 30), time.Thursday},
		{DateFor(1997, 1, 30), time.Thursday},
		{DateFor(1994, 11, 14), time.Monday},
		{DateFor(1992, 12, 16), time.Wednesday},
		{DateFor(2033, 1, 4), time.Tuesday},
		{DateFor(2033, 4, 8), time.Friday},
		{DateFor(2033, 4, 9), time.Saturday},
		{DateFor(2042, 7, 6), time.Sunday},
	}

	for _, tc := range testCases {
		assert.Equal(tc.Weekday, tc.Date.Weekday(), tc.Date.String())
	}
}

func TestScan(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		Value    interface{}
		Error    bool
		Expected Date
	}{
		{
			Value:    "2056-11-13",
			Expected: DateFor(2056, 11, 13),
		},
		{
			Value:    time.Date(2056, 10, 31, 0, 0, 0, 0, time.UTC),
			Expected: DateFor(2056, 10, 31),
		},
		{
			Value:    time.Date(2056, 9, 30, 1, 2, 3, 400000, time.FixedZone("Australia/Brisbane", 10*3600)),
			Expected: DateFor(2056, 9, 30),
		},
		{Value: []byte("2157-12-31"), Expected: DateFor(2157, 12, 31)},
		{Value: []byte("xxx"), Error: true},
		{Value: nil, Expected: Date{}},
		{Value: int64(11), Error: true},
		{Value: true, Error: true},
		{Value: float64(11.1), Error: true},
	}

	for _, tc := range testCases {
		var d Date
		err := d.Scan(tc.Value)
		if tc.Error {
			assert.Error(err)
		} else {
			assert.NoError(err)
			assert.True(d.Equal(tc.Expected))
		}
	}
}

func TestValue(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		Date     Date
		Expected time.Time
	}{
		{DateFor(2071, 1, 30), time.Date(2071, 1, 30, 0, 0, 0, 0, time.UTC)},
	}

	for _, tc := range testCases {
		v, err := tc.Date.Value()
		assert.NoError(err)
		assert.IsType(time.Time{}, v)
		t := v.(time.Time)
		assert.True(t.Equal(tc.Expected))
	}
}

func mustParseDate(s string) Date {
	d, err := DateParse(s)
	if err != nil {
		panic(err.Error())
	}
	return d
}

func TestDateProperties(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		Date      Date
		IsZero    bool
		Unix      int64
		Year      int
		Week      int
		YearDay   int
		Formatted string
	}{
		{
			Date:      Date{},
			IsZero:    true,
			Unix:      -62135596800,
			Year:      1,
			Week:      1,
			YearDay:   1,
			Formatted: "1 Jan 0001",
		},
		{
			Date:      DateFor(1970, 1, 1),
			IsZero:    false,
			Unix:      0,
			Year:      1970,
			Week:      1,
			YearDay:   1,
			Formatted: "1 Jan 1970",
		},
		{
			Date:      DateFor(2048, 1, 30),
			IsZero:    false,
			Unix:      2463955200,
			Year:      2048,
			Week:      5,
			YearDay:   30,
			Formatted: "30 Jan 2048",
		},
	}

	for _, tc := range testCases {
		assert.Equal(tc.IsZero, tc.Date.IsZero())
		assert.Equal(tc.Unix, tc.Date.Unix())
		year, week := tc.Date.ISOWeek()
		assert.Equal(tc.Year, year, "Year")
		assert.Equal(tc.Week, week, "Week")
		assert.Equal(tc.YearDay, tc.Date.YearDay())
		assert.Equal(tc.Formatted, tc.Date.Format("2 Jan 2006"))
	}
}

// Test for case where attempt made to unmarshal invalid binary data
func TestDateUnmarshalBinaryError(t *testing.T) {
	assert := assert.New(t)
	data := []byte("xxxx")
	var d Date
	err := d.UnmarshalBinary(data)
	assert.Error(err)
}

func TestDateAddDate(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		Date     Date
		Years    int
		Months   int
		Days     int
		Expected Date
	}{
		{
			Date:     DateFor(2029, 12, 16),
			Years:    1,
			Expected: DateFor(2030, 12, 16),
		},
		{
			Date:     DateFor(2029, 12, 16),
			Years:    1,
			Months:   3,
			Expected: DateFor(2031, 3, 16),
		},
		{
			Date:     DateFor(2029, 12, 16),
			Years:    1,
			Months:   3,
			Days:     30,
			Expected: DateFor(2031, 4, 15),
		},
		{
			Date:     DateFor(2029, 12, 16),
			Years:    -1,
			Months:   3,
			Days:     30,
			Expected: DateFor(2029, 4, 15),
		},
		{
			Date:     DateFor(2029, 12, 16),
			Months:   -13,
			Expected: DateFor(2028, 11, 16),
		},
		{
			Date:     DateFor(2029, 12, 16),
			Days:     15,
			Expected: DateFor(2029, 12, 31),
		},
	}
	for _, tc := range testCases {
		d := tc.Date.AddDate(tc.Years, tc.Months, tc.Days)
		assert.Equal(tc.Expected, d, tc.Expected.String()+" vs "+d.String())
	}
}

func TestDateSub(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		Date1 Date
		Date2 Date
		Days  int
	}{
		{
			Date1: DateFor(1994, 11, 14),
			Date2: DateFor(1994, 11, 13),
			Days:  1,
		},
		{
			Date1: DateFor(1994, 11, 14),
			Date2: DateFor(1994, 11, 15),
			Days:  -1,
		},
		{
			Date1: DateFor(1994, 11, 14),
			Date2: DateFor(1992, 12, 16),
			Days:  698,
		},
	}
	for _, tc := range testCases {
		d := tc.Date1.Sub(tc.Date2)
		assert.Equal(time.Duration(tc.Days)*time.Hour*24, d)
	}
}
