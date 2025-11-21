package purescript_functions

import . "github.com/purescript-native/go-runtime"

func init() {
	exports := Foreign("Data.Function.Uncurried")

	exports["mkFn2"] = func(fn Any) Any {
		return func(a Any, b Any) Any {
			return Apply(fn, a, b)
		}
	}

	exports["mkFn3"] = func(fn Any) Any {
		return func(a Any, b Any, c Any) Any {
			return Apply(fn, a, b, c)
		}
	}

	exports["mkFn4"] = func(fn Any) Any {
		return func(a Any, b Any, c Any, d Any) Any {
			return Apply(fn, a, b, c, d)
		}
	}

	exports["mkFn5"] = func(fn Any) Any {
		return func(a Any, b Any, c Any, d Any, e Any) Any {
			return Apply(fn, a, b, c, d, e)
		}
	}

	exports["mkFn6"] = func(fn Any) Any {
		return func(a Any, b Any, c Any, d Any, e Any, f Any) Any {
			return Apply(fn, a, b, c, d, e, f)
		}
	}

	exports["mkFn7"] = func(fn Any) Any {
		return func(a Any, b Any, c Any, d Any, e Any, f Any, g Any) Any {
			return Apply(fn, a, b, c, d, e, f, g)
		}
	}

	exports["mkFn8"] = func(fn Any) Any {
		return func(a Any, b Any, c Any, d Any, e Any, f Any, g Any, h Any) Any {
			return Apply(fn, a, b, c, d, e, f, g, h)
		}
	}

	exports["runFn2"] = func(fn Any) Any {
		return func(a Any) Any {
			return func(b Any) Any {
				f := fn.(Fn2)
				return f(a, b)
			}
		}
	}

	exports["runFn3"] = func(fn Any) Any {
		return func(a Any) Any {
			return func(b Any) Any {
				return func(c Any) Any {
					f := fn.(Fn3)
					return f(a, b, c)
				}
			}
		}
	}

	exports["runFn4"] = func(fn Any) Any {
		return func(a Any) Any {
			return func(b Any) Any {
				return func(c Any) Any {
					return func(d Any) Any {
						f := fn.(Fn4)
						return f(a, b, c, d)
					}
				}
			}
		}
	}

	exports["runFn5"] = func(fn Any) Any {
		return func(a Any) Any {
			return func(b Any) Any {
				return func(c Any) Any {
					return func(d Any) Any {
						return func(e Any) Any {
							f := fn.(Fn5)
							return f(a, b, c, d, e)
						}
					}
				}
			}
		}
	}

	exports["runFn6"] = func(fn Any) Any {
		return func(a Any) Any {

			return func(b Any) Any {
				return func(c Any) Any {
					return func(d Any) Any {
						return func(e Any) Any {
							return func(f Any) Any {
								g := fn.(Fn6)
								return g(a, b, c, d, e, f)
							}
						}
					}
				}
			}
		}
	}

	exports["runFn7"] = func(fn Any) Any {
		return func(a Any) Any {
			return func(b Any) Any {
				return func(c Any) Any {
					return func(d Any) Any {
						return func(e Any) Any {
							return func(f Any) Any {
								return func(g Any) Any {
									h := fn.(Fn7)
									return h(a, b, c, d, e, f, g)
								}
							}
						}
					}
				}
			}
		}
	}

	exports["runFn8"] = func(fn Any) Any {
		return func(a Any) Any {
			return func(b Any) Any {
				return func(c Any) Any {
					return func(d Any) Any {
						return func(e Any) Any {
							return func(f Any) Any {
								return func(g Any) Any {
									return func(h Any) Any {
										i := fn.(Fn8)
										return i(a, b, c, d, e, f, g, h)
									}
								}
							}
						}
					}
				}
			}
		}
	}

}
