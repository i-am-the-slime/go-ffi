package purescript_node_process

import (
	"os"
	"runtime"

	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Node.Process")

	// argv :: Array String
	exports["argv"] = func() Any {
		args := make([]Any, len(os.Args))
		for i, arg := range os.Args {
			args[i] = arg
		}
		return args
	}

	// execArgv :: Array String
	exports["execArgv"] = func() Any {
		// Go doesn't have exact equivalent, return empty array
		return []Any{}
	}

	// execPath :: String
	exports["execPath"] = func() Any {
		executable, err := os.Executable()
		if err != nil {
			return ""
		}
		return executable
	}

	// chdir :: String -> Effect Unit
	exports["chdir"] = func(dir_ Any) Any {
		return func() Any {
			dir := dir_.(string)
			return os.Chdir(dir)
		}
	}

	// cwd :: Effect String
	exports["cwd"] = func() Any {
		cwd, err := os.Getwd()
		if err != nil {
			return ""
		}
		return cwd
	}

	// exit :: Int -> Effect Unit
	exports["exit"] = func(code_ Any) Any {
		return func() Any {
			code := code_.(int)
			os.Exit(code)
			return nil
		}
	}

	// lookupEnv :: String -> Effect (Maybe String)
	exports["lookupEnv"] = func(key_ Any) Any {
		return func() Any {
			key := key_.(string)
			value, exists := os.LookupEnv(key)
			if exists {
				return Dict{"value0": value} // Just value
			}
			return Dict{} // Nothing
		}
	}

	// setEnv :: String -> String -> Effect Unit
	exports["setEnv"] = func(key_ Any, value_ Any) Any {
		return func() Any {
			key := key_.(string)
			value := value_.(string)
			return os.Setenv(key, value)
		}
	}

	// platform :: String
	exports["platform"] = runtime.GOOS

	// arch :: String  
	exports["arch"] = runtime.GOARCH

	// pid :: Int
	exports["pid"] = os.Getpid()

	// version :: String
	exports["version"] = runtime.Version()

	// stdin :: Readable
	exports["stdin"] = os.Stdin

	// stdout :: Writable
	exports["stdout"] = os.Stdout

	// stderr :: Writable
	exports["stderr"] = os.Stderr
}

