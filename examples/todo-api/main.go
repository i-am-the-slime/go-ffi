package main

import (
	gojson "encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	. "github.com/purescript-native/go-runtime"
	
	// Import all the libraries we need
	_ "github.com/i-am-the-slime/go-ffi/purescript-httpurple"
	_ "github.com/i-am-the-slime/go-ffi/purescript-simple-json"
	_ "github.com/i-am-the-slime/go-ffi/purescript-console"
	_ "github.com/i-am-the-slime/go-ffi/purescript-exceptions"
	_ "github.com/i-am-the-slime/go-ffi/purescript-datetime"
	_ "github.com/i-am-the-slime/go-ffi/purescript-node-fs"
	_ "github.com/i-am-the-slime/go-ffi/purescript-fetch"
)

// Todo represents a single todo item
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// In-memory store for todos
var todos = []Todo{}
var nextID = 1

func main() {
	fmt.Println("🚀 Starting Todo API Server Demo...")
	
	// Get FFI modules
	httpurple := Foreign("HTTPurple")
	simpleJSON := Foreign("Simple.JSON")
	console := Foreign("Effect.Console")
	
	// FFI functions
	serve := httpurple["serve"].(func(Any, Any) Any)
	ok := httpurple["ok"].(func(Any) Any)
	created := httpurple["created"].(func(Any) Any)
	notFound := httpurple["notFound"].(Dict)
	badRequest := httpurple["badRequest"].(func(Any) Any)
	json := httpurple["json"].(func(Any) Any)
	method := httpurple["method"].(func(Any) Any)
	path := httpurple["path"].(func(Any) Any)
	body := httpurple["body"].(func(Any) Any)
	writeJSON := simpleJSON["writeJSON"].(func(Any) Any)
	errorLog := console["error"].(func(Any) Any)
	info := console["info"].(func(Any) Any)
	
	// Initialize with some sample data
	todos = []Todo{
		{
			ID:          1,
			Title:       "Learn PureScript",
			Description: "Understand the basics of PureScript",
			Completed:   true,
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			UpdatedAt:   time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          2,
			Title:       "Build Go FFI",
			Description: "Port PureScript libraries to Go",
			Completed:   true,
			CreatedAt:   time.Now().Add(-12 * time.Hour),
			UpdatedAt:   time.Now().Add(-12 * time.Hour),
		},
		{
			ID:          3,
			Title:       "Create Demo App",
			Description: "Build a real-world example",
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	nextID = 4
	
	// Router function
	router := func(req Any) Any {
		reqMethod := method(req).(string)
		reqPath := path(req).(string)
		
		// Log request
		logMsg := fmt.Sprintf("%s %s", reqMethod, reqPath)
		Run(info(logMsg))
		
		// Route handling
		switch {
		// GET / - Welcome message
		case reqMethod == "GET" && reqPath == "/":
			welcome := Dict{
				"message": "🎉 Todo API Server",
				"version": "1.0.0",
				"endpoints": []string{
					"GET    /           - This message",
					"GET    /todos      - List all todos",
					"GET    /todos/:id  - Get todo by ID",
					"POST   /todos      - Create new todo",
					"PUT    /todos/:id  - Update todo",
					"DELETE /todos/:id  - Delete todo",
					"GET    /stats      - Get statistics",
					"GET    /external   - Fetch external data (demo)",
				},
			}
			jsonStr := writeJSON(welcome).(string)
			return json(jsonStr)
		
		// GET /todos - List all todos
		case reqMethod == "GET" && reqPath == "/todos":
			jsonStr := writeJSON(todos).(string)
			return json(jsonStr)
		
		// GET /todos/:id - Get specific todo
		case reqMethod == "GET" && len(reqPath) > 7 && reqPath[:7] == "/todos/":
			id := parseID(reqPath[7:])
			if id == -1 {
				return badRequest("Invalid todo ID")
			}
			
			for _, todo := range todos {
				if todo.ID == id {
					jsonStr := writeJSON(todo).(string)
					return json(jsonStr)
				}
			}
			return notFound
		
		// POST /todos - Create new todo
		case reqMethod == "POST" && reqPath == "/todos":
			// Read body
			bodyEffect := body(req).(func() Any)
			bodyStr := bodyEffect().(string)
			
			// Parse JSON using Go's encoding/json for the demo
			var data map[string]interface{}
			err := gojson.Unmarshal([]byte(bodyStr), &data)
			if err != nil {
				Run(errorLog(fmt.Sprintf("JSON parse error: %v", err)))
				return badRequest("Invalid JSON")
			}
			
			// Create new todo
			title := "Untitled"
			if t, ok := data["title"].(string); ok {
				title = t
			}
			description := ""
			if d, ok := data["description"].(string); ok {
				description = d
			}
			completed := false
			if c, ok := data["completed"].(bool); ok {
				completed = c
			}
			
			newTodo := Todo{
				ID:          nextID,
				Title:       title,
				Description: description,
				Completed:   completed,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			nextID++
			
			todos = append(todos, newTodo)
			
			Run(info(fmt.Sprintf("Created todo #%d: %s", newTodo.ID, newTodo.Title)))
			
			jsonStr := writeJSON(newTodo).(string)
			return created(jsonStr)
		
		// PUT /todos/:id - Update todo
		case reqMethod == "PUT" && len(reqPath) > 7 && reqPath[:7] == "/todos/":
			id := parseID(reqPath[7:])
			if id == -1 {
				return badRequest("Invalid todo ID")
			}
			
			// Find todo
			todoIndex := -1
			for i, todo := range todos {
				if todo.ID == id {
					todoIndex = i
					break
				}
			}
			
			if todoIndex == -1 {
				return notFound
			}
			
			// Read body
			bodyEffect := body(req).(func() Any)
			bodyStr := bodyEffect().(string)
			
			// Parse JSON
			var data map[string]interface{}
			err := gojson.Unmarshal([]byte(bodyStr), &data)
			if err != nil {
				return badRequest("Invalid JSON")
			}
			
			// Update todo
			if title, ok := data["title"].(string); ok {
				todos[todoIndex].Title = title
			}
			if desc, ok := data["description"].(string); ok {
				todos[todoIndex].Description = desc
			}
			if completed, ok := data["completed"].(bool); ok {
				todos[todoIndex].Completed = completed
			}
			todos[todoIndex].UpdatedAt = time.Now()
			
			Run(info(fmt.Sprintf("Updated todo #%d", id)))
			
			jsonStr := writeJSON(todos[todoIndex]).(string)
			return json(jsonStr)
		
		// DELETE /todos/:id - Delete todo
		case reqMethod == "DELETE" && len(reqPath) > 7 && reqPath[:7] == "/todos/":
			id := parseID(reqPath[7:])
			if id == -1 {
				return badRequest("Invalid todo ID")
			}
			
			// Find and remove todo
			for i, todo := range todos {
				if todo.ID == id {
					todos = append(todos[:i], todos[i+1:]...)
					Run(info(fmt.Sprintf("Deleted todo #%d", id)))
					return ok("Deleted")
				}
			}
			return notFound
		
		// GET /stats - Statistics
		case reqMethod == "GET" && reqPath == "/stats":
			completed := 0
			pending := 0
			for _, todo := range todos {
				if todo.Completed {
					completed++
				} else {
					pending++
				}
			}
			
			stats := Dict{
				"total":     len(todos),
				"completed": completed,
				"pending":   pending,
				"uptime":    "running",
			}
			
			jsonStr := writeJSON(stats).(string)
			return json(jsonStr)
		
		// GET /external - Demo of fetch (calling external API)
		case reqMethod == "GET" && reqPath == "/external":
			// This would normally use Aff/fetch, but for demo we'll return a message
			message := Dict{
				"note": "External API calls would use purescript-fetch here",
				"example": "fetch('https://api.github.com/users/octocat')",
			}
			jsonStr := writeJSON(message).(string)
			return json(jsonStr)
		
		default:
			return notFound
		}
	}
	
	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		Run(info("Shutting down gracefully..."))
		os.Exit(0)
	}()
	
	// Start server
	Run(info("Server starting on http://localhost:3000"))
	Run(info("Press Ctrl+C to stop"))
	effect := serve(3000, router)
	Run(effect)
}

// Helper functions
func parseID(idStr string) int {
	var id int
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		return -1
	}
	return id
}

func getStringOrDefault(data Dict, key string, defaultVal string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return defaultVal
}

func getBoolOrDefault(data Dict, key string, defaultVal bool) bool {
	if val, ok := data[key].(bool); ok {
		return val
	}
	return defaultVal
}

