package purescript_foreign_object

import . "github.com/purescript-native/go-runtime"

func init() {
	exports := Foreign("Foreign.Object.ST")

	// new :: forall a r. ST r (STObject r a)
	exports["new"] = func() Any {
		return make(Dict)
	}

	// peekImpl :: forall a r. (Maybe a) constructors -> String -> STObject r a -> ST r (Maybe a)
	// The first two parameters represent `Just` and `Nothing` constructors respectively.
	exports["peekImpl"] = func(just Any) Any {
		return func(nothing Any) Any {
			return func(k_ Any) Any {
				return func(m_ Any) Any {
					k, ok := k_.(string)
					if !ok {
						panic("Foreign.Object.ST.peekImpl: key is not a String")
					}
					m, ok := m_.(Dict)
					if !ok {
						panic("Foreign.Object.ST.peekImpl: target is not an object (Dict)")
					}
					return func() Any {
						if v, ok := m[k]; ok {
							return Apply(just, v)
						}
						return nothing
					}
				}
			}
		}
	}

	// poke :: forall a r. String -> a -> STObject r a -> ST r (STObject r a)
	exports["poke"] = func(k_ Any) Any {
		return func(v Any) Any {
			return func(m_ Any) Any {
				return func() Any {
					k, ok := k_.(string)
					if !ok {
						panic("Foreign.Object.ST.poke: key is not a String")
					}
					m, ok := m_.(Dict)
					if !ok {
						panic("Foreign.Object.ST.poke: target is not an object (Dict)")
					}
					m[k] = v
					return m
				}
			}
		}
	}

	// pokeAll :: forall a r. Array { key :: String, value :: a } -> STObject r a -> ST r Unit
	// Each element of the kvs array is expected to be a Dict with "key" and "value" entries.
	exports["pokeAll"] = func(kvs_ Any) Any {
		return func(m_ Any) Any {
			return func() Any {
				kvs, ok := kvs_.([]Any)
				if !ok {
					panic("Foreign.Object.ST.pokeAll: expected Array of key/value records")
				}
				m, ok := m_.(Dict)
				if !ok {
					panic("Foreign.Object.ST.pokeAll: target is not an object (Dict)")
				}
				for _, kv := range kvs {
					pair, ok := kv.(Dict)
					if !ok {
						panic("Foreign.Object.ST.pokeAll: element is not a record { key, value }")
					}
					rawKey, hasKey := pair["key"]
					if !hasKey {
						panic("Foreign.Object.ST.pokeAll: record missing 'key' field")
					}
					key, ok := rawKey.(string)
					if !ok {
						panic("Foreign.Object.ST.pokeAll: 'key' field is not a String")
					}
					val, hasVal := pair["value"]
					if !hasVal {
						panic("Foreign.Object.ST.pokeAll: record missing 'value' field")
					}
					m[key] = val
				}
				return m
			}
		}
	}

	// delete :: forall a r. String -> STObject r a -> ST r (STObject r a)
	exports["delete"] = func(k_ Any) Any {
		return func(m_ Any) Any {
			return func() Any {
				k, ok := k_.(string)
				if !ok {
					panic("Foreign.Object.ST.delete: key is not a String")
				}
				m, ok := m_.(Dict)
				if !ok {
					panic("Foreign.Object.ST.delete: target is not an object (Dict)")
				}
				delete(m, k)
				return m
			}
		}
	}
}
