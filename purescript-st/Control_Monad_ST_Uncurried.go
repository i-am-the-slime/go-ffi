package purescript_st

import (
	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Control.Monad.ST.Uncurried")

	exports["runSTFn1"] = func(fn Any) Any {
		return func(a Any) Any {
			return func() Any {
				Run(Apply(fn, a))
			}
		}
	}
}
