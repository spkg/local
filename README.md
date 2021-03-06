# Local date and time

[![GoDoc](https://godoc.org/github.com/spkg/local?status.svg)](https://godoc.org/github.com/spkg/local)
[![Build Status (Linux)](https://travis-ci.org/spkg/local.svg?branch=master)](https://travis-ci.org/spkg/local)
[![Build status (Windows)](https://ci.appveyor.com/api/projects/status/1l3spdhwwftuk6nt?svg=true)](https://ci.appveyor.com/project/jjeffery/local)
[![Coverage Status](https://coveralls.io/repos/github/spkg/local/badge.svg?branch=master)](https://coveralls.io/github/spkg/local?branch=master)
[![GoReportCard](https://goreportcard.com/badge/github.com/spkg/local)](https://goreportcard.com/report/github.com/spkg/local)
[![License](https://img.shields.io/badge/license-BSD-green.svg)](https://raw.githubusercontent.com/spkg/local/master/LICENSE.md)

Package local provides types for representing local dates, and times.

Sometimes it is useful to be able to represent a date or time without reference
to an instance in time with a timezone.

For example, when recording a person's date of birth all that is needed is a date.
There is no requirement to specify an instant in time with timezone.

There are also circumstances where an event will be scheduled for a date and time
in the local timezone, whatever that may be. And example of this might be a schedule for
taking medication.

Like the standard library time package, the local package uses the
[proleptic Gregorian calendar](https://en.wikipedia.org/wiki/Proleptic_Gregorian_calendar)
for all calculations. The local package makes use of the time package for all
of its date-time calculations. Because some of this code is based on the standard time package,
it has the identical license to the Go project.

For usage examples, refer to the [GoDoc](https://godoc.org/github.com/spkg/local) documentation.
