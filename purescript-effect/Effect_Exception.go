package purescript_effect

import (
	"errors"

	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Effect.Exception")

	exports["showErrorImpl"] = func(e_ Any) Any {
		e := e_.(error)
		return e.Error()
	}

	exports["error"] = func(s_ Any) Any {
		s := s_.(string)
		return errors.New(s)
	}

	exports["message"] = func(e_ Any) Any {
		e := e_.(error)
		return e.Error()
	}

}
