package purescript_partial

import . "github.com/purescript-native/go-runtime"

func init() {
	exports := Foreign("Partial")

	exports["_crashWith"] = func(msg Any) Any {
		panic(msg)
	}

}
