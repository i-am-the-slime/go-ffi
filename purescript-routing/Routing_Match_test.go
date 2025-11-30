package purescript_routing

import (
	"testing"

	. "github.com/purescript-native/go-runtime"
)

func TestLitMatch(t *testing.T) {
	exports := Foreign("Routing.Match")
	routingExports := Foreign("Routing")
	
	lit := exports["lit"].(func(Any) Any)
	match := routingExports["match"].(func(Any, Any) Any)
	
	// Create a matcher for "/hello"
	matcher := lit("hello")
	
	// Test matching URL
	result := match(matcher, "/hello").(Dict)
	if _, ok := result["Right"]; !ok {
		t.Errorf("Expected Right, got Left: %v", result["Left"])
	}
	
	// Test non-matching URL
	result2 := match(matcher, "/goodbye").(Dict)
	if _, ok := result2["Left"]; !ok {
		t.Error("Expected Left for non-matching path")
	}
}

func TestStrMatch(t *testing.T) {
	exports := Foreign("Routing.Match")
	routingExports := Foreign("Routing")
	
	str := exports["str"].(func(Any) Any)
	match := routingExports["match"].(func(Any, Any) Any)
	
	// Test matching any string segment
	result := match(str, "/anything").(Dict)
	if rightVal, ok := result["Right"]; ok {
		if rightVal != "anything" {
			t.Errorf("Expected 'anything', got %v", rightVal)
		}
	} else {
		t.Error("Expected Right for string match")
	}
}

func TestIntMatch(t *testing.T) {
	exports := Foreign("Routing.Match")
	routingExports := Foreign("Routing")
	
	intMatch := exports["int"].(func(Any) Any)
	match := routingExports["match"].(func(Any, Any) Any)
	
	// Test matching integer
	result := match(intMatch, "/42").(Dict)
	if rightVal, ok := result["Right"]; ok {
		if rightVal != 42 {
			t.Errorf("Expected 42, got %v", rightVal)
		}
	} else {
		t.Error("Expected Right for int match")
	}
	
	// Test non-integer
	result2 := match(intMatch, "/notanumber").(Dict)
	if _, ok := result2["Left"]; !ok {
		t.Error("Expected Left for non-integer")
	}
}

func TestNumMatch(t *testing.T) {
	exports := Foreign("Routing.Match")
	routingExports := Foreign("Routing")
	
	num := exports["num"].(func(Any) Any)
	match := routingExports["match"].(func(Any, Any) Any)
	
	// Test matching float
	result := match(num, "/3.14").(Dict)
	if rightVal, ok := result["Right"]; ok {
		if rightVal != 3.14 {
			t.Errorf("Expected 3.14, got %v", rightVal)
		}
	} else {
		t.Error("Expected Right for number match")
	}
}

func TestBoolMatch(t *testing.T) {
	exports := Foreign("Routing.Match")
	routingExports := Foreign("Routing")
	
	boolMatch := exports["bool"].(func(Any) Any)
	match := routingExports["match"].(func(Any, Any) Any)
	
	// Test true
	result := match(boolMatch, "/true").(Dict)
	if rightVal, ok := result["Right"]; ok {
		if rightVal != true {
			t.Errorf("Expected true, got %v", rightVal)
		}
	} else {
		t.Error("Expected Right for bool match")
	}
	
	// Test false
	result2 := match(boolMatch, "/false").(Dict)
	if rightVal, ok := result2["Right"]; ok {
		if rightVal != false {
			t.Errorf("Expected false, got %v", rightVal)
		}
	} else {
		t.Error("Expected Right for bool match")
	}
}

func TestParamMatch(t *testing.T) {
	exports := Foreign("Routing.Match")
	routingExports := Foreign("Routing")
	
	param := exports["param"].(func(Any) Any)
	match := routingExports["match"].(func(Any, Any) Any)
	
	// Test query parameter
	matcher := param("name")
	result := match(matcher, "/?name=test").(Dict)
	if rightVal, ok := result["Right"]; ok {
		if rightVal != "test" {
			t.Errorf("Expected 'test', got %v", rightVal)
		}
	} else {
		t.Error("Expected Right for param match")
	}
	
	// Test missing parameter
	result2 := match(matcher, "/?other=value").(Dict)
	if _, ok := result2["Left"]; !ok {
		t.Error("Expected Left for missing parameter")
	}
}

func TestParamsMatch(t *testing.T) {
	exports := Foreign("Routing.Match")
	routingExports := Foreign("Routing")
	
	params := exports["params"].(func(Any) Any)
	match := routingExports["match"].(func(Any, Any) Any)
	
	// Test all query parameters
	result := match(params, "/?key1=value1&key2=value2").(Dict)
	if rightVal, ok := result["Right"]; ok {
		paramsDict := rightVal.(Dict)
		if paramsDict["key1"] != "value1" || paramsDict["key2"] != "value2" {
			t.Errorf("Expected both parameters, got %v", paramsDict)
		}
	} else {
		t.Error("Expected Right for params match")
	}
}

func TestRootMatch(t *testing.T) {
	exports := Foreign("Routing.Match")
	routingExports := Foreign("Routing")
	
	root := exports["root"].(func(Any) Any)
	match := routingExports["match"].(func(Any, Any) Any)
	
	// Test root path
	result := match(root, "/").(Dict)
	if _, ok := result["Right"]; !ok {
		t.Error("Expected Right for root match")
	}
	
	// Test non-root
	result2 := match(root, "/something").(Dict)
	if _, ok := result2["Left"]; !ok {
		t.Error("Expected Left for non-root path")
	}
}

func TestEndMatch(t *testing.T) {
	exports := Foreign("Routing.Match")
	routingExports := Foreign("Routing")
	
	end := exports["end"].(func(Any) Any)
	match := routingExports["match"].(func(Any, Any) Any)
	
	// Test end of path
	result := match(end, "/").(Dict)
	if _, ok := result["Right"]; !ok {
		t.Error("Expected Right for end match on empty path")
	}
}

// Helper to create a chained matcher (simulating bind/>>= in PureScript)
func chainMatchers(matcher1 func(Any) Any, matcher2 func(Any) Any) func(Any) Any {
	return func(state Any) Any {
		result1 := matcher1(state)
		if dict, ok := result1.(Dict); ok {
			if rightVal, ok := dict["Right"]; ok {
				arr := rightVal.([]Any)
				newState := arr[1]
				return matcher2(newState)
			}
			return result1
		}
		return result1
	}
}

func TestComposedRoute(t *testing.T) {
	exports := Foreign("Routing.Match")
	routingExports := Foreign("Routing")
	
	lit := exports["lit"].(func(Any) Any)
	intMatch := exports["int"].(func(Any) Any)
	match := routingExports["match"].(func(Any, Any) Any)
	
	// Create matcher for "/users/123"
	matcher := chainMatchers(
		lit("users").(func(Any) Any),
		intMatch,
	)
	
	result := match(matcher, "/users/42").(Dict)
	if rightVal, ok := result["Right"]; ok {
		if rightVal != 42 {
			t.Errorf("Expected 42, got %v", rightVal)
		}
	} else {
		t.Errorf("Expected Right for composed route, got Left: %v", result["Left"])
	}
}

