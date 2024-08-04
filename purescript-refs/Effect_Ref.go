package purescript_refs

import . "github.com/purescript-native/go-runtime"

func init() {
	exports := Foreign("Effect.Ref")

	exports["_new"] = func(val Any) Any {
		return func() Any {
			ptr := new(Any)
			*ptr = val
			return ptr
		}
	}

	exports["newWithSelf"] = func(f Any) Any {
		return func() Any {
			ptr := new(Any)
			*ptr = nil
			ptr = Apply(f, ptr).(*Any)
			return ptr
		}
	}

	exports["read"] = func(ref_ Any) Any {
		return func() Any {
			ref := ref_.(*Any)
			return *ref
		}
	}

	exports["modifyImpl"] = func(f Any) Any {
		return func(ref_ Any) Any {
			return func() Any {
				ref := ref_.(*Any)
				t := Apply(f, *ref).(Dict)
				*ref = t["state"]
				return t["value"]
			}
		}
	}

	exports["write"] = func(a Any) Any {
		return func(ref_ Any) Any {
			return func() Any {
				ref := ref_.(*Any)
				*ref = a
				return nil
			}
		}
	}

}
