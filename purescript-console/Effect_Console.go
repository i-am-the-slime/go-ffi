package purescript_console

import (
	"fmt"
	"os"

	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Effect.Console")

	// log :: String -> Effect Unit
	exports["log"] = func(msg_ Any) Any {
		return func() Any {
			msg := msg_.(string)
			fmt.Println(msg)
			return nil
		}
	}

	// logShow :: forall a. Show a => a -> Effect Unit
	exports["logShow"] = func(val_ Any) Any {
		return func() Any {
			fmt.Println(val_)
			return nil
		}
	}

	// warn :: String -> Effect Unit
	exports["warn"] = func(msg_ Any) Any {
		return func() Any {
			msg := msg_.(string)
			fmt.Fprintln(os.Stderr, "[WARN]", msg)
			return nil
		}
	}

	// error :: String -> Effect Unit
	exports["error"] = func(msg_ Any) Any {
		return func() Any {
			msg := msg_.(string)
			fmt.Fprintln(os.Stderr, "[ERROR]", msg)
			return nil
		}
	}

	// info :: String -> Effect Unit
	exports["info"] = func(msg_ Any) Any {
		return func() Any {
			msg := msg_.(string)
			fmt.Println("[INFO]", msg)
			return nil
		}
	}

	// debug :: String -> Effect Unit
	exports["debug"] = func(msg_ Any) Any {
		return func() Any {
			msg := msg_.(string)
			fmt.Println("[DEBUG]", msg)
			return nil
		}
	}

	// time :: String -> Effect Unit
	exports["time"] = func(label_ Any) Any {
		return func() Any {
			label := label_.(string)
			fmt.Printf("[TIME] %s: timer started\n", label)
			return nil
		}
	}

	// timeLog :: String -> Effect Unit
	exports["timeLog"] = func(label_ Any) Any {
		return func() Any {
			label := label_.(string)
			fmt.Printf("[TIME] %s: ...\n", label)
			return nil
		}
	}

	// timeEnd :: String -> Effect Unit
	exports["timeEnd"] = func(label_ Any) Any {
		return func() Any {
			label := label_.(string)
			fmt.Printf("[TIME] %s: timer ended\n", label)
			return nil
		}
	}

	// clear :: Effect Unit
	exports["clear"] = func() Any {
		return func() Any {
			fmt.Print("\033[H\033[2J")
			return nil
		}
	}
}
