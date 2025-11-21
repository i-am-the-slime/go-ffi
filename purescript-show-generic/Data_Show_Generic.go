package purescript_show_generic

import (
	"strings"

	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Data.Show.Generic")

	// intercalate takes a separator and an array of strings and joins them
	// intercalate :: String -> Array String -> String
	exports["intercalate"] = func(sep Any) Any {
		return func(xs_ Any) Any {
			xs := xs_.([]Any)
			ss := make([]string, 0, len(xs))
			for _, x := range xs {
				ss = append(ss, x.(string))
			}
			return strings.Join(ss, sep.(string))
		}
	}
}
