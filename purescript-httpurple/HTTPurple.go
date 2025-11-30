package purescript_httpurple

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("HTTPurple")

	// serve :: Int -> (Request -> ResponseM) -> Effect Unit
	exports["serve"] = func(port_ Any, router_ Any) Any {
		return func() Any {
			port := port_.(int)
			router := router_.(func(Any) Any)

			server := &http.Server{
				Addr: fmt.Sprintf(":%d", port),
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Wrap request
					req := wrapRequest(r)
					
					// Call router to get response
					responseFn := router(req)
					
					// Execute the response monad
					if respEffect, ok := responseFn.(func() Any); ok {
						response := respEffect()
						
						// Apply the response
						if resp, ok := response.(Dict); ok {
							applyResponse(w, resp)
						}
					}
				}),
			}

			fmt.Printf("🚀 HTTPurple server listening on port %d\n", port)
			
			err := server.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				panic(err)
			}
			
			return nil
		}
	}

	// ok :: String -> Response
	exports["ok"] = func(body_ Any) Any {
		body := body_.(string)
		return Dict{
			"status": 200,
			"headers": Dict{
				"Content-Type": "text/plain; charset=utf-8",
			},
			"body": body,
		}
	}

	// notFound :: Response
	exports["notFound"] = Dict{
		"status": 404,
		"headers": Dict{
			"Content-Type": "text/plain; charset=utf-8",
		},
		"body": "Not Found",
	}

	// internalServerError :: String -> Response
	exports["internalServerError"] = func(msg_ Any) Any {
		msg := msg_.(string)
		return Dict{
			"status": 500,
			"headers": Dict{
				"Content-Type": "text/plain; charset=utf-8",
			},
			"body": msg,
		}
	}

	// created :: String -> Response
	exports["created"] = func(body_ Any) Any {
		body := body_.(string)
		return Dict{
			"status": 201,
			"headers": Dict{
				"Content-Type": "text/plain; charset=utf-8",
			},
			"body": body,
		}
	}

	// accepted :: Response
	exports["accepted"] = Dict{
		"status": 202,
		"headers": Dict{
			"Content-Type": "text/plain; charset=utf-8",
		},
		"body": "Accepted",
	}

	// noContent :: Response
	exports["noContent"] = Dict{
		"status": 204,
		"headers": Dict{},
		"body":    "",
	}

	// badRequest :: String -> Response
	exports["badRequest"] = func(msg_ Any) Any {
		msg := msg_.(string)
		return Dict{
			"status": 400,
			"headers": Dict{
				"Content-Type": "text/plain; charset=utf-8",
			},
			"body": msg,
		}
	}

	// unauthorized :: Response
	exports["unauthorized"] = Dict{
		"status": 401,
		"headers": Dict{
			"Content-Type": "text/plain; charset=utf-8",
		},
		"body": "Unauthorized",
	}

	// forbidden :: Response
	exports["forbidden"] = Dict{
		"status": 403,
		"headers": Dict{
			"Content-Type": "text/plain; charset=utf-8",
		},
		"body": "Forbidden",
	}

	// json :: String -> Response
	exports["json"] = func(jsonStr_ Any) Any {
		jsonStr := jsonStr_.(string)
		return Dict{
			"status": 200,
			"headers": Dict{
				"Content-Type": "application/json; charset=utf-8",
			},
			"body": jsonStr,
		}
	}

	// json' :: Int -> String -> Response
	exports["json'"] = func(status_ Any, jsonStr_ Any) Any {
		status := status_.(int)
		jsonStr := jsonStr_.(string)
		return Dict{
			"status": status,
			"headers": Dict{
				"Content-Type": "application/json; charset=utf-8",
			},
			"body": jsonStr,
		}
	}

	// html :: String -> Response
	exports["html"] = func(htmlStr_ Any) Any {
		htmlStr := htmlStr_.(string)
		return Dict{
			"status": 200,
			"headers": Dict{
				"Content-Type": "text/html; charset=utf-8",
			},
			"body": htmlStr,
		}
	}

	// redirect :: String -> Response
	exports["redirect"] = func(url_ Any) Any {
		url := url_.(string)
		return Dict{
			"status": 302,
			"headers": Dict{
				"Location": url,
			},
			"body": "",
		}
	}

	// redirect' :: Int -> String -> Response
	exports["redirect'"] = func(status_ Any, url_ Any) Any {
		status := status_.(int)
		url := url_.(string)
		return Dict{
			"status": status,
			"headers": Dict{
				"Location": url,
			},
			"body": "",
		}
	}

	// Request accessors

	// method :: Request -> Method
	exports["method"] = func(req_ Any) Any {
		req := req_.(Dict)["_request"].(*http.Request)
		return strings.ToUpper(req.Method)
	}

	// path :: Request -> String
	exports["path"] = func(req_ Any) Any {
		req := req_.(Dict)["_request"].(*http.Request)
		return req.URL.Path
	}

	// query :: Request -> Object String
	exports["query"] = func(req_ Any) Any {
		req := req_.(Dict)["_request"].(*http.Request)
		params := make(Dict)
		for key, values := range req.URL.Query() {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}
		return params
	}

	// headers :: Request -> Object String
	exports["headers"] = func(req_ Any) Any {
		req := req_.(Dict)["_request"].(*http.Request)
		headers := make(Dict)
		for key, values := range req.Header {
			headers[strings.ToLower(key)] = strings.Join(values, ", ")
		}
		return headers
	}

	// header :: String -> Request -> Maybe String
	exports["header"] = func(name_ Any, req_ Any) Any {
		name := strings.ToLower(name_.(string))
		req := req_.(Dict)["_request"].(*http.Request)
		
		for key, values := range req.Header {
			if strings.ToLower(key) == name && len(values) > 0 {
				return Dict{"value0": values[0]} // Just value
			}
		}
		return Dict{} // Nothing
	}

	// body :: Request -> Effect String
	exports["body"] = func(req_ Any) Any {
		return func() Any {
			req := req_.(Dict)["_request"].(*http.Request)
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				return ""
			}
			return string(bodyBytes)
		}
	}

	// Response modifiers

	// withStatus :: Int -> Response -> Response
	exports["withStatus"] = func(status_ Any, resp_ Any) Any {
		status := status_.(int)
		resp := resp_.(Dict)
		newResp := make(Dict)
		for k, v := range resp {
			newResp[k] = v
		}
		newResp["status"] = status
		return newResp
	}

	// withHeader :: String -> String -> Response -> Response
	exports["withHeader"] = func(name_ Any, value_ Any, resp_ Any) Any {
		name := name_.(string)
		value := value_.(string)
		resp := resp_.(Dict)
		
		newResp := make(Dict)
		for k, v := range resp {
			newResp[k] = v
		}
		
		headers := make(Dict)
		if h, ok := resp["headers"].(Dict); ok {
			for k, v := range h {
				headers[k] = v
			}
		}
		headers[name] = value
		newResp["headers"] = headers
		
		return newResp
	}

	// withHeaders :: Object String -> Response -> Response
	exports["withHeaders"] = func(newHeaders_ Any, resp_ Any) Any {
		newHeaders := newHeaders_.(Dict)
		resp := resp_.(Dict)
		
		newResp := make(Dict)
		for k, v := range resp {
			newResp[k] = v
		}
		
		headers := make(Dict)
		if h, ok := resp["headers"].(Dict); ok {
			for k, v := range h {
				headers[k] = v
			}
		}
		for k, v := range newHeaders {
			headers[k] = v
		}
		newResp["headers"] = headers
		
		return newResp
	}

	// withBody :: String -> Response -> Response
	exports["withBody"] = func(body_ Any, resp_ Any) Any {
		body := body_.(string)
		resp := resp_.(Dict)
		
		newResp := make(Dict)
		for k, v := range resp {
			newResp[k] = v
		}
		newResp["body"] = body
		
		return newResp
	}
}

func wrapRequest(r *http.Request) Dict {
	return Dict{
		"_request": r,
	}
}

func applyResponse(w http.ResponseWriter, response Dict) {
	// Set status code
	status := 200
	if s, ok := response["status"].(int); ok {
		status = s
	}

	// Set headers
	if headers, ok := response["headers"].(Dict); ok {
		for key, value := range headers {
			if val, ok := value.(string); ok {
				w.Header().Set(key, val)
			}
		}
	}

	// Write status
	w.WriteHeader(status)

	// Write body
	if body, ok := response["body"].(string); ok {
		w.Write([]byte(body))
	}
}

