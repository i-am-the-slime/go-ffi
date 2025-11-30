package purescript_node_http

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Node.HTTP")

	// createServer :: (Request -> Response -> Effect Unit) -> Effect Server
	exports["createServer"] = func(handler_ Any) Any {
		return func() Any {
			handler := handler_.(func(Any, Any) Any)
			
			server := &http.Server{
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Wrap request
					req := wrapRequest(r)
					// Wrap response
					res := wrapResponse(w)
					// Call PureScript handler
					effect := handler(req, res)
					// Run the effect
					if eff, ok := effect.(func() Any); ok {
						eff()
					}
				}),
			}
			
			return wrapServer(server)
		}
	}

	// listen :: Server -> Int -> String -> Effect Unit -> Effect Unit
	exports["listen"] = func(server_ Any, port_ Any, hostname_ Any, callback_ Any) Any {
		return func() Any {
			server := server_.(Dict)["_server"].(*http.Server)
			port := port_.(int)
			hostname := hostname_.(string)
			callback := callback_.(func() Any)
			
			addr := fmt.Sprintf("%s:%d", hostname, port)
			server.Addr = addr
			
			go func() {
				err := server.ListenAndServe()
				if err != nil && err != http.ErrServerClosed {
					panic(err)
				}
			}()
			
			// Call the callback
			callback()
			
			return nil
		}
	}

	// close :: Server -> Effect Unit
	exports["close"] = func(server_ Any) Any {
		return func() Any {
			server := server_.(Dict)["_server"].(*http.Server)
			return server.Close()
		}
	}

	// Request methods
	
	// requestMethod :: Request -> String
	exports["requestMethod"] = func(req_ Any) Any {
		req := req_.(Dict)["_request"].(*http.Request)
		return req.Method
	}

	// requestURL :: Request -> String
	exports["requestURL"] = func(req_ Any) Any {
		req := req_.(Dict)["_request"].(*http.Request)
		return req.URL.String()
	}

	// requestHeaders :: Request -> Object String
	exports["requestHeaders"] = func(req_ Any) Any {
		req := req_.(Dict)["_request"].(*http.Request)
		headers := make(Dict)
		for key, values := range req.Header {
			headers[strings.ToLower(key)] = strings.Join(values, ", ")
		}
		return headers
	}

	// requestBody :: Request -> Effect String
	exports["requestBody"] = func(req_ Any) Any {
		return func() Any {
			req := req_.(Dict)["_request"].(*http.Request)
			body, err := io.ReadAll(req.Body)
			if err != nil {
				panic(err)
			}
			return string(body)
		}
	}

	// Response methods
	
	// setStatusCode :: Response -> Int -> Effect Unit
	exports["setStatusCode"] = func(res_ Any, code_ Any) Any {
		return func() Any {
			res := res_.(Dict)["_response"].(http.ResponseWriter)
			code := code_.(int)
			res.WriteHeader(code)
			return nil
		}
	}

	// setHeader :: Response -> String -> String -> Effect Unit
	exports["setHeader"] = func(res_ Any, name_ Any, value_ Any) Any {
		return func() Any {
			res := res_.(Dict)["_response"].(http.ResponseWriter)
			name := name_.(string)
			value := value_.(string)
			res.Header().Set(name, value)
			return nil
		}
	}

	// writeString :: Response -> String -> Effect Unit
	exports["writeString"] = func(res_ Any, data_ Any) Any {
		return func() Any {
			res := res_.(Dict)["_response"].(http.ResponseWriter)
			data := data_.(string)
			_, err := res.Write([]byte(data))
			if err != nil {
				panic(err)
			}
			return nil
		}
	}

	// end :: Response -> Effect Unit
	exports["end"] = func(res_ Any) Any {
		return func() Any {
			// In Go, response is automatically ended when handler returns
			// This is a no-op but kept for API compatibility
			return nil
		}
	}

	// responseFinished :: Response -> Boolean
	exports["responseFinished"] = func(res_ Any) Any {
		// Go doesn't expose this directly, always return false
		return false
	}
}

func wrapRequest(r *http.Request) Dict {
	return Dict{
		"_request": r,
	}
}

func wrapResponse(w http.ResponseWriter) Dict {
	return Dict{
		"_response": w,
	}
}

func wrapServer(s *http.Server) Dict {
	return Dict{
		"_server": s,
	}
}

