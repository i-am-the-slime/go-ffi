# purescript-sqlite3-typed

Type-safe SQL queries for SQLite3 using PureScript's type system!

## Features

✅ **Type-level schema definitions** using row types  
✅ **Compile-time query validation** - wrong column names won't compile  
✅ **SQL injection protection** - automatic parameter binding  
✅ **Type-safe result parsing** - guaranteed correct types  
✅ **Automatic parameter ordering** - no manual array construction  

## Usage

### Define Your Schema

```purescript
import SQLTypes as SQL

-- Define table schema at type level
tableName_todos = SQL.TableName "todos"

tableDefinition_todos = SQL.table tableName_todos
  { id :: SQL.SQLColumn SQL.IntegerColumn [ SQL.PrimaryKey, SQL.NotNull ]
  , title :: SQL.SQLColumn SQL.TextColumn [ SQL.NotNull ]
  , description :: SQL.SQLColumn SQL.TextColumn []
  , completed :: SQL.SQLColumn SQL.BooleanColumn [ SQL.NotNull ]
  , created_at :: SQL.SQLColumn SQL.TextColumn []
  , updated_at :: SQL.SQLColumn SQL.TextColumn []
  }

-- Get column names (for use in queries)
columns_todos = SQL.columnNamesOf tableDefinition_todos
```

### Create Tables

```purescript
import Database.SQLite3 as SQLite

main = do
  db <- SQLite.open "todos.db" {}
  
  -- Generate CREATE TABLE statement from schema
  let SQL.CreateTableStatement createSQL = SQL.createTable tableDefinition_todos
  
  -- Execute it
  st <- SQLite.prepare createSQL db
  _ <- SQLite.run_ st
  
  pure unit
```

### Type-Safe Queries

```purescript
import SQLTypes (sql, argsFor, (^))
import SQLTypes as SQL

-- INSERT with type-safe parameters
insertTodo :: String -> String -> Boolean -> Effect Unit
insertTodo title desc completed = do
  db <- SQLite.open "todos.db" {}
  
  -- Define query with typed parameters
  let query = sql $
        "INSERT INTO todos (title, description, completed) VALUES (" 
          ^ SQL.str @"title"       -- Type: String
          ^ ","
          ^ SQL.str @"description" -- Type: String  
          ^ ","
          ^ SQL.arg @Boolean @"completed"  -- Type: Boolean
          ^ ")"
  
  -- Prepare statement
  st <- SQLite.prepare (SQL.sqlQueryToString query) db
  
  -- Build parameters - order is automatic!
  let params = argsFor query { title, description: desc, completed }
  
  -- Execute
  _ <- SQLite.run params st
  
  pure unit
```

### Type-Safe SELECT

```purescript
-- SELECT with typed results
getAllTodos :: Effect (Array Todo)
getAllTodos = do
  db <- SQLite.open "todos.db" {}
  
  let query = sql $ "SELECT id, title, description, completed FROM todos"
  
  st <- SQLite.prepare (SQL.sqlQueryToString query) db
  
  -- Type-safe result parsing
  result <- SQLite.all [] parseRow st
  
  case result of
    Right todos -> pure todos
    Left err -> throwError $ error err
  
  where
    -- Parser knows the types!
    parseRow :: Array SQLResult -> Either String Todo
    parseRow = SQL.fourResults @Int @String @String @Boolean \id title desc completed ->
      { id, title, description: desc, completed }
```

### WHERE Clauses

```purescript
getCompletedTodos :: Effect (Array Todo)  
getCompletedTodos = do
  db <- SQLite.open "todos.db" {}
  
  let query = sql $
        "SELECT id, title, description, completed FROM todos WHERE completed = " 
          ^ SQL.arg @Boolean @"completed"
  
  st <- SQLite.prepare (SQL.sqlQueryToString query) db
  
  let params = argsFor query { completed: true }
  
  result <- SQLite.all params parseRow st
  
  case result of
    Right todos -> pure todos
    Left err -> throwError $ error err
```

### UPDATE

```purescript
updateTodoCompleted :: Int -> Boolean -> Effect Unit
updateTodoCompleted todoId newCompleted = do
  db <- SQLite.open "todos.db" {}
  
  let query = sql $
        "UPDATE todos SET completed = "
          ^ SQL.arg @Boolean @"completed"
          ^ " WHERE id = "
          ^ SQL.int @"id"
  
  st <- SQLite.prepare (SQL.sqlQueryToString query) db
  
  let params = argsFor query { id: todoId, completed: newCompleted }
  
  _ <- SQLite.run params st
  
  pure unit
```

