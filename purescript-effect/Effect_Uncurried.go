package purescript_effect

import (
	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Effect.Uncurried")

	exports["mkEffectFn1"] = func(fn Any) Any {
		return func(x Any) Any {
			return Run(Apply(fn, x))
		}
	}

	exports["runEffectFn1"] = func(fn Any) Any {
		return func(a Any) Any {
			return func() Any {
				return Apply(fn, a)
			}
		}
	}

	exports["runEffectFn3"] = func(fn Any) Any {
		return func(a Any) Any {
			return func(b Any) Any {
				return func(c Any) Any {
					return func() Any {
						return Apply(fn, a, b, c)
					}
				}
			}
		}
	}

}
