package purescript_prelude

import . "github.com/purescript-native/go-runtime"

func init() {
	exports := Foreign("Data.Unit")
	// unit is just an empty struct - the simplest Go value
	exports["unit"] = struct{}{}
}
