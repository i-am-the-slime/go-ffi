package purescript_integers

import (
	"fmt"
	"math"
	"strconv"

	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Data.Int")

	exports["toNumber"] = func(n Any) Any {
		return float64(n.(int))
	}

	exports["fromNumberImpl"] = func(just_ Any) Any {
		return func(nothing Any) Any {
			return func(n_ Any) Any {
				just, _ := just_.(Fn)
				n := n_.(float64)
				if math.Round(n) == n {
					return just(int(n))
				} else {
					return nothing
				}
			}
		}
	}

	exports["pow"] = func(n_ Any) Any {
		return func(p_ Any) Any {
			n, _ := n_.(int)
			p, _ := p_.(int)
			return int(math.Pow(float64(n), float64(p)))
		}
	}

	exports["toStringAs"] = func(radix_ Any) Any {
		return func(i_ Any) Any {
			radix, _ := radix_.(int)
			i, _ := i_.(int)
			switch radix {
			case 2:
				return fmt.Sprintf("%b", i)
			case 8:
				return fmt.Sprintf("%o", i)
			case 10:
				return fmt.Sprintf("%d", i)
			case 16:
				return fmt.Sprintf("%x", i)
			default:
				panic(fmt.Sprintf("Unsupported radix (%d)", radix))
			}
		}
	}

	exports["fromStringAsImpl"] = func(just Any) Any {
		return func(nothing Any) Any {
			return func(radix_ Any) Any {
				return func(s_ Any) Any {
					radix := radix_.(int)
					s := s_.(string)
					res, err := strconv.ParseInt(s, 0, radix)
					if err != nil {
						return nothing
					}
					return Apply(just, res)
				}
			}
		}
	}

}
