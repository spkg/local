// Package local provides types for representing local dates, times and date-times.
// Sometimes it is useful to be able to represent a date or time without reference
// to an instance in time with a timezone.
//
// For example, when recording a person's date of birth all that is needed is a date.
// There is no requirement to specify an instant in time with timezone.
//
// Similarly, when scheduling activities to happen such as a meal time, wakeup,
// bedtime, etc, it is enough to be able to specify the time of day. The date and
// timezone might not be relevant in the context.
//
// Like the standard library time package, the local package uses a Gregorian calendar
// for all calculations. The local package makes use of the time package for all
// of its date-time calculations.
package local
