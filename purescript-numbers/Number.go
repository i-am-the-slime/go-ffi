package purescript_math

import (
	"math"
	"strconv"

	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Number")

	exports["nan"] = func() Any {
		return math.NaN()
	}

	exports["infinity"] = func() Any {
		return math.Inf(1)
	}

	exports["isNaN"] = func(x Any) Any {
		return math.IsNaN(x.(float64))
	}

	exports["isFinite"] = func(x Any) Any {
		return !math.IsInf(x.(float64), 0) && !math.IsNaN(x.(float64))
	}

	exports["fromStringImpl"] = func(str Any) Any {
		return func(isFinite Any) Any {
			return func(just Any) Any {
				return func(nothing Any) Any {
					num, err := strconv.ParseFloat(str.(string), 64)
					if err == nil && isFinite.(func(Any) Any)(num).(bool) {
						return just.(func(Any) Any)(num)
					}
					return nothing
				}
			}
		}
	}

	exports["abs"] = func(x Any) Any {
		return math.Abs(x.(float64))
	}

	exports["ceil"] = func(x Any) Any {
		return math.Ceil(x.(float64))
	}

	exports["floor"] = func(x Any) Any {
		return math.Floor(x.(float64))
	}

	exports["pow"] = func(n Any) Any {
		return func(p Any) Any {
			return math.Pow(n.(float64), p.(float64))
		}
	}

	exports["remainder"] = func(n Any) Any {
		return func(m Any) Any {
			return math.Remainder(n.(float64), m.(float64))
		}
	}

	exports["round"] = func(x Any) Any {
		return math.Round(x.(float64))
	}

	exports["acos"] = func(x Any) Any {
		return math.Acos(x.(float64))
	}

	exports["asin"] = func(x Any) Any {
		return math.Asin(x.(float64))
	}

	exports["atan"] = func(x Any) Any {
		return math.Atan(x.(float64))
	}

	exports["atan2"] = func(y Any) Any {
		return func(x Any) Any {
			return math.Atan2(y.(float64), x.(float64))
		}
	}

	exports["cos"] = func(x Any) Any {
		return math.Cos(x.(float64))
	}

	exports["exp"] = func(x Any) Any {
		return math.Exp(x.(float64))
	}

	exports["log"] = func(x Any) Any {
		return math.Log(x.(float64))
	}

	exports["max"] = func(n1 Any) Any {
		return func(n2 Any) Any {
			return math.Max(n1.(float64), n2.(float64))
		}
	}

	exports["min"] = func(n1 Any) Any {
		return func(n2 Any) Any {
			return math.Min(n1.(float64), n2.(float64))
		}
	}

	exports["sign"] = func(x Any) Any {
		if x.(float64) == 0 || math.IsNaN(x.(float64)) {
			return x
		}
		if x.(float64) < 0 {
			return -1
		}
		return 1
	}

	exports["sin"] = func(x Any) Any {
		return math.Sin(x.(float64))
	}

	exports["sqrt"] = func(x Any) Any {
		return math.Sqrt(x.(float64))
	}

	exports["tan"] = func(x Any) Any {
		return math.Tan(x.(float64))
	}

	exports["trunc"] = func(x Any) Any {
		if x.(float64) < 0 {
			return math.Ceil(x.(float64))
		}
		return math.Floor(x.(float64))
	}
}
