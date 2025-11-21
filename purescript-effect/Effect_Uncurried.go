package purescript_effect

import (
	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Effect.Uncurried")

	exports["mkEffectFn1"] = func(fn Any) Any {
		return func(x Any) Any {
			return Run(Apply(fn, x))
		}
	}

	exports["mkEffectFn2"] = func(fn Any) Any {
		return func(x Any, y Any) Any {
			return Run(Apply(Apply(fn, x), y))
		}
	}

	exports["mkEffectFn3"] = func(fn Any) Any {
		return func(x Any, y Any, z Any) Any {
			return Run(Apply(Apply(Apply(fn, x), y), z))
		}
	}

	exports["mkEffectFn4"] = func(fn Any) Any {
		return func(a Any, b Any, c Any, d Any) Any {
			return Run(Apply(Apply(Apply(Apply(fn, a), b), c), d))
		}
	}

	exports["mkEffectFn5"] = func(fn Any) Any {
		return func(a Any, b Any, c Any, d Any, e Any) Any {
			return Run(Apply(Apply(Apply(Apply(Apply(fn, a), b), c), d), e))
		}
	}

	exports["mkEffectFn6"] = func(fn Any) Any {
		return func(a Any, b Any, c Any, d Any, e Any, f Any) Any {
			return Run(Apply(Apply(Apply(Apply(Apply(Apply(fn, a), b), c), d), e), f))
		}
	}

	exports["runEffectFn1"] = func(fn Any) Any {
		return func(a Any) Any {
			return func() Any {
				return Apply(fn, a)
			}
		}
	}

	exports["runEffectFn2"] = func(fn_ Any) Any {
		return func(a Any) Any {
			return func(b Any) Any {
				return func() Any {
					fn := fn_.(Fn2)
					return fn(a, b)
				}
			}
		}
	}

	exports["runEffectFn3"] = func(fn_ Any) Any {
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

	exports["runEffectFn4"] = func(fn_ Any) Any {
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

	exports["runEffectFn5"] = func(fn_ Any) Any {
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

	exports["runEffectFn6"] = func(fn_ Any) Any {
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

	exports["runEffectFn7"] = func(fn_ Any) Any {
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

	exports["runEffectFn8"] = func(fn_ Any) Any {
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

	exports["runEffectFn9"] = func(fn_ Any) Any {
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

	exports["runEffectFn10"] = func(fn_ Any) Any {
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
