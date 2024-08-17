package purescript_st

import (
	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Control.Monad.ST.Uncurried")

	// export const runSTFn1 = function runSTFn1(fn) {
	//   return function(a) {
	//     return function() {
	//       return fn(a);
	//     };
	//   };
	// };
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

	exports["runSTFn3"] = func(fn_ Any) Any {
		return func(a Any) Any {
			return func(b Any) Any {
				return func(c Any) Any {
					return func() Any {
						fn := fn_.(Fn3)
						return fn(a, b, c)
					}
				}
			}
		}
	}

	exports["runSTFn4"] = func(fn_ Any) Any {
		return func(a Any) Any {
			return func(b Any) Any {
				return func(c Any) Any {
					return func(d Any) Any {
						return func() Any {
							fn := fn_.(Fn4)
							return fn(a, b, c, d)
						}
					}
				}
			}
		}
	}

	exports["runSTFn5"] = func(fn_ Any) Any {
		return func(a Any) Any {
			return func(b Any) Any {
				return func(c Any) Any {
					return func(d Any) Any {
						return func(e Any) Any {
							return func() Any {
								fn := fn_.(Fn5)
								return fn(a, b, c, d, e)
							}
						}
					}
				}
			}
		}
	}

}
