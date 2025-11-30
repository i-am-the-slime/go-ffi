package purescript_simple_json

import (
	"testing"

	. "github.com/purescript-native/go-runtime"
)

func TestParseJSON(t *testing.T) {
	exports := Foreign("Simple.JSON")
	parseJSON := exports["parseJSON"].(func(Any) Any)
	
	// Test valid JSON
	result := parseJSON(`{"name":"test","value":42}`).(Dict)
	if _, ok := result["Right"]; !ok {
		t.Error("parseJSON should return Right for valid JSON")
	}
	
	// Test invalid JSON
	result2 := parseJSON(`{invalid json}`).(Dict)
	if _, ok := result2["Left"]; !ok {
		t.Error("parseJSON should return Left for invalid JSON")
	}
}

func TestWriteJSON(t *testing.T) {
	exports := Foreign("Simple.JSON")
	writeJSON := exports["writeJSON"].(func(Any) Any)
	
	input := Dict{
		"name":  "test",
		"value": 42,
		"flag":  true,
	}
	
	result := writeJSON(input).(string)
	
	// Should be valid JSON string
	if result == "" {
		t.Error("writeJSON should not return empty string")
	}
	
	// Should contain the values
	if !contains(result, "test") || !contains(result, "42") {
		t.Errorf("writeJSON output missing expected values: %s", result)
	}
}

func TestWriteJSONPretty(t *testing.T) {
	exports := Foreign("Simple.JSON")
	writeJSON := exports["writeJSON'"].(func(Any) Any)
	
	input := Dict{
		"name": "test",
		"nested": Dict{
			"value": 123,
		},
	}
	
	result := writeJSON(input).(string)
	
	// Should have newlines (pretty printed)
	if !contains(result, "\n") {
		t.Error("writeJSON' should pretty print with newlines")
	}
}

func TestReadString(t *testing.T) {
	exports := Foreign("Simple.JSON")
	readString := exports["readString"].(func(Any) Any)
	
	// Test valid string
	result := readString("hello").(Dict)
	if val, ok := result["Right"].(string); !ok || val != "hello" {
		t.Error("readString should return Right for string value")
	}
	
	// Test non-string
	result2 := readString(42).(Dict)
	if _, ok := result2["Left"]; !ok {
		t.Error("readString should return Left for non-string")
	}
}

func TestReadNumber(t *testing.T) {
	exports := Foreign("Simple.JSON")
	readNumber := exports["readNumber"].(func(Any) Any)
	
	// Test float
	result := readNumber(42.5).(Dict)
	if val, ok := result["Right"].(float64); !ok || val != 42.5 {
		t.Error("readNumber should return Right for float value")
	}
	
	// Test int
	result2 := readNumber(42).(Dict)
	if _, ok := result2["Right"]; !ok {
		t.Error("readNumber should return Right for int value")
	}
	
	// Test non-number
	result3 := readNumber("not a number").(Dict)
	if _, ok := result3["Left"]; !ok {
		t.Error("readNumber should return Left for non-number")
	}
}

func TestReadInt(t *testing.T) {
	exports := Foreign("Simple.JSON")
	readInt := exports["readInt"].(func(Any) Any)
	
	// Test valid int
	result := readInt(42).(Dict)
	if val, ok := result["Right"].(int); !ok || val != 42 {
		t.Error("readInt should return Right for int value")
	}
	
	// Test float that's actually an int
	result2 := readInt(float64(100)).(Dict)
	if val, ok := result2["Right"].(int); !ok || val != 100 {
		t.Error("readInt should convert whole number floats")
	}
	
	// Test non-integer float
	result3 := readInt(42.5).(Dict)
	if _, ok := result3["Left"]; !ok {
		t.Error("readInt should return Left for non-integer float")
	}
}

func TestReadBoolean(t *testing.T) {
	exports := Foreign("Simple.JSON")
	readBoolean := exports["readBoolean"].(func(Any) Any)
	
	// Test true
	result := readBoolean(true).(Dict)
	if val, ok := result["Right"].(bool); !ok || val != true {
		t.Error("readBoolean should return Right for bool value")
	}
	
	// Test non-boolean
	result2 := readBoolean("true").(Dict)
	if _, ok := result2["Left"]; !ok {
		t.Error("readBoolean should return Left for non-bool")
	}
}

func TestReadArray(t *testing.T) {
	exports := Foreign("Simple.JSON")
	readArray := exports["readArray"].(func(Any) Any)
	
	// Test valid array
	arr := []Any{1, 2, 3}
	result := readArray(arr).(Dict)
	if val, ok := result["Right"].([]Any); !ok || len(val) != 3 {
		t.Error("readArray should return Right for array value")
	}
	
	// Test non-array
	result2 := readArray("not an array").(Dict)
	if _, ok := result2["Left"]; !ok {
		t.Error("readArray should return Left for non-array")
	}
}

func TestReadNull(t *testing.T) {
	exports := Foreign("Simple.JSON")
	readNull := exports["readNull"].(func(Any) Any)
	
	// Test null
	result := readNull(nil).(Dict)
	if _, ok := result["Right"]; !ok {
		t.Error("readNull should return Right for nil value")
	}
	
	// Test non-null
	result2 := readNull("not null").(Dict)
	if _, ok := result2["Left"]; !ok {
		t.Error("readNull should return Left for non-null")
	}
}

func TestRoundTrip(t *testing.T) {
	exports := Foreign("Simple.JSON")
	writeJSON := exports["writeJSON"].(func(Any) Any)
	parseJSON := exports["parseJSON"].(func(Any) Any)
	
	input := Dict{
		"string": "hello",
		"number": 42.5,
		"bool":   true,
		"array":  []Any{1, 2, 3},
		"nested": Dict{
			"key": "value",
		},
	}
	
	// Write to JSON
	jsonStr := writeJSON(input).(string)
	
	// Parse back
	result := parseJSON(jsonStr).(Dict)
	if _, ok := result["Right"]; !ok {
		t.Error("Round trip failed: could not parse generated JSON")
	}
	
	parsed := result["Right"].(Dict)
	if parsed["string"] != "hello" {
		t.Error("Round trip failed: string value mismatch")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && 
		(s[0:len(substr)] == substr || contains(s[1:], substr)))
}

