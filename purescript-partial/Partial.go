package purescript_partial

import . "github.com/purescript-native/go-runtime"

func init() {
	exports := Foreign("Partial")

	exports["crashWith"] = func(dict Any) Any {
		return func(msg Any) Any {
			panic(msg)
		}
	}

}
