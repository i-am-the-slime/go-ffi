package purescript_datetime

import (
	"testing"
	"time"

	. "github.com/purescript-native/go-runtime"
)

func TestNow(t *testing.T) {
	exports := Foreign("Data.DateTime")
	now := exports["now"].(func() Any)
	
	before := time.Now()
	effect := now().(func() Any)
	dt := effect().(time.Time)
	after := time.Now()
	
	if dt.Before(before) || dt.After(after) {
		t.Error("Now() returned time outside expected range")
	}
}

func TestFromString(t *testing.T) {
	exports := Foreign("Data.DateTime")
	fromString := exports["fromString"].(func(Any) Any)
	
	result := fromString("2023-12-25T10:30:00Z").(Dict)
	
	if _, ok := result["value0"]; !ok {
		t.Error("Expected Just for valid date string")
	}
	
	dt := result["value0"].(time.Time)
	if dt.Year() != 2023 || dt.Month() != 12 || dt.Day() != 25 {
		t.Errorf("Expected 2023-12-25, got %v", dt)
	}
}

func TestFromStringInvalid(t *testing.T) {
	exports := Foreign("Data.DateTime")
	fromString := exports["fromString"].(func(Any) Any)
	
	result := fromString("not a date").(Dict)
	
	if len(result) != 0 {
		t.Error("Expected Nothing for invalid date string")
	}
}

func TestToString(t *testing.T) {
	exports := Foreign("Data.DateTime")
	toString := exports["toString"].(func(Any) Any)
	
	dt := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
	result := toString(dt).(string)
	
	if result == "" {
		t.Error("Expected non-empty string")
	}
}

func TestYear(t *testing.T) {
	exports := Foreign("Data.DateTime")
	year := exports["year"].(func(Any) Any)
	
	dt := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
	result := year(dt).(int)
	
	if result != 2023 {
		t.Errorf("Expected 2023, got %d", result)
	}
}

func TestMonth(t *testing.T) {
	exports := Foreign("Data.DateTime")
	month := exports["month"].(func(Any) Any)
	
	dt := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
	result := month(dt).(int)
	
	if result != 12 {
		t.Errorf("Expected 12, got %d", result)
	}
}

func TestDay(t *testing.T) {
	exports := Foreign("Data.DateTime")
	day := exports["day"].(func(Any) Any)
	
	dt := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
	result := day(dt).(int)
	
	if result != 25 {
		t.Errorf("Expected 25, got %d", result)
	}
}

func TestHour(t *testing.T) {
	exports := Foreign("Data.DateTime")
	hour := exports["hour"].(func(Any) Any)
	
	dt := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
	result := hour(dt).(int)
	
	if result != 10 {
		t.Errorf("Expected 10, got %d", result)
	}
}

func TestMinute(t *testing.T) {
	exports := Foreign("Data.DateTime")
	minute := exports["minute"].(func(Any) Any)
	
	dt := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
	result := minute(dt).(int)
	
	if result != 30 {
		t.Errorf("Expected 30, got %d", result)
	}
}

func TestToEpochMilliseconds(t *testing.T) {
	exports := Foreign("Data.DateTime")
	toEpoch := exports["toEpochMilliseconds"].(func(Any) Any)
	
	dt := time.Date(1970, 1, 1, 0, 0, 1, 0, time.UTC)
	result := toEpoch(dt).(float64)
	
	if result != 1000.0 {
		t.Errorf("Expected 1000ms, got %f", result)
	}
}

func TestFromEpochMilliseconds(t *testing.T) {
	exports := Foreign("Data.DateTime")
	fromEpoch := exports["fromEpochMilliseconds"].(func(Any) Any)
	
	dt := fromEpoch(1000.0).(time.Time)
	
	if dt.Unix() != 1 {
		t.Errorf("Expected 1 second since epoch, got %d", dt.Unix())
	}
}

func TestDiff(t *testing.T) {
	exports := Foreign("Data.DateTime")
	diff := exports["diff"].(func(Any, Any) Any)
	
	dt1 := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
	dt2 := time.Date(2023, 12, 25, 10, 29, 0, 0, time.UTC)
	
	result := diff(dt1, dt2).(float64)
	
	if result != 60000.0 { // 1 minute = 60000 ms
		t.Errorf("Expected 60000ms, got %f", result)
	}
}

func TestAdd(t *testing.T) {
	exports := Foreign("Data.DateTime")
	add := exports["add"].(func(Any, Any) Any)
	
	dt := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
	result := add(60000.0, dt).(time.Time) // Add 1 minute
	
	if result.Minute() != 31 {
		t.Errorf("Expected minute 31, got %d", result.Minute())
	}
}

