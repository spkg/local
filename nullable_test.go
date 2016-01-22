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
