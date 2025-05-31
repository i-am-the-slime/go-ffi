package purescript_record

import (
	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Record.Builder")

	exports["copyRecord"] = func(rec Any) Any {
		rc, ok := rec.(Dict)
		if !ok {
			panic("Record.Builder.copyRecord: record is not a Dict")
		}
		cpy := make(Dict)
		for key, value := range rc {
			cpy[key] = value
		}
		return cpy
	}

	exports["unsafeInsert"] = func(l_ Any) Any {
		return func(a Any) Any {
			return func(rec_ Any) Any {
				l, ok := l_.(string)
				if !ok {
					panic("Record.Builder.unsafeInsert: label is not a String")
				}
				rec, ok := rec_.(Dict)
				if !ok {
					panic("Record.Builder.unsafeInsert: record is not a Dict")
				}
				rec[l] = a
				return rec
			}
		}
	}
}
