package purescript_sqlite3

import (
	"os"
	"testing"

	. "github.com/purescript-native/go-runtime"
)

func TestOpenClose(t *testing.T) {
	exports := Foreign("Database.SQLite3")
	open := exports["open"].(func(Any) Any)
	close := exports["close"].(func(Any) Any)
	
	// Open in-memory database
	dbEffect := open(":memory:").(func() Any)
	db := dbEffect()
	
	// Close database
	closeEffect := close(db).(func() Any)
	closeEffect()
}

func TestExec(t *testing.T) {
	exports := Foreign("Database.SQLite3")
	open := exports["open"].(func(Any) Any)
	exec := exports["exec"].(func(Any, Any) Any)
	close := exports["close"].(func(Any) Any)
	
	db := open(":memory:").(func() Any)()
	defer close(db).(func() Any)()
	
	// Create table
	execEffect := exec("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)", db).(func() Any)
	execEffect()
	
	// Insert data
	insertEffect := exec("INSERT INTO users (name) VALUES ('Alice')", db).(func() Any)
	insertEffect()
}

func TestExecWithParams(t *testing.T) {
	exports := Foreign("Database.SQLite3")
	open := exports["open"].(func(Any) Any)
	exec_ := exports["exec'"].(func(Any, Any, Any) Any)
	close := exports["close"].(func(Any) Any)
	
	db := open(":memory:").(func() Any)()
	defer close(db).(func() Any)()
	
	// Create table
	exec_("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, age INTEGER)", []Any{}, db).(func() Any)()
	
	// Insert with parameters
	execEffect := exec_("INSERT INTO users (name, age) VALUES (?, ?)", []Any{"Bob", 30}, db).(func() Any)
	execEffect()
}

func TestQuery(t *testing.T) {
	exports := Foreign("Database.SQLite3")
	open := exports["open"].(func(Any) Any)
	exec_ := exports["exec'"].(func(Any, Any, Any) Any)
	query := exports["query"].(func(Any, Any) Any)
	close := exports["close"].(func(Any) Any)
	
	db := open(":memory:").(func() Any)()
	defer close(db).(func() Any)()
	
	// Setup
	exec_("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)", []Any{}, db).(func() Any)()
	exec_("INSERT INTO users (name) VALUES (?)", []Any{"Alice"}, db).(func() Any)()
	exec_("INSERT INTO users (name) VALUES (?)", []Any{"Bob"}, db).(func() Any)()
	
	// Query
	queryEffect := query("SELECT * FROM users", db).(func() Any)
	results := queryEffect().([]Any)
	
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
	
	row1 := results[0].(Dict)
	if row1["name"] != "Alice" {
		t.Errorf("Expected name 'Alice', got %v", row1["name"])
	}
}

func TestQueryWithParams(t *testing.T) {
	exports := Foreign("Database.SQLite3")
	open := exports["open"].(func(Any) Any)
	exec_ := exports["exec'"].(func(Any, Any, Any) Any)
	query_ := exports["query'"].(func(Any, Any, Any) Any)
	close := exports["close"].(func(Any) Any)
	
	db := open(":memory:").(func() Any)()
	defer close(db).(func() Any)()
	
	// Setup
	exec_("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, age INTEGER)", []Any{}, db).(func() Any)()
	exec_("INSERT INTO users (name, age) VALUES (?, ?)", []Any{"Alice", 25}, db).(func() Any)()
	exec_("INSERT INTO users (name, age) VALUES (?, ?)", []Any{"Bob", 30}, db).(func() Any)()
	
	// Query with parameter
	queryEffect := query_("SELECT * FROM users WHERE age > ?", []Any{26}, db).(func() Any)
	results := queryEffect().([]Any)
	
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
	
	row := results[0].(Dict)
	if row["name"] != "Bob" {
		t.Errorf("Expected name 'Bob', got %v", row["name"])
	}
}

func TestQueryOne(t *testing.T) {
	exports := Foreign("Database.SQLite3")
	open := exports["open"].(func(Any) Any)
	exec_ := exports["exec'"].(func(Any, Any, Any) Any)
	queryOne := exports["queryOne"].(func(Any, Any) Any)
	close := exports["close"].(func(Any) Any)
	
	db := open(":memory:").(func() Any)()
	defer close(db).(func() Any)()
	
	// Setup
	exec_("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)", []Any{}, db).(func() Any)()
	exec_("INSERT INTO users (name) VALUES (?)", []Any{"Alice"}, db).(func() Any)()
	
	// Query one
	queryEffect := queryOne("SELECT * FROM users LIMIT 1", db).(func() Any)
	result := queryEffect().(Dict)
	
	if _, ok := result["value0"]; !ok {
		t.Error("Expected Just, got Nothing")
	}
	
	row := result["value0"].(Dict)
	if row["name"] != "Alice" {
		t.Errorf("Expected name 'Alice', got %v", row["name"])
	}
}

func TestQueryOneNotFound(t *testing.T) {
	exports := Foreign("Database.SQLite3")
	open := exports["open"].(func(Any) Any)
	exec_ := exports["exec'"].(func(Any, Any, Any) Any)
	queryOne := exports["queryOne"].(func(Any, Any) Any)
	close := exports["close"].(func(Any) Any)
	
	db := open(":memory:").(func() Any)()
	defer close(db).(func() Any)()
	
	// Setup empty table
	exec_("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)", []Any{}, db).(func() Any)()
	
	// Query one (should return Nothing)
	queryEffect := queryOne("SELECT * FROM users", db).(func() Any)
	result := queryEffect().(Dict)
	
	if len(result) != 0 {
		t.Error("Expected Nothing for empty result set")
	}
}

