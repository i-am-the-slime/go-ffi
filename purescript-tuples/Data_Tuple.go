package purescript_tuples

import (
	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Data.Tuple")

	// Tuple constructor - creates a tuple from two values
	exports["Tuple"] = func(a Any) Any {
		return func(b Any) Any {
			return Dict{
				"value0": a,
				"value1": b,
			}
		}
	}

	// fst - get first element of tuple
	exports["fst"] = func(tuple Any) Any {
		t, ok := tuple.(Dict)
		if !ok {
			panic("Data.Tuple.fst: tuple is not a Dict")
		}
		return t["value0"]
	}

	// snd - get second element of tuple
	exports["snd"] = func(tuple Any) Any {
		t, ok := tuple.(Dict)
		if !ok {
			panic("Data.Tuple.snd: tuple is not a Dict")
		}
		return t["value1"]
	}

	// curry - convert a function on tuples to a curried function
	exports["curry"] = func(f Any) Any {
		return func(a Any) Any {
			return func(b Any) Any {
				tuple := Dict{
					"value0": a,
					"value1": b,
				}
				return Apply(f, tuple)
			}
		}
	}

	// uncurry - convert a curried function to a function on tuples
	exports["uncurry"] = func(f Any) Any {
		return func(tuple Any) Any {
			t, ok := tuple.(Dict)
			if !ok {
				panic("Data.Tuple.uncurry: tuple is not a Dict")
			}
			return Apply(Apply(f, t["value0"]), t["value1"])
		}
	}
}
