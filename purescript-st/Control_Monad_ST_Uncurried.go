package purescript_st

import (
	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Control.Monad.ST.Uncurried")

	exports["runSTFn1"] = func(fn Any) Any {
		return func(a Any) Any {
			return func() Any {
				return Run(Apply(fn, a))
			}
		}
	}

	exports["runSTFn2"] = func(fn_ Any) Any {
		return func(a Any) Any {
			return func(b Any) Any {
				return func() Any {
					fn := fn_.(Fn2)
					return fn(a, b)
				}
			}
		}
	}
}
