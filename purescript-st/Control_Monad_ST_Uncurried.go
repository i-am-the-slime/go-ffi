package purescript_st

import (
	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Control.Monad.ST.Uncurried")

	// export const mkSTFn1 = function mkSTFn1(fn) {
	// 	return function(x) {
	// 	  return fn(x)();
	// 	};
	//   };
	exports["mkSTFn1"] = func(fn Any) Any {
		return func(x Any) Any {
			return Run(Apply(fn, x))
		}
	}

	exports["mkSTFn2"] = func(fn Any) Any {
		return func(a, b Any) Any {
			return Run(Apply(fn, a, b))
		}
	}

	exports["mkSTFn3"] = func(fn Any) Any {
		return func(a, b, c Any) Any {
			return Run(Apply(fn, a, b, c))
		}
	}

	exports["mkSTFn4"] = func(fn Any) Any {
		return func(a, b, c, d Any) Any {
			return Run(Apply(fn, a, b, c, d))
		}
	}

	exports["mkSTFn5"] = func(fn Any) Any {
		return func(a, b, c, d, e Any) Any {
			return Run(Apply(fn, a, b, c, d, e))
		}
	}

	exports["mkSTFn6"] = func(fn Any) Any {
		return func(a, b, c, d, e, f Any) Any {
			return Run(Apply(fn, a, b, c, d, e, f))
		}
	}

	exports["mkSTFn7"] = func(fn Any) Any {
		return func(a, b, c, d, e, f, g Any) Any {
			return Run(Apply(fn, a, b, c, d, e, f, g))
		}
	}

	exports["mkSTFn8"] = func(fn Any) Any {
		return func(a, b, c, d, e, f, g, h Any) Any {
			return Run(Apply(fn, a, b, c, d, e, f, g, h))
		}
	}

	exports["mkSTFn9"] = func(fn Any) Any {
		return func(a, b, c, d, e, f, g, h, i Any) Any {
			return Run(Apply(fn, a, b, c, d, e, f, g, h, i))
		}
	}

	exports["mkSTFn10"] = func(fn Any) Any {
		return func(a, b, c, d, e, f, g, h, i, j Any) Any {
			return Run(Apply(fn, a, b, c, d, e, f, g, h, i, j))
		}
	}

	exports["runSTFn1"] = func(fn Any) Any {
		return func(a Any) Any {
			return func() Any {
				return Apply(fn, a)
			}
		}
	}

	exports["runSTFn2"] = func(fn_ Any) Any {
		return func(a Any) Any {
			return func(b Any) Any {
				return func() Any {
					return fn_.(Fn2)(a, b)
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

	exports["runSTFn6"] = func(fn_ Any) Any {
		return func(a Any) Any {
			return func(b Any) Any {
				return func(c Any) Any {
					return func(d Any) Any {
						return func(e Any) Any {
							return func(f Any) Any {
								return func() Any {
									fn := fn_.(Fn6)
									return fn(a, b, c, d, e, f)
								}
							}
						}
					}
				}
			}
		}
	}

	exports["runSTFn7"] = func(fn_ Any) Any {
		return func(a Any) Any {
			return func(b Any) Any {
				return func(c Any) Any {
					return func(d Any) Any {
						return func(e Any) Any {
							return func(f Any) Any {
								return func(g Any) Any {
									return func() Any {
										fn := fn_.(Fn7)
										return fn(a, b, c, d, e, f, g)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	exports["runSTFn8"] = func(fn_ Any) Any {
		return func(a Any) Any {
			return func(b Any) Any {
				return func(c Any) Any {
					return func(d Any) Any {
						return func(e Any) Any {
							return func(f Any) Any {
								return func(g Any) Any {
									return func(h Any) Any {
										return func() Any {
											fn := fn_.(Fn8)
											return fn(a, b, c, d, e, f, g, h)
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

	exports["runSTFn9"] = func(fn_ Any) Any {
		return func(a Any) Any {
			return func(b Any) Any {
				return func(c Any) Any {
					return func(d Any) Any {
						return func(e Any) Any {
							return func(f Any) Any {
								return func(g Any) Any {
									return func(h Any) Any {
										return func(i Any) Any {
											return func() Any {
												fn := fn_.(Fn9)
												return fn(a, b, c, d, e, f, g, h, i)
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
	}

	exports["runSTFn10"] = func(fn_ Any) Any {
		return func(a Any) Any {
			return func(b Any) Any {
				return func(c Any) Any {
					return func(d Any) Any {
						return func(e Any) Any {
							return func(f Any) Any {
								return func(g Any) Any {
									return func(h Any) Any {
										return func(i Any) Any {
											return func(j Any) Any {
												return func() Any {
													fn := fn_.(Fn10)
													return fn(a, b, c, d, e, f, g, h, i, j)
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
		}
	}

}
