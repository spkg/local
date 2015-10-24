package local

import "time"

const (
	secondsPerDay        = 24 * 60 * 60
	nanosecondsPerSecond = 1000000000
	nanosecondsPerDay    = secondsPerDay * nanosecondsPerSecond
)

// toDays converts a duration that might contain a fractional number of days
// into an integral number of days. Truncation occurs towards zero. This function
// is used when using durations for date arithmetic.
func toDays(duration time.Duration) time.Duration {
	days := duration.Nanoseconds() / nanosecondsPerDay
	nanoseconds := days * nanosecondsPerDay
	return time.Duration(nanoseconds)
}

// toSeconds converts a duration that might contain a fractional number of seconds
// into an integral number of seconds. Truncation occurs towards zero. This function
// is used when using durations for date-time and time arithmetic.
func toSeconds(duration time.Duration) time.Duration {
	seconds := duration.Nanoseconds() / nanosecondsPerSecond
	nanoseconds := seconds * nanosecondsPerSecond
	return time.Duration(nanoseconds)
}
