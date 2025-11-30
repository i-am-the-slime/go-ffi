package purescript_routing

import (
	"net/url"
	"strconv"
	"strings"

	. "github.com/purescript-native/go-runtime"
)

// Route state that gets passed through the parser
type RouteState struct {
	path     []string
	params   Dict
	consumed []string
}

func init() {
	exports := Foreign("Routing.Match")

	// lit :: String -> Match Unit
	exports["lit"] = func(str_ Any) Any {
		return func(state_ Any) Any {
			str := str_.(string)
			state := state_.(RouteState)
			
			if len(state.path) == 0 {
				return Dict{"Left": "Expected path segment: " + str}
			}
			
			if state.path[0] == str {
				newState := RouteState{
					path:     state.path[1:],
					params:   state.params,
					consumed: append(state.consumed, str),
				}
				return Dict{"Right": []Any{Dict{}, newState}}
			}
			
			return Dict{"Left": "Expected '" + str + "', got '" + state.path[0] + "'"}
		}
	}

	// str :: Match String
	exports["str"] = func(state_ Any) Any {
		state := state_.(RouteState)
		
		if len(state.path) == 0 {
			return Dict{"Left": "Expected path segment"}
		}
		
		segment := state.path[0]
		newState := RouteState{
			path:     state.path[1:],
			params:   state.params,
			consumed: append(state.consumed, segment),
		}
		
		return Dict{"Right": []Any{segment, newState}}
	}

	// int :: Match Int
	exports["int"] = func(state_ Any) Any {
		state := state_.(RouteState)
		
		if len(state.path) == 0 {
			return Dict{"Left": "Expected integer path segment"}
		}
		
		segment := state.path[0]
		num, err := strconv.Atoi(segment)
		if err != nil {
			return Dict{"Left": "Expected integer, got '" + segment + "'"}
		}
		
		newState := RouteState{
			path:     state.path[1:],
			params:   state.params,
			consumed: append(state.consumed, segment),
		}
		
		return Dict{"Right": []Any{num, newState}}
	}

	// num :: Match Number
	exports["num"] = func(state_ Any) Any {
		state := state_.(RouteState)
		
		if len(state.path) == 0 {
			return Dict{"Left": "Expected numeric path segment"}
		}
		
		segment := state.path[0]
		num, err := strconv.ParseFloat(segment, 64)
		if err != nil {
			return Dict{"Left": "Expected number, got '" + segment + "'"}
		}
		
		newState := RouteState{
			path:     state.path[1:],
			params:   state.params,
			consumed: append(state.consumed, segment),
		}
		
		return Dict{"Right": []Any{num, newState}}
	}

	// bool :: Match Boolean
	exports["bool"] = func(state_ Any) Any {
		state := state_.(RouteState)
		
		if len(state.path) == 0 {
			return Dict{"Left": "Expected boolean path segment"}
		}
		
		segment := state.path[0]
		var result bool
		
		switch strings.ToLower(segment) {
		case "true", "1", "yes":
			result = true
		case "false", "0", "no":
			result = false
		default:
			return Dict{"Left": "Expected boolean, got '" + segment + "'"}
		}
		
		newState := RouteState{
			path:     state.path[1:],
			params:   state.params,
			consumed: append(state.consumed, segment),
		}
		
		return Dict{"Right": []Any{result, newState}}
	}

	// param :: String -> Match String
	exports["param"] = func(key_ Any) Any {
		return func(state_ Any) Any {
			key := key_.(string)
			state := state_.(RouteState)
			
			if val, ok := state.params[key]; ok {
				return Dict{"Right": []Any{val, state}}
			}
			
			return Dict{"Left": "Missing query parameter: " + key}
		}
	}

	// params :: Match (Object String)
	exports["params"] = func(state_ Any) Any {
		state := state_.(RouteState)
		return Dict{"Right": []Any{state.params, state}}
	}

	// root :: Match Unit
	exports["root"] = func(state_ Any) Any {
		state := state_.(RouteState)
		
		if len(state.path) == 0 || (len(state.path) == 1 && state.path[0] == "") {
			return Dict{"Right": []Any{Dict{}, state}}
		}
		
		return Dict{"Left": "Expected root path"}
	}

	// end :: Match Unit
	exports["end"] = func(state_ Any) Any {
		state := state_.(RouteState)
		
		if len(state.path) == 0 {
			return Dict{"Right": []Any{Dict{}, state}}
		}
		
		return Dict{"Left": "Expected end of path, but got: " + strings.Join(state.path, "/")}
	}

	// eitherMatch :: Match a -> Match b -> Match (Either a b)
	exports["eitherMatch"] = func(left_ Any, right_ Any) Any {
		return func(state Any) Any {
			leftFn := left_.(func(Any) Any)
			rightFn := right_.(func(Any) Any)
			
			// Try left parser first
			leftResult := leftFn(state)
			if resultDict, ok := leftResult.(Dict); ok {
				if _, hasRight := resultDict["Right"]; hasRight {
					// Left parser succeeded, wrap in Left
					arr := resultDict["Right"].([]Any)
					wrappedValue := Dict{"Left": arr[0]}
					return Dict{"Right": []Any{wrappedValue, arr[1]}}
				}
			}
			
			// Left failed, try right parser
			rightResult := rightFn(state)
			if resultDict, ok := rightResult.(Dict); ok {
				if _, hasRight := resultDict["Right"]; hasRight {
					// Right parser succeeded, wrap in Right
					arr := resultDict["Right"].([]Any)
					wrappedValue := Dict{"Right": arr[0]}
					return Dict{"Right": []Any{wrappedValue, arr[1]}}
				}
			}
			
			return Dict{"Left": "Both alternatives failed"}
		}
	}

	// fail :: String -> Match a
	exports["fail"] = func(msg_ Any) Any {
		return func(state Any) Any {
			msg := msg_.(string)
			return Dict{"Left": msg}
		}
	}
	
	// Also register core Routing module functions
	routingExports := Foreign("Routing")

	// match :: Match a -> String -> Either String a
	routingExports["match"] = func(matcher_ Any, url_ Any) Any {
		matcher := matcher_.(func(Any) Any)
		urlStr := url_.(string)
		
		// Parse URL
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return Dict{"Left": "Invalid URL: " + err.Error()}
		}
		
		// Split path into segments
		pathStr := strings.Trim(parsed.Path, "/")
		var segments []string
		if pathStr != "" {
			segments = strings.Split(pathStr, "/")
		}
		
		// Parse query parameters
		params := make(Dict)
		for key, values := range parsed.Query() {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}
		
		// Create initial state
		state := RouteState{
			path:     segments,
			params:   params,
			consumed: []string{},
		}
		
		// Run the matcher
		result := matcher(state)
		
		if resultDict, ok := result.(Dict); ok {
			if rightVal, hasRight := resultDict["Right"]; hasRight {
				arr := rightVal.([]Any)
				return Dict{"Right": arr[0]}
			}
			if leftVal, hasLeft := resultDict["Left"]; hasLeft {
				return Dict{"Left": leftVal}
			}
		}
		
		return Dict{"Left": "Route matching failed"}
	}

	// matchWith :: (a -> b -> b) -> b -> Match a -> String -> b
	routingExports["matchWith"] = func(fn_ Any, default_ Any, matcher_ Any, url_ Any) Any {
		fn := fn_.(func(Any, Any) Any)
		matcher := matcher_.(func(Any) Any)
		urlStr := url_.(string)
		
		matchResult := routingExports["match"].(func(Any, Any) Any)(matcher, urlStr)
		
		if resultDict, ok := matchResult.(Dict); ok {
			if rightVal, hasRight := resultDict["Right"]; hasRight {
				return fn(rightVal, default_)
			}
		}
		
		return default_
	}
}