## Type Safety Guarantees

### Wrong Column Name - Won't Compile!

```purescript
-- ❌ Compile error: No instance for Row.Cons "wrong_name" ...
let query = sql $
      "INSERT INTO todos (title) VALUES (" 
        ^ SQL.str @"wrong_name"  -- Error! Not in schema
        ^ ")"

-- ✅ Compiles: "title" exists in schema  
let query = sql $
      "INSERT INTO todos (title) VALUES (" 
        ^ SQL.str @"title"
        ^ ")"
```

### Wrong Type - Won't Compile!

```purescript
-- ❌ Compile error: Couldn't match type String with Boolean
let params = argsFor query 
  { title: "Todo"
  , completed: "yes"  -- Error! Should be Boolean
  }

-- ✅ Compiles: types match
let params = argsFor query
  { title: "Todo"
  , completed: true
  }
```

### Missing Parameter - Won't Compile!

```purescript
let query = sql $
      "INSERT INTO todos (title, completed) VALUES ("
        ^ SQL.str @"title"
        ^ ","
        ^ SQL.arg @Boolean @"completed"
        ^ ")"

-- ❌ Compile error: Record lacks required label "completed"
let params = argsFor query { title: "Todo" }

-- ✅ Compiles: all parameters provided
let params = argsFor query { title: "Todo", completed: false }
```

## How It Works

### Type-Level Magic

1. **Row Types** - Record types track available columns:
   ```purescript
   type TodoRow = ( id :: Int, title :: String, ... )
   ```

2. **Type-Level Strings** - Column names are types:
   ```purescript
   SQL.str @"title"  -- The string "title" is a type!
   ```

3. **Heterogeneous Folding** - Build parameter arrays in correct order:
   ```purescript
   argsFor query { b: 2, a: 1 }  
   -- Returns [1, 2] if query uses @"a" then @"b"
   ```

4. **Phantom Types** - Track query parameters without runtime cost:
   ```purescript
   data SQLQuery :: Row Type -> Type
   -- Row Type is erased at runtime, only for type checking!
   ```

## SQL Injection Protection

All parameters are automatically escaped:

```purescript
-- Safe! Uses prepared statements
let query = sql $
      "SELECT * FROM todos WHERE title = " ^ SQL.str @"title"

let params = argsFor query { title: "'; DROP TABLE todos; --" }
-- Parameters are properly escaped by SQLite
```

## Performance

- ✅ **Zero runtime overhead** - types erased at compile time
- ✅ **Prepared statements** - SQLite compiles queries once
- ✅ **Parameter binding** - native SQLite3 binding (fast!)

## Comparison

### Without Type Safety (Raw FFI)

```purescript
-- Error-prone!
query' "INSERT INTO todos (title, completed) VALUES (?, ?)" 
  [toSQLParam 42, toSQLParam "wrong order!"]  -- Oops!
  db
```

### With Type Safety (This Library)

```purescript
-- Compile-time safety!
let query = sql $
      "INSERT INTO todos (title, completed) VALUES ("
        ^ SQL.str @"title"
        ^ ","
        ^ SQL.arg @Boolean @"completed"
        ^ ")"

let params = argsFor query 
  { title: 42                -- ❌ Won't compile! Wrong type
  , completed: "wrong"       -- ❌ Won't compile! Wrong type  
  , title: "Right", completed: false  -- ✅ Compiles!
  }
```

## Advanced: Custom Types

```purescript
-- Define newtype with ToSQLParam instance
newtype TodoId = TodoId Int

derive instance Newtype TodoId _

instance ToSQLParam TodoId where
  toSQLParam (TodoId id) = toSQLParam id

-- Use in queries
let query = sql $
      "SELECT * FROM todos WHERE id = " ^ SQL.arg @TodoId @"todoId"
```

## Module Exports

- **SQLTypes** - Query builder, schema definitions
- **Database.SQLite3** - Go FFI wrapper with typed interface

## Requirements

- PureScript compiler
- purescript-sqlite3 (Go FFI)
- purescript-heterogeneous
- purescript-record-studio

Enjoy type-safe SQL! 🎉

