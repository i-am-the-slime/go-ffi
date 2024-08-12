package purescript_arrays

import (
	"sort"

	. "github.com/purescript-native/go-runtime"
)

func init() {

	exports := Foreign("Data.Array")

	//------------------------------------------------------------------------------
	// Array creation --------------------------------------------------------------
	//------------------------------------------------------------------------------

	exports["range"] = func(start_ Any) Any {
		return func(end_ Any) Any {
			start := start_.(int)
			end := end_.(int)
			var length int
			var step int
			if start > end {
				length = start - end + 1
				step = -1
			} else {
				length = end - start + 1
				step = 1
			}
			ns := make([]Any, 0, length)
			for i := start; i != end; i += step {
				ns = append(ns, i)
			}
			return append(ns, end)
		}
	}

	exports["replicate"] = func(count_ Any, value Any) Any {
		var count = count_.(int)
		if count < 1 {
			return []Any{}
		}
		var arr = make([]Any, count)
		for i := range arr {
			arr[i] = value
		}
		return arr
	}

	exports["fromFoldableImpl"] = func(foldr Any, foldable Any) Any {
		var f = func(x Any) Any {
			return func(acc Any) Any {
				return append(acc.([]Any), x)
			}
		}
		var xs = Apply(foldr, f, []Any{}, foldable)
		return exports["reverse"].(Fn)(xs)
	}

	//------------------------------------------------------------------------------
	// Array size ------------------------------------------------------------------
	//------------------------------------------------------------------------------

	exports["length"] = func(xs Any) Any {
		return len(xs.([]Any))
	}

	//------------------------------------------------------------------------------
	// Extending arrays ------------------------------------------------------------
	//------------------------------------------------------------------------------

	exports["cons"] = func(e Any, l Any) Any {
		return append([]Any{e}, l.([]Any)...)
	}

	exports["snoc"] = func(l_ Any, e Any) Any {
		l := l_.([]Any)
		xs := make([]Any, len(l), len(l)+1)
		copy(xs, l)
		return append(xs, e)
	}

	//------------------------------------------------------------------------------
	// Non-indexed reads -----------------------------------------------------------
	//------------------------------------------------------------------------------

	exports["unconsImpl"] = func(empty_ Any, next Any, xs_ Any) Any {
		empty, xs := empty_.(Fn), xs_.([]Any)
		if len(xs) == 0 {
			return empty(Dict{})
		}
		return Apply(next, xs[0], xs[1:])
	}

	//------------------------------------------------------------------------------
	// Indexed operations ----------------------------------------------------------
	//------------------------------------------------------------------------------

	exports["indexImpl"] = func(just Any, nothing Any, xs_ Any, i_ Any) Any {
		xs := xs_.([]Any)
		i := i_.(int)
		if i < 0 || i >= len(xs) {
			return nothing
		}
		return Apply(just, xs[i])
	}

	exports["findIndexImpl"] = func(just_ Any, nothing Any, f_ Any, xs_ Any) Any {
		xs, f, just := xs_.([]Any), f_.(Fn), just_.(Fn)
		for i, x := range xs {
			if f(x).(bool) {
				return just(i)
			}
		}
		return nothing
	}

	exports["findLastIndexImpl"] = func(just_ Any, nothing Any, f_ Any, xs_ Any) Any {
		xs, f, just := xs_.([]Any), f_.(Fn), just_.(Fn)
		for i := len(xs) - 1; i >= 0; i-- {
			if f(xs[i]).(bool) {
				return just(i)
			}
		}
		return nothing
	}

	exports["_insertAt"] = func(just_ Any, nothing Any, i_ Any, a Any, xs_ Any) Any {
		just, xs, i := just_.(Fn), xs_.([]Any), i_.(int)
		if i < 0 || i > len(xs) {
			return nothing
		}
		ys := make([]Any, len(xs)+1)
		copy(ys, xs[:i])
		ys[i] = a
		copy(ys[i+1:], xs[i:])
		return just(ys)
	}

	exports["_deleteAt"] = func(just_ Any, nothing Any, i_ Any, xs_ Any) Any {
		just, i, xs := just_.(Fn), i_.(int), xs_.([]Any)
		if i < 0 || i >= len(xs) {
			return nothing
		}
		ys := make([]Any, len(xs)-1)
		copy(ys, xs[:i])
		copy(ys[i:], xs[i+1:])
		return just(ys)
	}

	exports["_updateAt"] = func(just_ Any, nothing Any, i_ Any, a Any, xs_ Any) Any {
		just, i, xs := just_.(Fn), i_.(int), xs_.([]Any)
		if i < 0 || i >= len(xs) {
			return nothing
		}
		ys := make([]Any, len(xs))
		copy(ys, xs)
		ys[i] = a
		return just(ys)
	}

	//------------------------------------------------------------------------------
	// Transformations -------------------------------------------------------------
	//------------------------------------------------------------------------------

	exports["reverse"] = func(xs_ Any) Any {
		xs := xs_.([]Any)
		l := len(xs)
		ys := make([]Any, l)
		for i, j := 0, l-1; i < l; i, j = i+1, j-1 {
			ys[i] = xs[j]
		}
		return ys
	}

	exports["concat"] = func(xss_ Any) Any {
		xss := xss_.([]Any)
		result := []Any{}
		for _, xs := range xss {
			result = append(result, xs.([]Any)...)
		}
		return result
	}

	exports["filterImpl"] = func(f_ Any, xs_ Any) Any {
		xs, f := xs_.([]Any), f_.(Fn)
		result := []Any{}
		for _, x := range xs {
			if f(x).(bool) {
				result = append(result, x)
			}
		}
		return result
	}

	exports["partition"] = func(f_ Any, xs_ Any) Any {
		xs, f := xs_.([]Any), f_.(Fn)
		result := Dict{"yes": []Any{}, "no": []Any{}}
		for _, x := range xs {
			if f(x).(bool) {
				result["yes"] = append(result["yes"].([]Any), x)
			} else {
				result["no"] = append(result["no"].([]Any), x)
			}
		}
		return result
	}

	//------------------------------------------------------------------------------
	// Sorting ---------------------------------------------------------------------
	//------------------------------------------------------------------------------

	exports["sortImpl"] = func(f Any) Any {
		return func(l_ Any) Any {
			l := l_.([]Any)
			xs := make([]Any, len(l))
			copy(xs, l)
			sort.SliceStable(xs, func(i int, j int) bool {
				return Apply(f, xs[i], xs[j]).(int) < 0
			})
			return xs
		}
	}

	// foreign import sortByImpl :: forall a. Fn3 (a -> a -> Ordering) (Ordering -> Int) (Array a) (Array a)
	exports["sortByImpl"] = func(orderingFn_, orderingToInt_, l_ Any) Any {
		orderingFn := orderingFn_.(Fn)
		orderingToInt := orderingToInt_.(Fn)
		l := l_.([]Any)
		xs := make([]Any, len(l))
		copy(xs, l)
		sort.SliceStable(xs, func(i int, j int) bool {
			return Apply(orderingToInt(Apply(orderingFn, xs[i], xs[j]))).(int) < 0
		})
		return xs
	}

	//------------------------------------------------------------------------------
	// Subarrays -------------------------------------------------------------------
	//------------------------------------------------------------------------------

	exports["sliceImpl"] = func(s_ Any, e_ Any, l_ Any) Any {
		s := s_.(int)
		e := e_.(int)
		l := l_.([]Any)
		sz := len(l)
		// Adjust negative indices
		if s < 0 {
			s = sz + s
		}
		if e < 0 {
			e = sz + e
		}

		// Ensure indices are within bounds
		if s < 0 {
			s = 0
		}
		if e > sz {
			e = sz
		}
		if s > e {
			s = e
		}

		return l[s:e]
	}

	exports["take"] = func(n_ Any, l_ Any) Any {
		n, l := n_.(int), l_.([]Any)
		if n < 1 {
			return []Any{}
		}
		if n > len(l) {
			return l
		}
		return l[:n]
	}

	//------------------------------------------------------------------------------
	// Zipping ---------------------------------------------------------------------
	//------------------------------------------------------------------------------

	exports["zipWith"] = func(f_ Any, xs_ Any, ys_ Any) Any {
		f := f_.(Fn)
		xs := xs_.([]Any)
		ys := ys_.([]Any)
		lxs := len(xs)
		l := len(ys)
		if lxs < l {
			l = lxs
		}
		result := make([]Any, 0, l)
		for i := 0; i < l; i++ {
			fx := f(xs[i]).(Fn)
			result = append(result, fx(ys[i]))
		}
		return result
	}

	//------------------------------------------------------------------------------
	// Partial ---------------------------------------------------------------------
	//------------------------------------------------------------------------------

	exports["unsafeIndexImpl"] = func(xs_ Any, n_ Any) Any {
		xs, n := xs_.([]Any), n_.(int)
		return xs[n]
	}

}
