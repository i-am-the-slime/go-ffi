package purescript_exceptions

import (
	"fmt"

	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Effect.Exception")

	// Error type - represented as a dict with message and stack
	type Error struct {
		Message string
		Stack   string
	}

	// error :: String -> Error
	exports["error"] = func(msg_ Any) Any {
		msg := msg_.(string)
		return Dict{
			"message": msg,
			"stack":   "",
		}
	}

	// message :: Error -> String
	exports["message"] = func(err_ Any) Any {
		if err, ok := err_.(Dict); ok {
			if msg, ok := err["message"].(string); ok {
				return msg
			}
		}
		return fmt.Sprintf("%v", err_)
	}

	// stack :: Error -> Maybe String
	exports["stack"] = func(err_ Any) Any {
		if err, ok := err_.(Dict); ok {
			if stack, ok := err["stack"].(string); ok && stack != "" {
				return Dict{"value0": stack} // Just
			}
		}
		return Dict{} // Nothing
	}

	// throwException :: forall a. Error -> Effect a
	exports["throwException"] = func(err_ Any) Any {
		return func() Any {
			panic(err_)
		}
	}

	// throw :: forall a. String -> Effect a
	exports["throw"] = func(msg_ Any) Any {
		return func() Any {
			msg := msg_.(string)
			panic(Dict{
				"message": msg,
				"stack":   "",
			})
		}
	}

	// try :: forall a. Effect a -> Effect (Either Error a)
	exports["try"] = func(effect_ Any) Any {
		return func() Any {
			effect := effect_.(func() Any)
			
			// Catch panics
			defer func() {
				if r := recover(); r != nil {
					// This would be handled in the outer context
					panic(r)
				}
			}()
			
			// Try to run the effect
			result := tryRun(effect)
			return result
		}
	}

	// catchException :: forall a. (Error -> Effect a) -> Effect a -> Effect a
	exports["catchException"] = func(handler_ Any, effect_ Any) Any {
		return func() Any {
			handler := handler_.(func(Any) Any)
			effect := effect_.(func() Any)
			
			// Catch panics and handle them
			var result Any
			func() {
				defer func() {
					if r := recover(); r != nil {
						// Convert panic to error
						var err Dict
						if e, ok := r.(Dict); ok {
							err = e
						} else {
							err = Dict{
								"message": fmt.Sprintf("%v", r),
								"stack":   "",
							}
						}
						
						// Call handler
						handlerEffect := handler(err).(func() Any)
						result = handlerEffect()
					}
				}()
				
				result = effect()
			}()
			
			return result
		}
	}

	// finally :: forall a. Effect Unit -> Effect a -> Effect a
	exports["finally"] = func(finalizer_ Any, effect_ Any) Any {
		return func() Any {
			finalizer := finalizer_.(func() Any)
			effect := effect_.(func() Any)
			
			defer func() {
				finalizer()
			}()
			
			return effect()
		}
	}
}

// tryRun executes an effect and catches panics, returning Either Error a
func tryRun(effect func() Any) Any {
	var result Any
	caught := false
	var caughtErr Any
	
	func() {
		defer func() {
			if r := recover(); r != nil {
				caught = true
				// Convert panic to error dict
				if e, ok := r.(Dict); ok {
					caughtErr = e
				} else {
					caughtErr = Dict{
						"message": fmt.Sprintf("%v", r),
						"stack":   "",
					}
				}
			}
		}()
		
		result = effect()
	}()
	
	if caught {
		return Dict{"Left": caughtErr}
	}
	return Dict{"Right": result}
}
