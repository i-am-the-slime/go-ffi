package purescript_now

import (
	"time"

	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Effect.Now")

	// now :: Effect Instant
	// Returns the current time as an Instant (milliseconds since Unix epoch)
	exports["now"] = func() Any {
		// Get current time in milliseconds since Unix epoch
		nowMs := float64(time.Now().UnixNano()) / 1e6
		// Return as Milliseconds constructor: data Instant = Instant Milliseconds
		return Apply(Dict{"Milliseconds": nowMs}["Milliseconds"], nowMs)
	}

	// getTimezoneOffset :: Effect Minutes
	// Returns the timezone offset in minutes from UTC
	exports["getTimezoneOffset"] = func() Any {
		_, offsetSeconds := time.Now().Zone()
		// Convert seconds to minutes and negate (to match JavaScript convention)
		offsetMinutes := float64(-offsetSeconds / 60)
		// Return as Minutes constructor: newtype Minutes = Minutes Number
		return Apply(Dict{"Minutes": offsetMinutes}["Minutes"], offsetMinutes)
	}
}


