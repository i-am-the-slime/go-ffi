package purescript_prelude

import (
	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Record.Unsafe")

	exports["unsafeHas"] = func(label_ Any) Any {
		return func(rec_ Any) Any {
			rec, ok := rec_.(Dict)
			if !ok {
				panic("Record.Unsafe.unsafeHas: record is not a Dict")
			}
			label, ok := label_.(string)
			if !ok {
				panic("Record.Unsafe.unsafeHas: label is not a String")
			}
			_, ok = rec[label]
			return ok
		}
	}

	exports["unsafeGet"] = func(label_ Any) Any {
		return func(rec_ Any) Any {
			rec, ok := rec_.(map[string]Any)
			if !ok {
				panic("Record.Unsafe.unsafeGet: record is not a Dict")
			}
			label, ok := label_.(string)
			if !ok {
				panic("Record.Unsafe.unsafeGet: label is not a String")
			}
			return rec[label]
		}
	}

	exports["unsafeSet"] = func(label Any) Any {
		return func(value Any) Any {
			return func(rec Any) Any {
				rc, ok := rec.(Dict)
				if !ok {
					panic("Record.Unsafe.unsafeSet: record is not a Dict")
				}
				copy := make(Dict)
				for key, val := range rc {
					copy[key] = val
				}
				copy[label.(string)] = value
				return copy
			}
		}
	}

}
