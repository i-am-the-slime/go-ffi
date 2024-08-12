package purescript_strings

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/dlclark/regexp2"
	. "github.com/purescript-native/go-runtime"
)

type regex_pair struct {
	regex  *regexp2.Regexp
	global bool
}

func regexp2FindAllString(re *regexp2.Regexp, s string) []string {
	var matches []string
	m, _ := re.FindStringMatch(s)
	for m != nil {
		matches = append(matches, m.String())
		m, _ = re.FindNextMatch(m)
	}
	return matches
}

func init() {
	exports := Foreign("Data.String.Regex")

	exports["regexImpl"] = func(left Any) Any {
		return func(right Any) Any {
			return func(s_ Any) Any {
				return func(flags_ Any) Any {
					s := s_.(string)
					flags := flags_.(string)
					global := false
					if strings.Contains(flags, "g") {
						global = true
						flags = strings.ReplaceAll(flags, "g", "")
						if flags != "" {
							flags = fmt.Sprintf("(?%s)", flags)
						}
					}
					r, err := regexp2.Compile(flags+s, 0)
					if err == nil {
						return Apply(right, regex_pair{r, global})
					} else {
						return Apply(left, err.Error())
					}
				}
			}
		}
	}

	exports["test"] = func(p_ Any) Any {
		return func(s_ Any) Any {
			p := p_.(regex_pair)
			r := p.regex
			s := s_.(string)
			_, result := r.MatchString(s)
			return result
		}
	}

	exports["_match"] = func(just Any) Any {
		return func(nothing Any) Any {
			return func(p_ Any) Any {
				return func(s_ Any) Any {
					p := p_.(regex_pair)
					r := p.regex
					s := s_.(string)
					ms := regexp2FindAllString(r, s)
					if ms == nil {
						return nothing
					}
					result := make([]Any, 0, len(ms))
					for _, m := range ms {
						if m == "" {
							result = append(result, nothing)
						} else {
							result = append(result, Apply(just, m))
						}
					}
					return Apply(just, result)
				}
			}
		}
	}

	exports["replace"] = func(p_ Any) Any {
		return func(s1_ Any) Any {
			return func(s2_ Any) Any {
				p := p_.(regex_pair)
				r := p.regex
				global := p.global

				s1 := s1_.(string)
				s2 := s2_.(string)

				if global {
					result, _ := r.Replace(s2, s1, -1, -1)
					return result
				}

				result, _ := r.Replace(s2, s1, -1, 1)
				return result
			}
		}
	}

	exports["replaceBy"] = func(p_ Any) Any {
		return func(f Any) Any {
			return func(s_ Any) Any {
				p := p_.(regex_pair)
				r := p.regex
				global := p.global
				s := s_.(string)

				all := regexp2FindAllString(r, s)
				submatches := make([]Any, 0, len(all))
				for _, submatch := range all {
					submatches = append(submatches, submatch)
				}

				frepl := func(str regexp2.Match) string {
					return Apply(f, str.String(), submatches).(string)
				}

				if global {
					result, _ := r.ReplaceFunc(s, frepl, -1, -1)
					return result
				}

				result, _ := r.ReplaceFunc(s, frepl, -1, 1)
				return result
			}
		}
	}

	exports["_search"] = func(just Any) Any {
		return func(nothing Any) Any {
			return func(p_ Any) Any {
				return func(s_ Any) Any {
					p := p_.(regex_pair)
					r := p.regex
					s := s_.(string)
					foundMatch, _ := r.FindStringMatch(s)
					if foundMatch == nil {
						return nothing
					}
					// TODO: is there a way to do this that is faster?
					return Apply(just, utf8.RuneCountInString(s[:foundMatch.String()[0]]))
				}
			}
		}
	}

	exports["split"] = func(p_ Any) Any {
		return func(s_ Any) Any {
			p := p_.(regex_pair)
			r := p.regex
			s := s_.(string)
			matches, _ := r.FindStringMatch(s)
			var result []Any
			lastIndex := 0
			for matches != nil {
				result = append(result, s[lastIndex:matches.Index])
				lastIndex = matches.Index + matches.Length
				matches, _ = r.FindNextMatch(matches)
			}
			result = append(result, s[lastIndex:])
			return result
		}
	}

}
