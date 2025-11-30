package main

import (
	gojson "encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	. "github.com/purescript-native/go-runtime"
	
	// Import all the libraries we need
	_ "github.com/i-am-the-slime/go-ffi/purescript-httpurple"
	_ "github.com/i-am-the-slime/go-ffi/purescript-simple-json"
	_ "github.com/i-am-the-slime/go-ffi/purescript-console"
	_ "github.com/i-am-the-slime/go-ffi/purescript-sqlite3"
)

func main() {
	fmt.Println("🚀 Starting Todo API Server with SQLite...")
	
	// Get FFI modules
	httpurple := Foreign("HTTPurple")
	simpleJSON := Foreign("Simple.JSON")
	console := Foreign("Effect.Console")
	sqlite := Foreign("Database.SQLite3")
	
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
	info := console["info"].(func(Any) Any)
	
	// SQLite functions
	openDB := sqlite["open"].(func(Any) Any)
	exec_ := sqlite["exec'"].(func(Any, Any, Any) Any)
	query_ := sqlite["query'"].(func(Any, Any, Any) Any)
	queryOne_ := sqlite["queryOne'"].(func(Any, Any, Any) Any)
	lastInsertRowId := sqlite["lastInsertRowId"].(func(Any) Any)
	
	// Open database
	dbEffect := openDB("todos.db").(func() Any)
	db := dbEffect()
	
	Run(info("Database opened: todos.db"))
	
	// Initialize database schema
	createTable := `
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			completed INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`
	Run(exec_(createTable, []Any{}, db))
	Run(info("Database schema initialized"))
	
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
				"message": "🎉 Todo API Server with SQLite",
				"version": "2.0.0",
				"database": "todos.db",
				"endpoints": []string{
					"GET    /           - This message",
					"GET    /todos      - List all todos",
					"GET    /todos/:id  - Get todo by ID",
					"POST   /todos      - Create new todo",
					"PUT    /todos/:id  - Update todo",
					"DELETE /todos/:id  - Delete todo",
					"GET    /stats      - Get statistics",
				},
			}
			jsonStr := writeJSON(welcome).(string)
			return json(jsonStr)
		
		// GET /todos - List all todos
		case reqMethod == "GET" && reqPath == "/todos":
			results := query_("SELECT * FROM todos ORDER BY id", []Any{}, db).(func() Any)().([]Any)
			jsonStr := writeJSON(results).(string)
			return json(jsonStr)
		
		// GET /todos/:id - Get specific todo
		case reqMethod == "GET" && len(reqPath) > 7 && reqPath[:7] == "/todos/":
			id := parseID(reqPath[7:])
			if id == -1 {
				return badRequest("Invalid todo ID")
			}
			
			result := queryOne_("SELECT * FROM todos WHERE id = ?", []Any{id}, db).(func() Any)().(Dict)
			
			if _, ok := result["value0"]; !ok {
				return notFound
			}
			
			todo := result["value0"]
			jsonStr := writeJSON(todo).(string)
			return json(jsonStr)
		
		// POST /todos - Create new todo
		case reqMethod == "POST" && reqPath == "/todos":
			// Read body
			bodyEffect := body(req).(func() Any)
			bodyStr := bodyEffect().(string)
			
			// Parse JSON (simplified for demo)
			title := "Untitled"
			description := ""
			completed := 0
			
			// Extract fields from body
			if bodyStr != "" {
				var data map[string]interface{}
				if err := gojson.Unmarshal([]byte(bodyStr), &data); err == nil {
					if t, ok := data["title"].(string); ok {
						title = t
					}
					if d, ok := data["description"].(string); ok {
						description = d
					}
					if c, ok := data["completed"].(bool); ok && c {
						completed = 1
					}
				}
			}
			
			// Insert into database
			Run(exec_(
				"INSERT INTO todos (title, description, completed) VALUES (?, ?, ?)",
				[]Any{title, description, completed},
				db,
			))
			
			// Get the inserted ID
			newID := lastInsertRowId(db).(func() Any)()
			
			Run(info(fmt.Sprintf("Created todo #%d: %s", newID, title)))
			
			// Fetch the created todo
			result := queryOne_("SELECT * FROM todos WHERE id = ?", []Any{newID}, db).(func() Any)().(Dict)
			
			if todo, ok := result["value0"]; ok {
				jsonStr := writeJSON(todo).(string)
				return created(jsonStr)
			}
			
			return badRequest("Failed to create todo")
		
		// PUT /todos/:id - Update todo
		case reqMethod == "PUT" && len(reqPath) > 7 && reqPath[:7] == "/todos/":
			id := parseID(reqPath[7:])
			if id == -1 {
				return badRequest("Invalid todo ID")
			}
			
			// Check if todo exists
			result := queryOne_("SELECT * FROM todos WHERE id = ?", []Any{id}, db).(func() Any)().(Dict)
			if _, ok := result["value0"]; !ok {
				return notFound
			}
			
			// Read body
			bodyEffect := body(req).(func() Any)
			bodyStr := bodyEffect().(string)
			
			// Parse updates
			var data map[string]interface{}
			if err := gojson.Unmarshal([]byte(bodyStr), &data); err != nil {
				return badRequest("Invalid JSON")
			}
			
			// Build update query dynamically
			updates := []string{}
			args := []Any{}
			
			if title, ok := data["title"].(string); ok {
				updates = append(updates, "title = ?")
				args = append(args, title)
			}
			if desc, ok := data["description"].(string); ok {
				updates = append(updates, "description = ?")
				args = append(args, desc)
			}
			if completed, ok := data["completed"].(bool); ok {
				val := 0
				if completed {
					val = 1
				}
				updates = append(updates, "completed = ?")
				args = append(args, val)
			}
			
			if len(updates) == 0 {
				return badRequest("No fields to update")
			}
			
			updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
			args = append(args, id)
			
			updateSQL := "UPDATE todos SET " + strings.Join(updates, ", ") + " WHERE id = ?"
			Run(exec_(updateSQL, args, db))
			
			Run(info(fmt.Sprintf("Updated todo #%d", id)))
			
			// Return updated todo
			updatedResult := queryOne_("SELECT * FROM todos WHERE id = ?", []Any{id}, db).(func() Any)().(Dict)
			if todo, ok := updatedResult["value0"]; ok {
				jsonStr := writeJSON(todo).(string)
				return json(jsonStr)
			}
			
			return badRequest("Failed to fetch updated todo")
		
		// DELETE /todos/:id - Delete todo
		case reqMethod == "DELETE" && len(reqPath) > 7 && reqPath[:7] == "/todos/":
			id := parseID(reqPath[7:])
			if id == -1 {
				return badRequest("Invalid todo ID")
			}
			
			// Check if exists
			result := queryOne_("SELECT * FROM todos WHERE id = ?", []Any{id}, db).(func() Any)().(Dict)
			if _, ok := result["value0"]; !ok {
				return notFound
			}
			
			// Delete
			Run(exec_("DELETE FROM todos WHERE id = ?", []Any{id}, db))
			Run(info(fmt.Sprintf("Deleted todo #%d", id)))
			
			return ok("Deleted")
		
		// GET /stats - Statistics
		case reqMethod == "GET" && reqPath == "/stats":
			totalResult := queryOne_("SELECT COUNT(*) as count FROM todos", []Any{}, db).(func() Any)().(Dict)
			completedResult := queryOne_("SELECT COUNT(*) as count FROM todos WHERE completed = 1", []Any{}, db).(func() Any)().(Dict)
			
			total := 0
			completed := 0
			
			if t, ok := totalResult["value0"].(Dict); ok {
				if count, ok := t["count"].(int64); ok {
					total = int(count)
				}
			}
			if c, ok := completedResult["value0"].(Dict); ok {
				if count, ok := c["count"].(int64); ok {
					completed = int(count)
				}
			}
			
			stats := Dict{
				"total":     total,
				"completed": completed,
				"pending":   total - completed,
				"database":  "SQLite (todos.db)",
			}
			
			jsonStr := writeJSON(stats).(string)
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
		// Close database (would add this with proper close function)
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

