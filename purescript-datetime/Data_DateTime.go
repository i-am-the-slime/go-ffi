package purescript_datetime

import (
	"strings"
	"time"

	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Data.DateTime")

	// DateTime is represented as time.Time in Go
	
	// now :: Effect DateTime
	exports["now"] = func() Any {
		return func() Any {
			return time.Now()
		}
	}

	// fromString :: String -> Maybe DateTime
	exports["fromString"] = func(str_ Any) Any {
		str := str_.(string)
		
		// Try various common formats
		formats := []string{
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}
		
		for _, format := range formats {
			if t, err := time.Parse(format, str); err == nil {
				return Dict{"value0": t} // Just
			}
		}
		
		return Dict{} // Nothing
	}

	// toString :: DateTime -> String
	exports["toString"] = func(dt_ Any) Any {
		dt := dt_.(time.Time)
		return dt.Format(time.RFC3339)
	}

	// year :: DateTime -> Int
	exports["year"] = func(dt_ Any) Any {
		dt := dt_.(time.Time)
		return dt.Year()
	}

	// month :: DateTime -> Int
	exports["month"] = func(dt_ Any) Any {
		dt := dt_.(time.Time)
		return int(dt.Month())
	}

	// day :: DateTime -> Int
	exports["day"] = func(dt_ Any) Any {
		dt := dt_.(time.Time)
		return dt.Day()
	}

	// hour :: DateTime -> Int
	exports["hour"] = func(dt_ Any) Any {
		dt := dt_.(time.Time)
		return dt.Hour()
	}

	// minute :: DateTime -> Int
	exports["minute"] = func(dt_ Any) Any {
		dt := dt_.(time.Time)
		return dt.Minute()
	}

	// second :: DateTime -> Int
	exports["second"] = func(dt_ Any) Any {
		dt := dt_.(time.Time)
		return dt.Second()
	}

	// millisecond :: DateTime -> Int
	exports["millisecond"] = func(dt_ Any) Any {
		dt := dt_.(time.Time)
		return dt.Nanosecond() / 1000000
	}

	// toEpochMilliseconds :: DateTime -> Number
	exports["toEpochMilliseconds"] = func(dt_ Any) Any {
		dt := dt_.(time.Time)
		return float64(dt.UnixMilli())
	}

	// fromEpochMilliseconds :: Number -> DateTime
	exports["fromEpochMilliseconds"] = func(ms_ Any) Any {
		ms := int64(ms_.(float64))
		return time.UnixMilli(ms)
	}

	// diff :: DateTime -> DateTime -> Number
	exports["diff"] = func(dt1_ Any, dt2_ Any) Any {
		dt1 := dt1_.(time.Time)
		dt2 := dt2_.(time.Time)
		return float64(dt1.Sub(dt2).Milliseconds())
	}

	// add :: Number -> DateTime -> DateTime
	exports["add"] = func(ms_ Any, dt_ Any) Any {
		ms := int64(ms_.(float64))
		dt := dt_.(time.Time)
		return dt.Add(time.Duration(ms) * time.Millisecond)
	}

	// format :: String -> DateTime -> String
	exports["format"] = func(fmt_ Any, dt_ Any) Any {
		format := fmt_.(string)
		dt := dt_.(time.Time)
		
		// Convert common format strings to Go format
		goFormat := convertDateFormat(format)
		return dt.Format(goFormat)
	}
}

// convertDateFormat converts common date format strings to Go's format
func convertDateFormat(format string) string {
	// Simple conversions for common patterns
	// This is a basic implementation - could be extended
	replacements := map[string]string{
		"YYYY": "2006",
		"YY":   "06",
		"MM":   "01",
		"DD":   "02",
		"HH":   "15",
		"mm":   "04",
		"ss":   "05",
		"SSS":  "000",
	}
	
	result := format
	for old, new := range replacements {
		result = strings.Replace(result, old, new, -1)
	}
	
	return result
}

