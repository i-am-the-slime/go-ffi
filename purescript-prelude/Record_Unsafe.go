package purescript_prelude

import (
	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Record.Unsafe")

	exports["unsafeHas"] = func(label_ Any) Any {
		return func(rec_ Any) Any {
			rec := rec_.(Dict)
			label := label_.(string)
			_, ok := rec[label]
			return ok
		}
	}

	exports["unsafeGet"] = func(label_ Any) Any {
		return func(rec_ Any) Any {
			rec := rec_.(map[string]Any)
			label := label_.(string)
			return rec[label]
		}
	}

	exports["unsafeSet"] = func(label_ Any) Any {
		return func(value Any) Any {
			return func(rec_ Any) Any {
				label, _ := label_.(string)
				rec, _ := rec_.(Dict)
				copy := make(Dict)
				for key, val := range rec {
					copy[key] = val
				}
				copy[label] = value
				return copy
			}
		}
	}

}
