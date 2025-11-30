package purescript_node_path

import (
	"path/filepath"
	"strings"

	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Node.Path")

	// sep :: String
	exports["sep"] = string(filepath.Separator)

	// delimiter :: String
	exports["delimiter"] = string(filepath.ListSeparator)

	// basename :: FilePath -> String
	exports["basename"] = func(path_ Any) Any {
		path := path_.(string)
		return filepath.Base(path)
	}

	// basenameWithoutExt :: FilePath -> String -> String
	exports["basenameWithoutExt"] = func(path_ Any, ext_ Any) Any {
		path := path_.(string)
		ext := ext_.(string)
		base := filepath.Base(path)
		if strings.HasSuffix(base, ext) {
			return strings.TrimSuffix(base, ext)
		}
		return base
	}

	// dirname :: FilePath -> FilePath
	exports["dirname"] = func(path_ Any) Any {
		path := path_.(string)
		return filepath.Dir(path)
	}

	// extname :: FilePath -> String
	exports["extname"] = func(path_ Any) Any {
		path := path_.(string)
		return filepath.Ext(path)
	}

	// normalize :: FilePath -> FilePath
	exports["normalize"] = func(path_ Any) Any {
		path := path_.(string)
		return filepath.Clean(path)
	}

	// isAbsolute :: FilePath -> Boolean
	exports["isAbsolute"] = func(path_ Any) Any {
		path := path_.(string)
		return filepath.IsAbs(path)
	}

	// join :: Array FilePath -> FilePath
	exports["join"] = func(paths_ Any) Any {
		paths := paths_.([]Any)
		strs := make([]string, len(paths))
		for i, p := range paths {
			strs[i] = p.(string)
		}
		return filepath.Join(strs...)
	}

	// relative :: FilePath -> FilePath -> FilePath
	exports["relative"] = func(from_ Any, to_ Any) Any {
		from := from_.(string)
		to := to_.(string)
		rel, err := filepath.Rel(from, to)
		if err != nil {
			// If can't make relative, return the target path
			return to
		}
		return rel
	}

	// resolve :: Array FilePath -> FilePath
	exports["resolve"] = func(paths_ Any) Any {
		paths := paths_.([]Any)
		if len(paths) == 0 {
			return "."
		}
		
		// Join all paths
		strs := make([]string, len(paths))
		for i, p := range paths {
			strs[i] = p.(string)
		}
		joined := filepath.Join(strs...)
		
		// Make absolute
		abs, err := filepath.Abs(joined)
		if err != nil {
			return joined
		}
		return abs
	}

	// parse :: FilePath -> { root :: String, dir :: String, base :: String, ext :: String, name :: String }
	exports["parse"] = func(path_ Any) Any {
		path := path_.(string)
		dir := filepath.Dir(path)
		base := filepath.Base(path)
		ext := filepath.Ext(path)
		name := strings.TrimSuffix(base, ext)
		
		// Determine root (volume on Windows, "/" on Unix)
		root := ""
		if filepath.IsAbs(path) {
			vol := filepath.VolumeName(path)
			if vol != "" {
				root = vol + string(filepath.Separator)
			} else {
				root = string(filepath.Separator)
			}
		}
		
		return Dict{
			"root": root,
			"dir":  dir,
			"base": base,
			"ext":  ext,
			"name": name,
		}
	}

	// concat :: FilePath -> FilePath -> FilePath
	exports["concat"] = func(dir_ Any, base_ Any) Any {
		dir := dir_.(string)
		base := base_.(string)
		return filepath.Join(dir, base)
	}
}

