package purescript_node_fs

import (
	"io"
	"os"

	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Node.FS.Sync")

	// readFile :: Encoding -> FilePath -> Effect String
	exports["readFile"] = func(encoding_ Any, path_ Any) Any {
		return func() Any {
			path := path_.(string)
			data, err := os.ReadFile(path)
			if err != nil {
				panic(err)
			}
			return string(data)
		}
	}

	// readTextFile :: FilePath -> Effect String  
	exports["readTextFile"] = func(path_ Any) Any {
		return func() Any {
			path := path_.(string)
			data, err := os.ReadFile(path)
			if err != nil {
				panic(err)
			}
			return string(data)
		}
	}

	// writeFile :: Encoding -> FilePath -> String -> Effect Unit
	exports["writeFile"] = func(encoding_ Any, path_ Any, data_ Any) Any {
		return func() Any {
			path := path_.(string)
			data := data_.(string)
			err := os.WriteFile(path, []byte(data), 0644)
			if err != nil {
				panic(err)
			}
			return nil
		}
	}

	// writeTextFile :: FilePath -> String -> Effect Unit
	exports["writeTextFile"] = func(path_ Any, data_ Any) Any {
		return func() Any {
			path := path_.(string)
			data := data_.(string)
			err := os.WriteFile(path, []byte(data), 0644)
			if err != nil {
				panic(err)
			}
			return nil
		}
	}

	// appendFile :: Encoding -> FilePath -> String -> Effect Unit
	exports["appendFile"] = func(encoding_ Any, path_ Any, data_ Any) Any {
		return func() Any {
			path := path_.(string)
			data := data_.(string)
			f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			_, err = f.WriteString(data)
			if err != nil {
				panic(err)
			}
			return nil
		}
	}

	// appendTextFile :: FilePath -> String -> Effect Unit
	exports["appendTextFile"] = func(path_ Any, data_ Any) Any {
		return func() Any {
			path := path_.(string)
			data := data_.(string)
			f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			_, err = f.WriteString(data)
			if err != nil {
				panic(err)
			}
			return nil
		}
	}

	// exists :: FilePath -> Effect Boolean
	exports["exists"] = func(path_ Any) Any {
		return func() Any {
			path := path_.(string)
			_, err := os.Stat(path)
			return err == nil
		}
	}

	// mkdir :: FilePath -> Effect Unit
	exports["mkdir"] = func(path_ Any) Any {
		return func() Any {
			path := path_.(string)
			err := os.Mkdir(path, 0755)
			if err != nil {
				panic(err)
			}
			return nil
		}
	}

	// mkdir' :: FilePath -> { recursive :: Boolean, mode :: Int } -> Effect Unit
	exports["mkdir'"] = func(path_ Any, opts_ Any) Any {
		return func() Any {
			path := path_.(string)
			opts := opts_.(Dict)
			
			var err error
			if recursive, ok := opts["recursive"].(bool); ok && recursive {
				mode := 0755
				if m, ok := opts["mode"].(int); ok {
					mode = m
				}
				err = os.MkdirAll(path, os.FileMode(mode))
			} else {
				err = os.Mkdir(path, 0755)
			}
			
			if err != nil {
				panic(err)
			}
			return nil
		}
	}

	// readdir :: FilePath -> Effect (Array String)
	exports["readdir"] = func(path_ Any) Any {
		return func() Any {
			path := path_.(string)
			entries, err := os.ReadDir(path)
			if err != nil {
				panic(err)
			}
			names := make([]Any, len(entries))
			for i, entry := range entries {
				names[i] = entry.Name()
			}
			return names
		}
	}

	// rmdir :: FilePath -> Effect Unit
	exports["rmdir"] = func(path_ Any) Any {
		return func() Any {
			path := path_.(string)
			err := os.Remove(path)
			if err != nil {
				panic(err)
			}
			return nil
		}
	}

	// rmdir' :: FilePath -> { recursive :: Boolean } -> Effect Unit
	exports["rmdir'"] = func(path_ Any, opts_ Any) Any {
		return func() Any {
			path := path_.(string)
			opts := opts_.(Dict)
			
			var err error
			if recursive, ok := opts["recursive"].(bool); ok && recursive {
				err = os.RemoveAll(path)
			} else {
				err = os.Remove(path)
			}
			
			if err != nil {
				panic(err)
			}
			return nil
		}
	}

	// unlink :: FilePath -> Effect Unit
	exports["unlink"] = func(path_ Any) Any {
		return func() Any {
			path := path_.(string)
			err := os.Remove(path)
			if err != nil {
				panic(err)
			}
			return nil
		}
	}

	// rename :: FilePath -> FilePath -> Effect Unit
	exports["rename"] = func(oldPath_ Any, newPath_ Any) Any {
		return func() Any {
			oldPath := oldPath_.(string)
			newPath := newPath_.(string)
			err := os.Rename(oldPath, newPath)
			if err != nil {
				panic(err)
			}
			return nil
		}
	}

	// copyFile :: FilePath -> FilePath -> Effect Unit
	exports["copyFile"] = func(src_ Any, dst_ Any) Any {
		return func() Any {
			src := src_.(string)
			dst := dst_.(string)
			
			source, err := os.Open(src)
			if err != nil {
				panic(err)
			}
			defer source.Close()
			
			destination, err := os.Create(dst)
			if err != nil {
				panic(err)
			}
			defer destination.Close()
			
			_, err = io.Copy(destination, source)
			if err != nil {
				panic(err)
			}
			return nil
		}
	}

	// stat :: FilePath -> Effect Stats
	exports["stat"] = func(path_ Any) Any {
		return func() Any {
			path := path_.(string)
			info, err := os.Stat(path)
			if err != nil {
				panic(err)
			}
			
			return Dict{
				"isFile":      info.Mode().IsRegular(),
				"isDirectory": info.Mode().IsDir(),
				"isSymbolicLink": info.Mode()&os.ModeSymlink != 0,
				"size":        int(info.Size()),
				"atime":       info.ModTime(), // Go doesn't track atime separately
				"mtime":       info.ModTime(),
				"ctime":       info.ModTime(), // Go doesn't have ctime
			}
		}
	}
}

