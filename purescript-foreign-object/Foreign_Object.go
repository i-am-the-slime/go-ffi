package purescript_foreign_object

import (
	"fmt"

	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Foreign.Object")

	exports["_copyST"] = func(m Any) Any {
		return func() Any {
			if fn, ok := m.(func() Any); ok {
				m = fn()
			}
			if m == nil {
				return make(Dict)
			}
			switch src := m.(type) {
			case Dict:
				r := make(Dict)
				for k, v := range src {
					r[k] = v
				}
				return r
			default:
				panic(fmt.Sprintf("Foreign.Object._copyST: input is not a Dict (received %T)", m))
			}
		}
	}

	exports["empty"] = func() Any {
		return make(Dict)
	}

	exports["runST"] = func(f Any) Any {
		fn, ok := f.(func() Any)
		if !ok {
			panic("Foreign.Object.runST: argument is not a function of type func() Any")
		}
		return fn()
	}

	exports["_fmapObject"] = func(m0 Any, f Any) Any {
		src, ok := m0.(Dict)
		if !ok {
			panic("Foreign.Object._fmapObject: first argument is not a Dict")
		}
		m := make(Dict)
		for k, v := range src {
			m[k] = f.(func(Any) Any)(v)
		}
		return m
	}

	exports["_mapWithKey"] = func(m0 Any, f Any) Any {
		src, ok := m0.(Dict)
		if !ok {
			panic("Foreign.Object._mapWithKey: first argument is not a Dict")
		}
		m := make(Dict)
		for k, v := range src {
			m[k] = Apply(f, k, v)
		}
		return m
	}

	exports["_foldM"] = func(bind Any) Any {
		return func(f Any) Any {
			return func(mz Any) Any {
				return func(m Any) Any {
					acc := mz
					md, ok := m.(Dict)
					if !ok {
						panic("Foreign.Object._foldM: target is not a Dict")
					}
					for k, v := range md {
						currentAcc := acc
						step := func() Any {
							return Apply(f, currentAcc, k, v)
						}
						acc = Apply(bind, currentAcc, step)
					}
					return acc
				}
			}
		}
	}

	exports["_foldSCObject"] = func(m Any, z Any, f Any, fromMaybe Any) Any {
		acc := z
		md, ok := m.(Dict)
		if !ok {
			panic("Foreign.Object._foldSCObject: target is not a Dict")
		}
		for k, v := range md {
			maybeR := Apply(f, acc, k, v)
			r := Apply(fromMaybe, nil, maybeR)
			if r == nil {
				return acc
			}
			acc = r
		}
		return acc
	}

	exports["all"] = func(f Any) Any {
		return func(m Any) Any {
			md, ok := m.(Dict)
			if !ok {
				panic("Foreign.Object.all: target is not a Dict")
			}
			for k, v := range md {
				if !Apply(f, k, v).(bool) {
					return false
				}
			}
			return true
		}
	}

	exports["size"] = func(m Any) Any {
		md, ok := m.(Dict)
		if !ok {
			panic("Foreign.Object.size: argument is not a Dict")
		}
		return len(md)
	}

	exports["_lookup"] = func(no Any, yes Any, k Any, m Any) Any {
		md, ok := m.(Dict)
		if !ok {
			panic("Foreign.Object._lookup: input is not a Dict")
		}
		ks, ok := k.(string)
		if !ok {
			panic("Foreign.Object._lookup: key is not a String")
		}
		if v, ok := md[ks]; ok {
			return yes.(func(Any) Any)(v)
		}
		return no
	}

	exports["_lookupST"] = func(no Any, yes Any, k Any, m Any) Any {
		return func() Any {
			md, ok := m.(Dict)
			if !ok {
				panic("Foreign.Object._lookupST: input is not a Dict")
			}
			ks, ok := k.(string)
			if !ok {
				panic("Foreign.Object._lookupST: key is not a String")
			}
			if v, ok := md[ks]; ok {
				return yes.(func(Any) Any)(v)
			}
			return no
		}
	}

	exports["toArrayWithKey"] = func(f_ Any) Any {
		fmt.Println("toArrayWithKey", f_)
		return func(m_ Any) Any {
			mDict, ok := m_.(Dict)
			if !ok {
				panic("Foreign.Object.toArrayWithKey: input is not a Dict")
			}
			r := make([]Any, 0, len(mDict))
			for k, v := range mDict {
				r = append(r, Apply(f_, k, v))
			}
			return r
		}
	}

	exports["keys"] = func(m Any) Any {
		md, ok := m.(Dict)
		if !ok {
			panic("Foreign.Object.keys: argument is not a Dict")
		}
		keys := make([]Any, 0, len(md))
		for k := range md {
			keys = append(keys, k)
		}
		return keys
	}
}
