package purescript_foreign_object

import (
	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Foreign.Object")

	exports["_copyST"] = func(m Any) Any {
		return func() Any {
			r := make(Dict)
			for k, v := range m.(Dict) {
				r[k] = v
			}
			return r
		}
	}

	exports["empty"] = func() Any {
		return make(Dict)
	}

	exports["runST"] = func(f Any) Any {
		return f.(func() Any)()
	}

	exports["_fmapObject"] = func(m0 Any, f Any) Any {
		m := make(Dict)
		for k, v := range m0.(Dict) {
			m[k] = f.(func(Any) Any)(v)
		}
		return m
	}

	exports["_mapWithKey"] = func(m0 Any, f Any) Any {
		m := make(Dict)
		for k, v := range m0.(Dict) {
			m[k] = f.(func(Any) func(Any) Any)(k)(v)
		}
		return m
	}

	exports["_foldM"] = func(bind Any) Any {
		return func(f Any) Any {
			return func(mz Any) Any {
				return func(m Any) Any {
					acc := mz
					for k, v := range m.(Dict) {
						acc = bind.(func(Any) func(func() Any) Any)(acc)(func() Any {
							return f.(func(Any) func(Any) func(Any) Any)(acc)(k)(v)
						})
					}
					return acc
				}
			}
		}
	}

	exports["_foldSCObject"] = func(m Any, z Any, f Any, fromMaybe Any) Any {
		acc := z
		for k, v := range m.(Dict) {
			maybeR := f.(func(Any) func(Any) func(Any) Any)(acc)(k)(v)
			r := fromMaybe.(func(Any) func(Any) Any)(nil)(maybeR)
			if r == nil {
				return acc
			}
			acc = r
		}
		return acc
	}

	exports["all"] = func(f Any) Any {
		return func(m Any) Any {
			for k, v := range m.(Dict) {
				if !f.(func(Any) func(Any) bool)(k)(v) {
					return false
				}
			}
			return true
		}
	}

	exports["size"] = func(m Any) Any {
		return len(m.(Dict))
	}

	exports["_lookup"] = func(no Any, yes Any, k Any, m Any) Any {
		if v, ok := m.(Dict)[k.(string)]; ok {
			return yes.(func(Any) Any)(v)
		}
		return no
	}

	exports["_lookupST"] = func(no Any, yes Any, k Any, m Any) Any {
		return func() Any {
			if v, ok := m.(Dict)[k.(string)]; ok {
				return yes.(func(Any) Any)(v)
			}
			return no
		}
	}

	exports["toArrayWithKey"] = func(f Any) Any {
		return func(m Any) Any {
			r := make([]Any, 0, len(m.(Dict)))
			for k, v := range m.(Dict) {
				r = append(r, Apply(f, k, v))
			}
			return r
		}
	}

	exports["keys"] = func(m Any) Any {
		keys := make([]Any, 0, len(m.(Dict)))
		for k := range m.(Dict) {
			keys = append(keys, k)
		}
		return keys
	}
}
