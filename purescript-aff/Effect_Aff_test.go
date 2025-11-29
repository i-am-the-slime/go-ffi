package purescript_aff

import (
	"fmt"
	"testing"
	"time"

	. "github.com/purescript-native/go-runtime"
)

// Test utility functions
func makeUtil() Dict {
	return Dict{
		"isLeft": func(e Any) Any {
			if dict, ok := e.(Dict); ok {
				_, hasLeft := dict["Left"]
				return hasLeft
			}
			return false
		},
		"fromLeft": func(e Any) Any {
			if dict, ok := e.(Dict); ok {
				if val, hasLeft := dict["Left"]; hasLeft {
					return val
				}
			}
			return nil
		},
		"fromRight": func(e Any) Any {
			if dict, ok := e.(Dict); ok {
				if val, hasRight := dict["Right"]; hasRight {
					return val
				}
			}
			return nil
		},
		"left": func(e Any) Any {
			return Dict{"Left": e}
		},
		"right": func(v Any) Any {
			return Dict{"Right": v}
		},
	}
}

func TestPureAff(t *testing.T) {
	util := makeUtil()
	
	// Create a Pure Aff
	aff := Pure{value: 42}
	
	// Create a fiber
	fiber := Fiber(util, nil, aff)
	fiberDict := fiber.(Dict)
	
	// Track completion
	completed := false
	var result Any
	
	// Register completion handler
	onComplete := fiberDict["onComplete"].(func(OnComplete) func() Any)
	onComplete(OnComplete{
		rethrow: false,
		handler: func(res Any) func() Any {
			return func() Any {
				completed = true
				result = res
				return nil
			}
		},
	})()
	
	// Run the fiber
	runFn := fiberDict["run"].(func() Any)
	runFn()
	
	// Check result
	if !completed {
		t.Fatal("Fiber did not complete")
	}
	
	// Result should be Right 42
	resultDict := result.(Dict)
	if _, hasRight := resultDict["Right"]; !hasRight {
		t.Fatal("Result should be Right")
	}
	
	if resultDict["Right"] != 42 {
		t.Fatalf("Expected 42, got %v", resultDict["Right"])
	}
	
	t.Log("✓ Pure Aff works correctly")
}

func TestBindAff(t *testing.T) {
	util := makeUtil()
	
	// Create Pure(10) >>= (\x -> Pure(x * 2))
	aff1 := Pure{value: 10}
	aff2 := Bind{
		affOfB: aff1,
		bToAff: func(x Any) Any {
			return Pure{value: x.(int) * 2}
		},
	}
	
	// Create a fiber
	fiber := Fiber(util, nil, aff2)
	fiberDict := fiber.(Dict)
	
	completed := false
	var result Any
	
	onComplete := fiberDict["onComplete"].(func(OnComplete) func() Any)
	onComplete(OnComplete{
		rethrow: false,
		handler: func(res Any) func() Any {
			return func() Any {
				completed = true
				result = res
				return nil
			}
		},
	})()
	
	runFn := fiberDict["run"].(func() Any)
	runFn()
	
	if !completed {
		t.Fatal("Fiber did not complete")
	}
	
	resultDict := result.(Dict)
	if resultDict["Right"] != 20 {
		t.Fatalf("Expected 20, got %v", resultDict["Right"])
	}
	
	t.Log("✓ Bind Aff works correctly")
}

func TestAsyncAff(t *testing.T) {
	util := makeUtil()
	right := util["right"].(func(Any) Any)
	
	// Create an async Aff that completes after a short delay
	aff := Async{
		asyncFn: func(cb Any) Any {
			return func() Any {
				// Simulate async work
				go func() {
					time.Sleep(10 * time.Millisecond)
					// Call callback with result
					effect := cb.(func(Any) func() Any)(right(123))
					effect()
				}()
				
				// Return canceler
				return func(error Any) Any {
					return Pure{value: Dict{}}
				}
			}
		},
	}
	
	fiber := Fiber(util, nil, aff)
	fiberDict := fiber.(Dict)
	
	completed := false
	var result Any
	done := make(chan bool)
	
	onComplete := fiberDict["onComplete"].(func(OnComplete) func() Any)
	onComplete(OnComplete{
		rethrow: false,
		handler: func(res Any) func() Any {
			return func() Any {
				completed = true
				result = res
				done <- true
				return nil
			}
		},
	})()
	
	runFn := fiberDict["run"].(func() Any)
	runFn()
	
	// Wait for completion with timeout
	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Async Aff timed out")
	}
	
	if !completed {
		t.Fatal("Fiber did not complete")
	}
	
	resultDict := result.(Dict)
	if resultDict["Right"] != 123 {
		t.Fatalf("Expected 123, got %v", resultDict["Right"])
	}
	
	t.Log("✓ Async Aff works correctly")
}

