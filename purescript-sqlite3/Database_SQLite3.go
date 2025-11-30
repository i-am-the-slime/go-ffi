package purescript_sqlite3

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	. "github.com/purescript-native/go-runtime"
)

func init() {
	exports := Foreign("Database.SQLite3")

	// open :: String -> Effect Database
	exports["open"] = func(path_ Any) Any {
		return func() Any {
			path := path_.(string)
			db, err := sql.Open("sqlite3", path)
			if err != nil {
				panic(Dict{
					"message": fmt.Sprintf("Failed to open database: %v", err),
					"stack":   "",
				})
			}
			return db
		}
	}

	// close :: Database -> Effect Unit
	exports["close"] = func(db_ Any) Any {
		return func() Any {
			db := db_.(*sql.DB)
			err := db.Close()
			if err != nil {
				panic(Dict{
					"message": fmt.Sprintf("Failed to close database: %v", err),
					"stack":   "",
				})
			}
			return nil
		}
	}

	// exec :: String -> Database -> Effect Unit
	exports["exec"] = func(query_ Any, db_ Any) Any {
		return func() Any {
			query := query_.(string)
			db := db_.(*sql.DB)
			_, err := db.Exec(query)
			if err != nil {
				panic(Dict{
					"message": fmt.Sprintf("Query failed: %v", err),
					"stack":   "",
				})
			}
			return nil
		}
	}

	// exec' :: String -> Array Any -> Database -> Effect Unit
	exports["exec'"] = func(query_ Any, args_ Any, db_ Any) Any {
		return func() Any {
			query := query_.(string)
			args := args_.([]Any)
			db := db_.(*sql.DB)
			_, err := db.Exec(query, args...)
			if err != nil {
				panic(Dict{
					"message": fmt.Sprintf("Query failed: %v", err),
					"stack":   "",
				})
			}
			return nil
		}
	}

	// query :: String -> Database -> Effect (Array (Object Any))
	exports["query"] = func(query_ Any, db_ Any) Any {
		return func() Any {
			query := query_.(string)
			db := db_.(*sql.DB)
			
			rows, err := db.Query(query)
			if err != nil {
				panic(Dict{
					"message": fmt.Sprintf("Query failed: %v", err),
					"stack":   "",
				})
			}
			defer rows.Close()
			
			return scanRows(rows)
		}
	}

	// query' :: String -> Array Any -> Database -> Effect (Array (Object Any))
	exports["query'"] = func(query_ Any, args_ Any, db_ Any) Any {
		return func() Any {
			query := query_.(string)
			args := args_.([]Any)
			db := db_.(*sql.DB)
			
			rows, err := db.Query(query, args...)
			if err != nil {
				panic(Dict{
					"message": fmt.Sprintf("Query failed: %v", err),
					"stack":   "",
				})
			}
			defer rows.Close()
			
			return scanRows(rows)
		}
	}

	// queryOne :: String -> Database -> Effect (Maybe (Object Any))
	exports["queryOne"] = func(query_ Any, db_ Any) Any {
		return func() Any {
			query := query_.(string)
			db := db_.(*sql.DB)
			
			rows, err := db.Query(query)
			if err != nil {
				panic(Dict{
					"message": fmt.Sprintf("Query failed: %v", err),
					"stack":   "",
				})
			}
			defer rows.Close()
			
			results := scanRows(rows)
			if len(results) > 0 {
				return Dict{"value0": results[0]} // Just
			}
			return Dict{} // Nothing
		}
	}

	// queryOne' :: String -> Array Any -> Database -> Effect (Maybe (Object Any))
	exports["queryOne'"] = func(query_ Any, args_ Any, db_ Any) Any {
		return func() Any {
			query := query_.(string)
			args := args_.([]Any)
			db := db_.(*sql.DB)
			
			rows, err := db.Query(query, args...)
			if err != nil {
				panic(Dict{
					"message": fmt.Sprintf("Query failed: %v", err),
					"stack":   "",
				})
			}
			defer rows.Close()
			
			results := scanRows(rows)
			if len(results) > 0 {
				return Dict{"value0": results[0]} // Just
			}
			return Dict{} // Nothing
		}
	}

	// lastInsertRowId :: Database -> Effect Int
	exports["lastInsertRowId"] = func(db_ Any) Any {
		return func() Any {
			db := db_.(*sql.DB)
			var id int64
			err := db.QueryRow("SELECT last_insert_rowid()").Scan(&id)
			if err != nil {
				panic(Dict{
					"message": fmt.Sprintf("Failed to get last insert ID: %v", err),
					"stack":   "",
				})
			}
			return int(id)
		}
	}

	// beginTransaction :: Database -> Effect Transaction
	exports["beginTransaction"] = func(db_ Any) Any {
		return func() Any {
			db := db_.(*sql.DB)
			tx, err := db.Begin()
			if err != nil {
				panic(Dict{
					"message": fmt.Sprintf("Failed to begin transaction: %v", err),
					"stack":   "",
				})
			}
			return tx
		}
	}

	// commit :: Transaction -> Effect Unit
	exports["commit"] = func(tx_ Any) Any {
		return func() Any {
			tx := tx_.(*sql.Tx)
			err := tx.Commit()
			if err != nil {
				panic(Dict{
					"message": fmt.Sprintf("Failed to commit transaction: %v", err),
					"stack":   "",
				})
			}
			return nil
		}
	}

	// rollback :: Transaction -> Effect Unit
	exports["rollback"] = func(tx_ Any) Any {
		return func() Any {
			tx := tx_.(*sql.Tx)
			err := tx.Rollback()
			if err != nil {
				panic(Dict{
					"message": fmt.Sprintf("Failed to rollback transaction: %v", err),
					"stack":   "",
				})
			}
			return nil
		}
	}

	// execTx :: String -> Array Any -> Transaction -> Effect Unit
	exports["execTx"] = func(query_ Any, args_ Any, tx_ Any) Any {
		return func() Any {
			query := query_.(string)
			args := args_.([]Any)
			tx := tx_.(*sql.Tx)
			_, err := tx.Exec(query, args...)
			if err != nil {
				panic(Dict{
					"message": fmt.Sprintf("Transaction query failed: %v", err),
					"stack":   "",
				})
			}
			return nil
		}
	}
}

// scanRows converts sql.Rows to []Any (array of objects)
func scanRows(rows *sql.Rows) []Any {
	columns, err := rows.Columns()
	if err != nil {
		panic(Dict{
			"message": fmt.Sprintf("Failed to get columns: %v", err),
			"stack":   "",
		})
	}

	results := []Any{}
	
	for rows.Next() {
		// Create a slice of interface{} to hold each column value
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		// Scan the row into the value pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			panic(Dict{
				"message": fmt.Sprintf("Failed to scan row: %v", err),
				"stack":   "",
			})
		}

		// Create a map for this row
		row := make(Dict)
		for i, col := range columns {
			val := values[i]
			
			// Convert []byte to string (SQLite TEXT)
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		panic(Dict{
			"message": fmt.Sprintf("Row iteration error: %v", err),
			"stack":   "",
		})
	}

	return results
}

