package local

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNullDateScan(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		Input         interface{}
		ExpectedError string
		ExpectedDate  Date
	}{
		{
			Input:        time.Date(2091, 11, 14, 18, 47, 0, 0, time.UTC),
			ExpectedDate: DateFor(2091, 11, 14),
		},
		{
			Input:        "2091-11-14",
			ExpectedDate: DateFor(2091, 11, 14),
		},
		{
			Input:        []byte("2091-11-14"),
			ExpectedDate: DateFor(2091, 11, 14),
		},
		{
			Input:         "xxxx",
			ExpectedError: "invalid date format",
		},
		{
			Input:         24,
			ExpectedError: "cannot convert to local.Date",
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
			assert.True(d.Date.Equal(tc.ExpectedDate))
		}
	}
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