func TestCatchError(t *testing.T) {
	util := makeUtil()
	
	// Create Throw >>= Catch
	throwAff := Throw{err: fmt.Errorf("test error")}
	catchAff := Catch{
		aff: throwAff,
		errorToAff: func(e Any) Any {
			return Pure{value: "caught"}
		},
	}
	
	fiber := Fiber(util, nil, catchAff)
	fiberDict := fiber.(Dict)
	
	completed := false
	var result Any
	
	onComplete := fiberDict["onComplete"].(func(OnComplete) func() Any)
	onComplete(OnComplete{
		rethrow: false,
		handler: func(res Any) func() Any {
			return func() Any {
				completed = true
				result = res
				return nil
			}
		},
	})()
	
	runFn := fiberDict["run"].(func() Any)
	runFn()
	
	if !completed {
		t.Fatal("Fiber did not complete")
	}
	
	resultDict := result.(Dict)
	if resultDict["Right"] != "caught" {
		t.Fatalf("Expected 'caught', got %v", resultDict["Right"])
	}
	
	t.Log("✓ Catch Error works correctly")
}

func TestParallelMap(t *testing.T) {
	
	util := makeUtil()
	
	// Create ParMap that doubles a value
	parAff := ParMap{
		bToA:      func(x Any) Any { return x.(int) * 2 },
		parAffOfB: Pure{value: 21},
		result:    EMPTY,
	}
	
	// Convert to sequential Aff
	seqAff := Sequential{parAff: parAff}
	
	fiber := Fiber(util, nil, seqAff)
	fiberDict := fiber.(Dict)
	
	completed := false
	var result Any
	done := make(chan bool, 1)
	
	onComplete := fiberDict["onComplete"].(func(OnComplete) func() Any)
	onComplete(OnComplete{
		rethrow: false,
		handler: func(res Any) func() Any {
			return func() Any {
				completed = true
				result = res
				select {
				case done <- true:
				default:
				}
				return nil
			}
		},
	})()
	
	runFn := fiberDict["run"].(func() Any)
	runFn()
	
	// Wait for completion
	select {
	case <-done:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Parallel Map timed out - likely deadlock in runPar")
	}
	
	if !completed {
		t.Fatal("Fiber did not complete")
	}
	
	resultDict := result.(Dict)
	if resultDict["Right"] != 42 {
		t.Fatalf("Expected 42, got %v", resultDict["Right"])
	}
	
	t.Log("✓ Parallel Map works correctly")
}

func TestSupervisor(t *testing.T) {
	
	util := makeUtil()
	
	supervisor := SupervisorNew(util)
	supervisorDict := supervisor.(Dict)
	
	// Check isEmpty
	isEmpty := supervisorDict["isEmpty"].(func() Any)()
	if !isEmpty.(bool) {
		t.Fatal("New supervisor should be empty")
	}
	
	// Register a fiber
	fiber := Fiber(util, supervisor, Pure{value: 1})
	registerFn := supervisorDict["register"].(func(Any))
	registerFn(fiber)
	
	// Should not be empty now
	isEmpty = supervisorDict["isEmpty"].(func() Any)()
	if isEmpty.(bool) {
		t.Fatal("Supervisor should not be empty after registration")
	}
	
	// Run the fiber to completion
	fiberDict := fiber.(Dict)
	runFn := fiberDict["run"].(func() Any)
	runFn()
	
	// Give it a moment to complete
	time.Sleep(10 * time.Millisecond)
	
	// Should be empty again
	isEmpty = supervisorDict["isEmpty"].(func() Any)()
	if !isEmpty.(bool) {
		t.Fatal("Supervisor should be empty after fiber completes")
	}
	
	t.Log("✓ Supervisor works correctly")
}

func TestEffectQueue(t *testing.T) {
	// Test the effect queue mechanism
	executed := false
	
	QueueEffect(func() Any {
		executed = true
		return nil
	})
	
	DrainEffectQueue()
	
	if !executed {
		t.Fatal("Effect was not executed")
	}
	
	t.Log("✓ Effect queue works correctly")
}

