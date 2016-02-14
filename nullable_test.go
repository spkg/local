package local

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNullDateScan(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		Input         interface{}
		ExpectedError string
		ExpectedDate  NullDate
	}{
		{
			Input:        time.Date(2091, 11, 14, 18, 47, 0, 0, time.UTC),
			ExpectedDate: NullDate{DateFor(2091, 11, 14), true},
		},
		{
			Input:        "2091-11-14",
			ExpectedDate: NullDate{DateFor(2091, 11, 14), true},
		},
		{
			Input:        []byte("2091-11-14"),
			ExpectedDate: NullDate{DateFor(2091, 11, 14), true},
		},
		{
			Input:         "xxxx",
			ExpectedError: "invalid date format",
		},
		{
			Input:         24,
			ExpectedError: "cannot convert to local.Date",
		},
		{
			Input:        nil,
			ExpectedDate: NullDate{Valid: false},
		},
	}

	for _, tc := range testCases {
		var d NullDate
		err := d.Scan(tc.Input)
		if tc.ExpectedError != "" {
			assert.Error(err, tc.ExpectedError)
			assert.Equal(tc.ExpectedError, err.Error())
		} else {
			assert.NoError(err)
			if tc.ExpectedDate.Valid {
				assert.True(d.Valid)
				assert.True(d.Date.Equal(tc.ExpectedDate.Date))
			} else {
				assert.False(d.Valid)
			}
		}
	}

	// check that nil NullDate does not panic but returns error
	var nilND *NullDate
	assert.Error(nilND.Scan(nil))
}

func TestNullDateValue(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		NullDate     NullDate
		ExpectNil    bool
		ExpectedTime time.Time
	}{
		{
			NullDate:     NullDate{Date: DateFor(2087, 12, 16), Valid: true},
			ExpectedTime: time.Date(2087, 12, 16, 0, 0, 0, 0, time.UTC),
		},
		{
			NullDate:  NullDate{Date: DateFor(2087, 12, 16), Valid: false},
			ExpectNil: true,
		},
		{
			NullDate:  NullDate{},
			ExpectNil: true,
		},
	}

	for _, tc := range testCases {
		v, err := tc.NullDate.Value()
		assert.NoError(err)
		if tc.ExpectNil {
			assert.Nil(v)
		} else {
			t, ok := v.(time.Time)
			assert.True(ok)
			assert.True(tc.ExpectedTime.Equal(t))
		}
	}
}

func TestNullDateJSON(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		Text          string
		NullDate      NullDate
		ExpectedError string
	}{
		{
			Text:     "null",
			NullDate: NullDate{},
		},
		{
			Text:     `"2091-09-30"`,
			NullDate: NullDate{Date: mustParseDate("2091-09-30"), Valid: true},
		},
		{
			Text:          `25`,
			ExpectedError: "invalid date format",
		},
	}

	for _, tc := range testCases {
		var nd NullDate
		err := json.Unmarshal([]byte(tc.Text), &nd)
		if tc.ExpectedError != "" {
			assert.Error(err)
			assert.True(strings.Contains(err.Error(), tc.ExpectedError), err.Error())
		} else {
			assert.NoError(err)
			assert.Equal(tc.NullDate.Valid, nd.Valid, tc.NullDate.Date.String())
			assert.True(tc.NullDate.Date.Equal(nd.Date), tc.NullDate.Date.String(), nd.Date.String())

			p, err := json.Marshal(tc.NullDate)
			assert.NoError(err)
			assert.Equal(tc.Text, string(p))
		}
	}
}

func TestNullDateTimeScan(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		Input            interface{}
		ExpectedError    string
		ExpectedDateTime NullDateTime
	}{
		{
			Input:            time.Date(2091, 11, 14, 18, 47, 0, 0, time.UTC),
			ExpectedDateTime: NullDateTime{DateTimeFor(2091, 11, 14, 18, 47, 0), true},
		},
		{
			Input:            "2091-11-14T12:34:56",
			ExpectedDateTime: NullDateTime{DateTimeFor(2091, 11, 14, 12, 34, 56), true},
		},
		{
			Input:            []byte("2091-11-14T11:47"),
			ExpectedDateTime: NullDateTime{DateTimeFor(2091, 11, 14, 11, 47, 0), true},
		},
		{
			Input:         "xxxx",
			ExpectedError: "invalid date-time format",
		},
		{
			Input:         24,
			ExpectedError: "cannot convert to local.DateTime",
		},
		{
			Input:            nil,
			ExpectedDateTime: NullDateTime{Valid: false},
		},
	}

	for _, tc := range testCases {
		var d NullDateTime
		err := d.Scan(tc.Input)
		if tc.ExpectedError != "" {
			assert.Error(err, tc.ExpectedError)
			assert.Equal(tc.ExpectedError, err.Error())
		} else {
			assert.NoError(err)
			if tc.ExpectedDateTime.Valid {
				assert.True(d.Valid)
				assert.Equal(tc.ExpectedDateTime.DateTime, d.DateTime, tc.ExpectedDateTime.DateTime.String()+" vs "+d.DateTime.String())
			} else {
				assert.False(d.Valid)
			}
		}
	}

	// check that nil NullDate does not panic but returns error
	var nilND *NullDateTime
	assert.Error(nilND.Scan(nil))
}

func TestNullDateTimeValue(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		NullDateTime NullDateTime
		ExpectNil    bool
		ExpectedTime time.Time
	}{
		{
			NullDateTime: NullDateTime{DateTime: DateTimeFor(2087, 12, 16, 11, 47, 42), Valid: true},
			ExpectedTime: time.Date(2087, 12, 16, 11, 47, 42, 0, time.UTC),
		},
		{
			NullDateTime: NullDateTime{DateTime: DateTimeFor(2087, 12, 16, 11, 47, 43), Valid: false},
			ExpectNil:    true,
		},
		{
			NullDateTime: NullDateTime{},
			ExpectNil:    true,
		},
	}

	for _, tc := range testCases {
		v, err := tc.NullDateTime.Value()
		assert.NoError(err)
		if tc.ExpectNil {
			assert.Nil(v)
		} else {
			t, ok := v.(time.Time)
			assert.True(ok)
			assert.True(tc.ExpectedTime.Equal(t))
		}
	}
}

func TestNullDateTimeJSON(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		Text          string
		NullDateTime  NullDateTime
		ExpectedError string
	}{
		{
			Text:         "null",
			NullDateTime: NullDateTime{},
		},
		{
			Text:         `"2091-01-30T22:02:00"`,
			NullDateTime: NullDateTime{DateTime: mustParseDateTime("2091-01-30T22:02"), Valid: true},
		},
		{
			Text:          `25`,
			ExpectedError: "invalid date-time format",
		},
	}

	for _, tc := range testCases {
		var nd NullDateTime
		err := json.Unmarshal([]byte(tc.Text), &nd)
		if tc.ExpectedError != "" {
			assert.Error(err)
			assert.True(strings.Contains(err.Error(), tc.ExpectedError), err.Error())
		} else {
			assert.NoError(err)
			assert.Equal(tc.NullDateTime.Valid, nd.Valid, tc.NullDateTime.DateTime.String())
			assert.True(tc.NullDateTime.DateTime.Equal(nd.DateTime), tc.NullDateTime.DateTime.String())

			p, err := json.Marshal(tc.NullDateTime)
			assert.NoError(err)
			assert.Equal(tc.Text, string(p))
		}
	}
}
