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

	assert.Equal(t, year, date.Year())
	assert.Equal(t, month, int(date.Month()))
	assert.Equal(t, day, date.Day())

	// Calculate expected text representation
	var text string
	if year < 0 {
		text = fmt.Sprintf("%05d-%02d-%02d", year, month, day)

	} else {
		text = fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	}

	assert.Equal(t, text, date.String())

	if date2, err := ParseDate(text); err != nil || !date.Equal(date2) {
		if err != nil {
			t.Errorf("ParseDate: %s: unexpected error: %v", text, err)
		} else {
			t.Errorf("ParseDate: expected=%v, actual=%v", date, date2)
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
		assert.Equal(t, `"`+text+`"`, string(data))
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
		assert.Equal(t, text, string(data))
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
			text := tc.Text + suffix
			ld, err := ParseDate(text)
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
				Element:        MustParseDate("2021-01-02"),
				AnotherElement: MustParseDate("2021-01-03"),
				Attribute1:     MustParseDate("2021-01-04"),
				Attribute2:     MustParseDate("2021-01-05"),
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
