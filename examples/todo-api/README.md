# 🎯 Todo API Server - Complete Demo

A full-featured REST API built with PureScript libraries running on Go!

## Features Demonstrated

✅ **HTTPurple** - HTTP server with routing  
✅ **Simple.JSON** - JSON parsing and encoding  
✅ **Effect.Console** - Structured logging  
✅ **Effect.Exception** - Error handling with try/catch  
✅ **Data.DateTime** - Timestamps for todos  
✅ **In-memory CRUD** - Create, Read, Update, Delete operations  

## Quick Start

```bash
# Build and run
cd examples/todo-api
go build -o todo-server main.go
./todo-server

# Server will start on http://localhost:3000
```

## API Endpoints

### GET /
Welcome message with available endpoints

```bash
curl http://localhost:3000/
```

### GET /todos
List all todos

```bash
curl http://localhost:3000/todos
```

**Response:**
```json
[
  {
    "id": 1,
    "title": "Learn PureScript",
    "description": "Understand the basics of PureScript",
    "completed": true,
    "createdAt": "2025-11-29T10:30:00Z",
    "updatedAt": "2025-11-29T10:30:00Z"
  }
]
```

### GET /todos/:id
Get a specific todo

```bash
curl http://localhost:3000/todos/1
```

### POST /todos
Create a new todo

```bash
curl -X POST http://localhost:3000/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn Go FFI",
    "description": "Master PureScript to Go FFI",
    "completed": false
  }'
```

**Response:** `201 Created`
```json
{
  "id": 4,
  "title": "Learn Go FFI",
  "description": "Master PureScript to Go FFI",
  "completed": false,
  "createdAt": "2025-11-30T20:00:00Z",
  "updatedAt": "2025-11-30T20:00:00Z"
}
```

### PUT /todos/:id
Update an existing todo

```bash
curl -X PUT http://localhost:3000/todos/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn PureScript - Updated",
    "completed": true
  }'
```

### DELETE /todos/:id
Delete a todo

```bash
curl -X DELETE http://localhost:3000/todos/1
```

**Response:** `200 OK`
```
Deleted
```

### GET /stats
Get statistics about todos

```bash
curl http://localhost:3000/stats
```

**Response:**
```json
{
  "total": 3,
  "completed": 2,
  "pending": 1,
  "uptime": "running"
}
```

### GET /external
Demo endpoint showing where external API calls would go

```bash
curl http://localhost:3000/external
```

## Full Demo Script

```bash
# 1. Start server (in one terminal)
./todo-server

# 2. In another terminal, run these commands:

# Welcome message
curl http://localhost:3000/

# List todos
curl http://localhost:3000/todos

# Get stats
curl http://localhost:3000/stats

# Create a new todo
curl -X POST http://localhost:3000/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Build amazing apps","description":"Using PureScript on Go","completed":false}'

# List todos again (should see the new one)
curl http://localhost:3000/todos

# Update todo #4
curl -X PUT http://localhost:3000/todos/4 \
  -H "Content-Type: application/json" \
  -d '{"completed":true}'

# Get specific todo
curl http://localhost:3000/todos/4

# Delete todo
curl -X DELETE http://localhost:3000/todos/4

# Check stats
curl http://localhost:3000/stats
```

## Error Handling Examples

### Invalid JSON
```bash
curl -X POST http://localhost:3000/todos \
  -H "Content-Type: application/json" \
  -d 'not valid json'
```
**Response:** `400 Bad Request` - "Invalid JSON"

### Non-existent todo
```bash
curl http://localhost:3000/todos/999
```
**Response:** `404 Not Found`

### Invalid ID
```bash
curl http://localhost:3000/todos/abc
```
**Response:** `400 Bad Request` - "Invalid todo ID"

## What This Demonstrates

### 1. **Routing** (HTTPurple)
- Pattern matching on method + path
- Path parameter extraction (`/todos/:id`)
- Multiple HTTP methods (GET, POST, PUT, DELETE)

### 2. **JSON Handling** (Simple.JSON)
- Automatic serialization with `writeJSON`
- Safe parsing with `parseJSON`
- Error handling for invalid JSON

### 3. **Error Handling** (Effect.Exception)
- Try/catch with `try` function
- Returns `Either Error a` for safe operations
- Graceful error responses

### 4. **Logging** (Effect.Console)
- Request logging with `info`
- Error logging with `error`
- Structured log messages

### 5. **Date/Time** (Data.DateTime)
- Automatic timestamps on create/update
- Go's `time.Time` integration

### 6. **HTTP Responses** (HTTPurple)
- `ok` - 200 OK
- `created` - 201 Created
- `badRequest` - 400 Bad Request
- `notFound` - 404 Not Found
- `json` - JSON content type

## Code Highlights

### Clean Error Handling
```go
tryEffect := tryFn(func() Any {
    return parseJSON(bodyStr)
}).(func() Any)
result := tryEffect().(Dict)

if leftErr, hasLeft := result["Left"]; hasLeft {
    Run(errorLog(fmt.Sprintf("JSON parse error: %v", leftErr)))
    return badRequest("Invalid JSON")
}
```

### Pattern Matching Routes
```go
switch {
case reqMethod == "GET" && reqPath == "/todos":
    // List todos
case reqMethod == "POST" && reqPath == "/todos":
    // Create todo
case reqMethod == "PUT" && len(reqPath) > 7 && reqPath[:7] == "/todos/":
    // Update specific todo
}
```

### Graceful Shutdown
```go
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

go func() {
    <-sigChan
    Run(info("Shutting down gracefully..."))
    os.Exit(0)
}()
```

## Next Steps

This demo can be extended with:
- **File persistence** using `purescript-node-fs`
- **External API calls** using `purescript-fetch`
- **Validation** using `purescript-validation`
- **Authentication** with JWT tokens
- **Database integration** (PostgreSQL, MongoDB)
- **WebSocket support**
- **Middleware** for logging/auth/CORS

## Performance

Lightweight and fast:
- Go's native HTTP server
- In-memory operations
- Minimal overhead from FFI
- Production-ready

Enjoy building with PureScript → Go! 🚀

