package purescript_arrays

import (
	. "github.com/purescript-native/go-runtime"
)

// STArray is a pointer to a slice, allowing in-place mutation including growth
type STArray = *[]Any

func init() {

	exports := Foreign("Data.Array.ST")

	exports["new"] = func() Any {
		result := make([]Any, 0)
		return STArray(&result)
	}

	exports["peekImpl"] = func(just_ Any) Any {
		return func(nothing Any) Any {
			return func(i_ Any) Any {
				return func(xs_ Any) Any {
					return func() Any {
						just, i := just_.(Fn), i_.(int)
						xs := *xs_.(STArray)
						if i >= 0 && i < len(xs) {
							return just((xs)[i])
						}
						return nothing
					}
				}
			}
		}
	}

	// foreign import thawImpl :: forall h a. STFn1 (Array a) h (STArray h a)
	exports["thawImpl"] = func(xs_ Any) Any {
		xs := xs_.([]Any)
		result := make([]Any, len(xs))
		copy(result, xs)
		return STArray(&result)
	}

	exports["poke"] = func(i_ Any) Any {
		return func(a Any) Any {
			return func(xs_ Any) Any {
				return func() Any {
					i := i_.(int)
					xs := xs_.(STArray)
					result := i >= 0 && i < len(*xs)
					if result {
						(*xs)[i] = a
					}
					return result
				}
			}
		}
	}

	exports["pushAll"] = func(as_ Any) Any {
		return func(xs_ Any) Any {
			return func() Any {
				as := as_.([]Any)
				xs := xs_.(STArray)
				*xs = append(*xs, as...)
				return len(*xs)
			}
		}
	}

	exports["pushImpl"] = func(a Any, xs_ Any) Any {
		xs := xs_.(STArray)
		*xs = append(*xs, a)
		return len(*xs)
	}

	exports["unsafeFreeze"] = func(xs_ Any) Any {
		return func() Any {
			xs := xs_.(STArray)
			return *xs
		}
	}

	exports["unsafeFreezeImpl"] = func(xs_ Any) Any {
		xs := xs_.(STArray)
		return *xs
	}

	exports["unsafeThaw"] = func(xs_ Any) Any {
		return func() Any {
			xs := xs_.([]Any)
			return STArray(&xs)
		}
	}

	exports["unsafeThawImpl"] = func(xs_ Any) Any {
		xs := xs_.([]Any)
		return STArray(&xs)
	}

	exports["freeze"] = func(xs_ Any) Any {
		return func() Any {
			xs := *xs_.(STArray)
			result := make([]Any, len(xs))
			copy(result, xs)
			return result
		}
	}

	exports["thaw"] = func(xs_ Any) Any {
		return func() Any {
			xs := xs_.([]Any)
			result := make([]Any, len(xs))
			copy(result, xs)
			return STArray(&result)
		}
	}

	exports["any"] = func(xs_ Any) Any {
		return func() Any {
			xs := xs_.([]Any)
			result := make([]Any, len(xs))
			copy(result, xs)
			return STArray(&result)
		}
	}

}
