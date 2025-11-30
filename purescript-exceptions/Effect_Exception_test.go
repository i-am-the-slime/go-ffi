package purescript_exceptions

import (
	"testing"

	. "github.com/purescript-native/go-runtime"
)

func TestError(t *testing.T) {
	exports := Foreign("Effect.Exception")
	errorFn := exports["error"].(func(Any) Any)
	
	err := errorFn("Test error").(Dict)
	
	if err["message"] != "Test error" {
		t.Errorf("Expected message 'Test error', got '%v'", err["message"])
	}
}

func TestMessage(t *testing.T) {
	exports := Foreign("Effect.Exception")
	errorFn := exports["error"].(func(Any) Any)
	message := exports["message"].(func(Any) Any)
	
	err := errorFn("Test error")
	msg := message(err).(string)
	
	if msg != "Test error" {
		t.Errorf("Expected 'Test error', got '%s'", msg)
	}
}

func TestTry(t *testing.T) {
	exports := Foreign("Effect.Exception")
	tryFn := exports["try"].(func(Any) Any)
	
	// Test successful effect
	successEffect := func() Any {
		return 42
	}
	
	tryEffect := tryFn(successEffect).(func() Any)
	result := tryEffect().(Dict)
	
	if rightVal, ok := result["Right"]; !ok {
		t.Error("Expected Right for successful effect")
	} else if rightVal != 42 {
		t.Errorf("Expected 42, got %v", rightVal)
	}
}

func TestTryWithError(t *testing.T) {
	exports := Foreign("Effect.Exception")
	tryFn := exports["try"].(func(Any) Any)
	
	// Test failing effect
	failEffect := func() Any {
		panic(Dict{"message": "Something went wrong", "stack": ""})
	}
	
	tryEffect := tryFn(failEffect).(func() Any)
	result := tryEffect().(Dict)
	
	if leftVal, ok := result["Left"]; !ok {
		t.Error("Expected Left for failing effect")
	} else {
		err := leftVal.(Dict)
		if err["message"] != "Something went wrong" {
			t.Errorf("Expected error message 'Something went wrong', got '%v'", err["message"])
		}
	}
}

func TestCatchException(t *testing.T) {
	exports := Foreign("Effect.Exception")
	catchException := exports["catchException"].(func(Any, Any) Any)
	
	// Handler that returns a default value
	handler := func(err Any) Any {
		return func() Any {
			return "handled"
		}
	}
	
	// Effect that throws
	effect := func() Any {
		panic(Dict{"message": "error", "stack": ""})
	}
	
	catchEffect := catchException(handler, effect).(func() Any)
	result := catchEffect()
	
	if result != "handled" {
		t.Errorf("Expected 'handled', got '%v'", result)
	}
}

func TestCatchExceptionSuccess(t *testing.T) {
	exports := Foreign("Effect.Exception")
	catchException := exports["catchException"].(func(Any, Any) Any)
	
	// Handler that shouldn't be called
	handler := func(err Any) Any {
		return func() Any {
			return "handled"
		}
	}
	
	// Effect that succeeds
	effect := func() Any {
		return "success"
	}
	
	catchEffect := catchException(handler, effect).(func() Any)
	result := catchEffect()
	
	if result != "success" {
		t.Errorf("Expected 'success', got '%v'", result)
	}
}

func TestFinally(t *testing.T) {
	exports := Foreign("Effect.Exception")
	finally := exports["finally"].(func(Any, Any) Any)
	
	finalized := false
	
	// Finalizer
	finalizer := func() Any {
		finalized = true
		return nil
	}
	
	// Effect
	effect := func() Any {
		return "done"
	}
	
	finallyEffect := finally(finalizer, effect).(func() Any)
	result := finallyEffect()
	
	if result != "done" {
		t.Errorf("Expected 'done', got '%v'", result)
	}
	
	if !finalized {
		t.Error("Expected finalizer to be called")
	}
}

func TestFinallyWithError(t *testing.T) {
	exports := Foreign("Effect.Exception")
	finally := exports["finally"].(func(Any, Any) Any)
	
	finalized := false
	
	// Finalizer
	finalizer := func() Any {
		finalized = true
		return nil
	}
	
	// Effect that panics
	effect := func() Any {
		panic("error")
	}
	
	finallyEffect := finally(finalizer, effect).(func() Any)
	
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic to propagate")
		}
		
		if !finalized {
			t.Error("Expected finalizer to be called even on panic")
		}
	}()
	
	finallyEffect()
}

func TestThrow(t *testing.T) {
	exports := Foreign("Effect.Exception")
	throwFn := exports["throw"].(func(Any) Any)
	
	effect := throwFn("Test error").(func() Any)
	
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic")
		} else {
			err := r.(Dict)
			if err["message"] != "Test error" {
				t.Errorf("Expected 'Test error', got '%v'", err["message"])
			}
		}
	}()
	
	effect()
}

