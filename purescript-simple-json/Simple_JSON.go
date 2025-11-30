package purescript_simple_json

import (
	"encoding/json"

	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Simple.JSON")

	// parseJSON :: String -> Either String Foreign
	exports["parseJSON"] = func(str_ Any) Any {
		str := str_.(string)
		var result Any
		err := json.Unmarshal([]byte(str), &result)
		if err != nil {
			return Dict{"Left": err.Error()}
		}
		return Dict{"Right": convertFromGoJSON(result)}
	}

	// writeJSON :: Foreign -> String
	exports["writeJSON"] = func(value Any) Any {
		converted := convertToGoJSON(value)
		bytes, err := json.Marshal(converted)
		if err != nil {
			panic(err)
		}
		return string(bytes)
	}

	// writeJSON' :: Foreign -> String
	exports["writeJSON'"] = func(value Any) Any {
		converted := convertToGoJSON(value)
		bytes, err := json.MarshalIndent(converted, "", "  ")
		if err != nil {
			panic(err)
		}
		return string(bytes)
	}

	// read' :: Foreign -> Either String a
	exports["read'"] = func(foreign Any) Any {
		// Just return the foreign value wrapped in Right
		// The actual type checking happens in PureScript
		return Dict{"Right": foreign}
	}

	// write :: a -> Foreign
	exports["write"] = func(value Any) Any {
		// Just return the value as-is
		// PureScript handles the conversion
		return value
	}

	// undefined :: Foreign
	exports["undefined"] = nil

	// readString :: Foreign -> Either String String
	exports["readString"] = func(value Any) Any {
		if str, ok := value.(string); ok {
			return Dict{"Right": str}
		}
		return Dict{"Left": "Expected String"}
	}

	// readNumber :: Foreign -> Either String Number
	exports["readNumber"] = func(value Any) Any {
		switch v := value.(type) {
		case float64:
			return Dict{"Right": v}
		case int:
			return Dict{"Right": float64(v)}
		default:
			return Dict{"Left": "Expected Number"}
		}
	}

	// readInt :: Foreign -> Either String Int
	exports["readInt"] = func(value Any) Any {
		switch v := value.(type) {
		case int:
			return Dict{"Right": v}
		case float64:
			if v == float64(int(v)) {
				return Dict{"Right": int(v)}
			}
		}
		return Dict{"Left": "Expected Int"}
	}

	// readBoolean :: Foreign -> Either String Boolean
	exports["readBoolean"] = func(value Any) Any {
		if b, ok := value.(bool); ok {
			return Dict{"Right": b}
		}
		return Dict{"Left": "Expected Boolean"}
	}

	// readArray :: Foreign -> Either String (Array Foreign)
	exports["readArray"] = func(value Any) Any {
		if arr, ok := value.([]Any); ok {
			return Dict{"Right": arr}
		}
		return Dict{"Left": "Expected Array"}
	}

	// readNull :: Foreign -> Either String Unit
	exports["readNull"] = func(value Any) Any {
		if value == nil {
			return Dict{"Right": Dict{}}
		}
		return Dict{"Left": "Expected Null"}
	}
}

// convertFromGoJSON converts Go's json.Unmarshal output to PureScript-friendly format
func convertFromGoJSON(value Any) Any {
	switch v := value.(type) {
	case map[string]interface{}:
		result := make(Dict)
		for k, val := range v {
			result[k] = convertFromGoJSON(val)
		}
		return result
	case []interface{}:
		result := make([]Any, len(v))
		for i, val := range v {
			result[i] = convertFromGoJSON(val)
		}
		return result
	default:
		return v
	}
}

// convertToGoJSON converts PureScript values to Go JSON-compatible format
func convertToGoJSON(value Any) Any {
	switch v := value.(type) {
	case Dict:
		result := make(map[string]interface{})
		for k, val := range v {
			result[k] = convertToGoJSON(val)
		}
		return result
	case []Any:
		result := make([]interface{}, len(v))
		for i, val := range v {
			result[i] = convertToGoJSON(val)
		}
		return result
	default:
		return v
	}
}