func TestLastInsertRowId(t *testing.T) {
	exports := Foreign("Database.SQLite3")
	open := exports["open"].(func(Any) Any)
	exec_ := exports["exec'"].(func(Any, Any, Any) Any)
	lastInsertRowId := exports["lastInsertRowId"].(func(Any) Any)
	close := exports["close"].(func(Any) Any)
	
	db := open(":memory:").(func() Any)()
	defer close(db).(func() Any)()
	
	// Setup
	exec_("CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)", []Any{}, db).(func() Any)()
	exec_("INSERT INTO users (name) VALUES (?)", []Any{"Alice"}, db).(func() Any)()
	
	// Get last insert ID
	idEffect := lastInsertRowId(db).(func() Any)
	id := idEffect().(int)
	
	if id != 1 {
		t.Errorf("Expected ID 1, got %d", id)
	}
	
	// Insert another
	exec_("INSERT INTO users (name) VALUES (?)", []Any{"Bob"}, db).(func() Any)()
	id2 := lastInsertRowId(db).(func() Any)().(int)
	
	if id2 != 2 {
		t.Errorf("Expected ID 2, got %d", id2)
	}
}

func TestTransaction(t *testing.T) {
	exports := Foreign("Database.SQLite3")
	open := exports["open"].(func(Any) Any)
	exec_ := exports["exec'"].(func(Any, Any, Any) Any)
	beginTransaction := exports["beginTransaction"].(func(Any) Any)
	execTx := exports["execTx"].(func(Any, Any, Any) Any)
	commit := exports["commit"].(func(Any) Any)
	query := exports["query"].(func(Any, Any) Any)
	close := exports["close"].(func(Any) Any)
	
	db := open(":memory:").(func() Any)()
	defer close(db).(func() Any)()
	
	// Setup
	exec_("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)", []Any{}, db).(func() Any)()
	
	// Begin transaction
	txEffect := beginTransaction(db).(func() Any)
	tx := txEffect()
	
	// Execute in transaction
	execTx("INSERT INTO users (name) VALUES (?)", []Any{"Alice"}, tx).(func() Any)()
	execTx("INSERT INTO users (name) VALUES (?)", []Any{"Bob"}, tx).(func() Any)()
	
	// Commit
	commit(tx).(func() Any)()
	
	// Verify
	results := query("SELECT * FROM users", db).(func() Any)().([]Any)
	if len(results) != 2 {
		t.Errorf("Expected 2 results after commit, got %d", len(results))
	}
}

func TestTransactionRollback(t *testing.T) {
	exports := Foreign("Database.SQLite3")
	open := exports["open"].(func(Any) Any)
	exec_ := exports["exec'"].(func(Any, Any, Any) Any)
	beginTransaction := exports["beginTransaction"].(func(Any) Any)
	execTx := exports["execTx"].(func(Any, Any, Any) Any)
	rollback := exports["rollback"].(func(Any) Any)
	query := exports["query"].(func(Any, Any) Any)
	close := exports["close"].(func(Any) Any)
	
	db := open(":memory:").(func() Any)()
	defer close(db).(func() Any)()
	
	// Setup
	exec_("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)", []Any{}, db).(func() Any)()
	
	// Begin transaction
	txEffect := beginTransaction(db).(func() Any)
	tx := txEffect()
	
	// Execute in transaction
	execTx("INSERT INTO users (name) VALUES (?)", []Any{"Alice"}, tx).(func() Any)()
	
	// Rollback
	rollback(tx).(func() Any)()
	
	// Verify (should be empty)
	results := query("SELECT * FROM users", db).(func() Any)().([]Any)
	if len(results) != 0 {
		t.Errorf("Expected 0 results after rollback, got %d", len(results))
	}
}

func TestFileDatabase(t *testing.T) {
	exports := Foreign("Database.SQLite3")
	open := exports["open"].(func(Any) Any)
	exec_ := exports["exec'"].(func(Any, Any, Any) Any)
	query := exports["query"].(func(Any, Any) Any)
	close := exports["close"].(func(Any) Any)
	
	dbFile := "test.db"
	defer os.Remove(dbFile)
	
	// Create database file
	db := open(dbFile).(func() Any)()
	
	// Create table and insert
	exec_("CREATE TABLE test (id INTEGER PRIMARY KEY, value TEXT)", []Any{}, db).(func() Any)()
	exec_("INSERT INTO test (value) VALUES (?)", []Any{"persistent"}, db).(func() Any)()
	
	// Close
	close(db).(func() Any)()
	
	// Reopen
	db2 := open(dbFile).(func() Any)()
	defer close(db2).(func() Any)()
	
	// Query (data should persist)
	results := query("SELECT * FROM test", db2).(func() Any)().([]Any)
	if len(results) != 1 {
		t.Errorf("Expected 1 result from persistent database, got %d", len(results))
	}
	
	row := results[0].(Dict)
	if row["value"] != "persistent" {
		t.Errorf("Expected value 'persistent', got %v", row["value"])
	}
}

