package purescript_node_buffer

import (
	"encoding/base64"
	"encoding/hex"

	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Node.Buffer")

	// Buffer is represented as []byte in Go
	
	// create :: Int -> Effect Buffer
	exports["create"] = func(size_ Any) Any {
		return func() Any {
			size := size_.(int)
			return make([]byte, size)
		}
	}

	// fromString :: String -> Encoding -> Effect Buffer
	exports["fromString"] = func(str_ Any, encoding_ Any) Any {
		return func() Any {
			str := str_.(string)
			encoding := encoding_.(string)
			
			switch encoding {
			case "UTF8", "utf8":
				return []byte(str)
			case "Base64", "base64":
				data, err := base64.StdEncoding.DecodeString(str)
				if err != nil {
					return []byte{}
				}
				return data
			case "Hex", "hex":
				data, err := hex.DecodeString(str)
				if err != nil {
					return []byte{}
				}
				return data
			default:
				return []byte(str)
			}
		}
	}

	// fromArray :: Array Int -> Effect Buffer
	exports["fromArray"] = func(arr_ Any) Any {
		return func() Any {
			arr := arr_.([]Any)
			buf := make([]byte, len(arr))
			for i, v := range arr {
				if num, ok := v.(int); ok {
					buf[i] = byte(num)
				}
			}
			return buf
		}
	}

	// fromArrayBuffer :: ArrayBuffer -> Effect Buffer
	exports["fromArrayBuffer"] = func(ab_ Any) Any {
		return func() Any {
			// ArrayBuffer is already []byte in our representation
			if buf, ok := ab_.([]byte); ok {
				result := make([]byte, len(buf))
				copy(result, buf)
				return result
			}
			return []byte{}
		}
	}

	// toString :: Encoding -> Buffer -> Effect String
	exports["toString"] = func(encoding_ Any, buf_ Any) Any {
		return func() Any {
			encoding := encoding_.(string)
			buf := buf_.([]byte)
			
			switch encoding {
			case "UTF8", "utf8":
				return string(buf)
			case "Base64", "base64":
				return base64.StdEncoding.EncodeToString(buf)
			case "Hex", "hex":
				return hex.EncodeToString(buf)
			default:
				return string(buf)
			}
		}
	}

	// toArray :: Buffer -> Effect (Array Int)
	exports["toArray"] = func(buf_ Any) Any {
		return func() Any {
			buf := buf_.([]byte)
			arr := make([]Any, len(buf))
			for i, b := range buf {
				arr[i] = int(b)
			}
			return arr
		}
	}

	// read :: BufferValueType -> Offset -> Buffer -> Effect Number
	exports["read"] = func(valueType_ Any, offset_ Any, buf_ Any) Any {
		return func() Any {
			offset := offset_.(int)
			buf := buf_.([]byte)
			
			if offset >= len(buf) {
				return 0.0
			}
			
			return float64(buf[offset])
		}
	}

	// readString :: Encoding -> Offset -> Offset -> Buffer -> Effect String
	exports["readString"] = func(encoding_ Any, start_ Any, end_ Any, buf_ Any) Any {
		return func() Any {
			encoding := encoding_.(string)
			start := start_.(int)
			end := end_.(int)
			buf := buf_.([]byte)
			
			if start < 0 {
				start = 0
			}
			if end > len(buf) {
				end = len(buf)
			}
			if start >= end {
				return ""
			}
			
			slice := buf[start:end]
			
			switch encoding {
			case "UTF8", "utf8":
				return string(slice)
			case "Base64", "base64":
				return base64.StdEncoding.EncodeToString(slice)
			case "Hex", "hex":
				return hex.EncodeToString(slice)
			default:
				return string(slice)
			}
		}
	}

	// write :: BufferValueType -> Number -> Offset -> Buffer -> Effect Unit
	exports["write"] = func(valueType_ Any, value_ Any, offset_ Any, buf_ Any) Any {
		return func() Any {
			value := value_.(float64)
			offset := offset_.(int)
			buf := buf_.([]byte)
			
			if offset < len(buf) {
				buf[offset] = byte(value)
			}
			
			return nil
		}
	}

	// writeString :: Encoding -> Offset -> Int -> String -> Buffer -> Effect Int
	exports["writeString"] = func(encoding_ Any, offset_ Any, length_ Any, str_ Any, buf_ Any) Any {
		return func() Any {
			encoding := encoding_.(string)
			offset := offset_.(int)
			maxLength := length_.(int)
			str := str_.(string)
			buf := buf_.([]byte)
			
			var data []byte
			switch encoding {
			case "UTF8", "utf8":
				data = []byte(str)
			case "Base64", "base64":
				decoded, err := base64.StdEncoding.DecodeString(str)
				if err != nil {
					return 0
				}
				data = decoded
			case "Hex", "hex":
				decoded, err := hex.DecodeString(str)
				if err != nil {
					return 0
				}
				data = decoded
			default:
				data = []byte(str)
			}
			
			written := 0
			for i := 0; i < len(data) && i < maxLength && offset+i < len(buf); i++ {
				buf[offset+i] = data[i]
				written++
			}
			
			return written
		}
	}

	// toArrayBuffer :: Buffer -> Effect ArrayBuffer
	exports["toArrayBuffer"] = func(buf_ Any) Any {
		return func() Any {
			buf := buf_.([]byte)
			result := make([]byte, len(buf))
			copy(result, buf)
			return result
		}
	}

	// size :: Buffer -> Effect Int
	exports["size"] = func(buf_ Any) Any {
		return func() Any {
			buf := buf_.([]byte)
			return len(buf)
		}
	}

	// concat :: Array Buffer -> Effect Buffer
	exports["concat"] = func(bufs_ Any) Any {
		return func() Any {
			bufs := bufs_.([]Any)
			
			// Calculate total size
			totalSize := 0
			for _, buf_ := range bufs {
				if buf, ok := buf_.([]byte); ok {
					totalSize += len(buf)
				}
			}
			
			// Concatenate
			result := make([]byte, 0, totalSize)
			for _, buf_ := range bufs {
				if buf, ok := buf_.([]byte); ok {
					result = append(result, buf...)
				}
			}
			
			return result
		}
	}

	// concat' :: Array Buffer -> Int -> Effect Buffer
	exports["concat'"] = func(bufs_ Any, length_ Any) Any {
		return func() Any {
			bufs := bufs_.([]Any)
			maxLength := length_.(int)
			
			result := make([]byte, 0, maxLength)
			for _, buf_ := range bufs {
				if buf, ok := buf_.([]byte); ok {
					if len(result)+len(buf) <= maxLength {
						result = append(result, buf...)
					} else {
						remaining := maxLength - len(result)
						if remaining > 0 {
							result = append(result, buf[:remaining]...)
						}
						break
					}
				}
			}
			
			return result
		}
	}

	// copy :: Offset -> Offset -> Buffer -> Offset -> Buffer -> Effect Int
	exports["copy"] = func(srcStart_ Any, srcEnd_ Any, src_ Any, dstStart_ Any, dst_ Any) Any {
		return func() Any {
			srcStart := srcStart_.(int)
			srcEnd := srcEnd_.(int)
			src := src_.([]byte)
			dstStart := dstStart_.(int)
			dst := dst_.([]byte)
			
			if srcStart < 0 {
				srcStart = 0
			}
			if srcEnd > len(src) {
				srcEnd = len(src)
			}
			if dstStart < 0 {
				dstStart = 0
			}
			
			copied := 0
			for i := srcStart; i < srcEnd && dstStart+copied < len(dst); i++ {
				dst[dstStart+copied] = src[i]
				copied++
			}
			
			return copied
		}
	}

	// fill :: Int -> Offset -> Offset -> Buffer -> Effect Unit
	exports["fill"] = func(value_ Any, start_ Any, end_ Any, buf_ Any) Any {
		return func() Any {
			value := byte(value_.(int))
			start := start_.(int)
			end := end_.(int)
			buf := buf_.([]byte)
			
			if start < 0 {
				start = 0
			}
			if end > len(buf) {
				end = len(buf)
			}
			
			for i := start; i < end; i++ {
				buf[i] = value
			}
			
			return nil
		}
	}

	// slice :: Offset -> Offset -> Buffer -> Effect Buffer
	exports["slice"] = func(start_ Any, end_ Any, buf_ Any) Any {
		return func() Any {
			start := start_.(int)
			end := end_.(int)
			buf := buf_.([]byte)
			
			if start < 0 {
				start = 0
			}
			if end > len(buf) {
				end = len(buf)
			}
			if start >= end {
				return []byte{}
			}
			
			result := make([]byte, end-start)
			copy(result, buf[start:end])
			return result
		}
	}
}

