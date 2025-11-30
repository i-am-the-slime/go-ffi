package purescript_fetch

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Fetch")

	// Helper to create async fetch operation
	makeAsyncFetch := func(url_ Any, options_ Any) Any {
		url := url_.(string)
		options := options_.(Dict)
		
		// Create the async effect that returns a canceler
		return func(callback_ Any) Any {
			return func() Any {
				callback := callback_.(func(Any) Any)
				
				// Build request
				method := "GET"
				if m, ok := options["method"].(string); ok {
					method = strings.ToUpper(m)
				}
				
				var body io.Reader
				if b, ok := options["body"].(string); ok {
					body = strings.NewReader(b)
				}
				
				req, err := http.NewRequest(method, url, body)
				if err != nil {
					// Call callback with Left (error)
					Apply(callback, Dict{"Left": err.Error()})
					// Return canceler (no-op)
					return func() Any { return nil }
				}
				
				// Set headers
				if headers, ok := options["headers"].(Dict); ok {
					for key, value := range headers {
						if val, ok := value.(string); ok {
							req.Header.Set(key, val)
						}
					}
				}
				
				// Make request in goroutine
				done := make(chan struct{})
				go func() {
					defer close(done)
					
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						// Call callback with Left (error)
						Apply(callback, Dict{"Left": err.Error()})
						return
					}
					
					// Read body
					bodyBytes, err := io.ReadAll(resp.Body)
					resp.Body.Close()
					if err != nil {
						Apply(callback, Dict{"Left": err.Error()})
						return
					}
					
					// Build headers dict
					respHeaders := make(Dict)
					for key, values := range resp.Header {
						respHeaders[strings.ToLower(key)] = strings.Join(values, ", ")
					}
					
					// Build response object
					response := Dict{
						"status":     resp.StatusCode,
						"statusText": resp.Status,
						"headers":    respHeaders,
						"body":       string(bodyBytes),
						"url":        url,
						"ok":         resp.StatusCode >= 200 && resp.StatusCode < 300,
					}
					
					// Call callback with Right (success)
					Apply(callback, Dict{"Right": response})
				}()
				
				// Return canceler function
				return func() Any {
					// Could cancel the request here if needed
					return nil
				}
			}
		}
	}

	// fetch :: String -> Aff Response
	exports["fetch"] = func(url_ Any) Any {
		return makeAsyncFetch(url_, Dict{})
	}

	// fetch' :: String -> Options -> Aff Response
	exports["fetch'"] = func(url_ Any, options_ Any) Any {
		return makeAsyncFetch(url_, options_)
	}

	// text :: Response -> Aff String
	exports["text"] = func(response_ Any) Any {
		return func(callback_ Any) Any {
			return func() Any {
				callback := callback_.(func(Any) Any)
				response := response_.(Dict)
				
				if body, ok := response["body"].(string); ok {
					Apply(callback, Dict{"Right": body})
				} else {
					Apply(callback, Dict{"Left": "No body in response"})
				}
				
				return func() Any { return nil }
			}
		}
	}

	// json :: Response -> Aff Json
	exports["json"] = func(response_ Any) Any {
		return func(callback_ Any) Any {
			return func() Any {
				callback := callback_.(func(Any) Any)
				response := response_.(Dict)
				
				// Just return the body string - PureScript will parse it
				if body, ok := response["body"].(string); ok {
					Apply(callback, Dict{"Right": body})
				} else {
					Apply(callback, Dict{"Left": "No body in response"})
				}
				
				return func() Any { return nil }
			}
		}
	}

	// arrayBuffer :: Response -> Aff ArrayBuffer
	exports["arrayBuffer"] = func(response_ Any) Any {
		return func(callback_ Any) Any {
			return func() Any {
				callback := callback_.(func(Any) Any)
				response := response_.(Dict)
				
				if body, ok := response["body"].(string); ok {
					Apply(callback, Dict{"Right": []byte(body)})
				} else {
					Apply(callback, Dict{"Left": "No body in response"})
				}
				
				return func() Any { return nil }
			}
		}
	}

	// status :: Response -> Int
	exports["status"] = func(response_ Any) Any {
		response := response_.(Dict)
		if status, ok := response["status"].(int); ok {
			return status
		}
		return 0
	}

	// statusText :: Response -> String
	exports["statusText"] = func(response_ Any) Any {
		response := response_.(Dict)
		if statusText, ok := response["statusText"].(string); ok {
			return statusText
		}
		return ""
	}

	// headers :: Response -> Headers
	exports["headers"] = func(response_ Any) Any {
		response := response_.(Dict)
		if headers, ok := response["headers"].(Dict); ok {
			return headers
		}
		return Dict{}
	}

	// header :: String -> Response -> Maybe String
	exports["header"] = func(name_ Any, response_ Any) Any {
		name := strings.ToLower(name_.(string))
		response := response_.(Dict)
		
		if headers, ok := response["headers"].(Dict); ok {
			if val, ok := headers[name].(string); ok {
				return Dict{"value0": val} // Just value
			}
		}
		return Dict{} // Nothing
	}

	// url :: Response -> String
	exports["url"] = func(response_ Any) Any {
		response := response_.(Dict)
		if url, ok := response["url"].(string); ok {
			return url
		}
		return ""
	}

	// ok :: Response -> Boolean
	exports["ok"] = func(response_ Any) Any {
		response := response_.(Dict)
		if ok, exists := response["ok"].(bool); exists {
			return ok
		}
		return false
	}
}

// Helper function to build request with options
func buildRequest(url string, options Dict) (*http.Request, error) {
	method := "GET"
	if m, ok := options["method"].(string); ok {
		method = strings.ToUpper(m)
	}
	
	var body io.Reader
	if b, ok := options["body"].(string); ok {
		body = bytes.NewBufferString(b)
	} else if b, ok := options["body"].([]byte); ok {
		body = bytes.NewBuffer(b)
	}
	
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	
	// Set headers
	if headers, ok := options["headers"].(Dict); ok {
		for key, value := range headers {
			if val, ok := value.(string); ok {
				req.Header.Set(key, val)
			}
		}
	}
	
	return req, nil
}

