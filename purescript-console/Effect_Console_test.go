package purescript_console

import (
	"bytes"
	"os"
	"testing"

	. "github.com/purescript-native/go-runtime"
)

func TestLog(t *testing.T) {
	exports := Foreign("Effect.Console")
	log := exports["log"].(func(Any) Any)
	
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	effect := log("Test message").(func() Any)
	effect()
	
	w.Close()
	os.Stdout = old
	
	var buf bytes.Buffer
	buf.ReadFrom(r)
	
	output := buf.String()
	if output != "Test message\n" {
		t.Errorf("Expected 'Test message\\n', got '%s'", output)
	}
}

func TestLogShow(t *testing.T) {
	exports := Foreign("Effect.Console")
	logShow := exports["logShow"].(func(Any) Any)
	
	effect := logShow(42).(func() Any)
	effect()
	// Just checking it doesn't crash
}

func TestWarn(t *testing.T) {
	exports := Foreign("Effect.Console")
	warn := exports["warn"].(func(Any) Any)
	
	effect := warn("Warning message").(func() Any)
	effect()
	// Just checking it doesn't crash
}

func TestError(t *testing.T) {
	exports := Foreign("Effect.Console")
	errorFn := exports["error"].(func(Any) Any)
	
	effect := errorFn("Error message").(func() Any)
	effect()
	// Just checking it doesn't crash
}

func TestInfo(t *testing.T) {
	exports := Foreign("Effect.Console")
	info := exports["info"].(func(Any) Any)
	
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	effect := info("Info message").(func() Any)
	effect()
	
	w.Close()
	os.Stdout = old
	
	var buf bytes.Buffer
	buf.ReadFrom(r)
	
	output := buf.String()
	if output != "[INFO] Info message\n" {
		t.Errorf("Expected '[INFO] Info message\\n', got '%s'", output)
	}
}

func TestDebug(t *testing.T) {
	exports := Foreign("Effect.Console")
	debug := exports["debug"].(func(Any) Any)
	
	effect := debug("Debug message").(func() Any)
	effect()
	// Just checking it doesn't crash
}

func TestTime(t *testing.T) {
	exports := Foreign("Effect.Console")
	timeFn := exports["time"].(func(Any) Any)
	timeEnd := exports["timeEnd"].(func(Any) Any)
	
	effect1 := timeFn("test-timer").(func() Any)
	effect1()
	
	effect2 := timeEnd("test-timer").(func() Any)
	effect2()
	// Just checking it doesn't crash
}

func TestClear(t *testing.T) {
	exports := Foreign("Effect.Console")
	clear := exports["clear"].(func() Any)
	
	effect := clear().(func() Any)
	effect()
	// Just checking it doesn't crash
}

