# Local date and time

[![GoDoc](https://godoc.org/github.com/spkg/local?status.svg)](https://godoc.org/github.com/spkg/local)
[![Build Status](https://travis-ci.org/spkg/local.svg?branch=master)](https://travis-ci.org/spkg/local)
[![Coverage](http://gocover.io/_badge/github.com/spkg/local)](http://gocover.io/github.com/spkg/local)
[![License](https://img.shields.io/badge/license-BSD-green.svg)](https://raw.githubusercontent.com/spkg/local/master/LICENSE.md)

Package local provides types for representing local dates, times and date-times.
Sometimes it is useful to be able to represent a date or time without reference
to an instance in time with a timezone.

For example, when recording a person's date of birth all that is needed is a date.
There is no requirement to specify an instant in time with timezone.

Similarly, when scheduling activities to happen such as a meal time, wakeup,
bedtime, etc, it is enough to be able to specify the time of day. The date and
timezone might not be relevant in the context.

Like the standard library time package, the local package uses a Gregorian calendar
for all calculations. The local package makes use of the time package for all
of its date-time calculations. Because this code is based on the standard time package,
it has the identical license to the Go project.

For usage examples, refer to the [GoDoc](https://godoc.org/github.com/spkg/local) documentation.
