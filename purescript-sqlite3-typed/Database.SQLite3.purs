module Database.SQLite3 where

import Prelude

import Data.Either (Either(..))
import Data.Maybe (Maybe)
import Data.Newtype (class Newtype)
import Data.Nullable (Nullable)
import Data.Nullable as Nullable
import Data.Traversable (traverse)
import Effect (Effect)
import SQLTypes (class FromResultArray, SQLParameter, SQLResult, fromResultArray)
import Safe.Coerce (coerce)

-- Foreign types matching our Go FFI
foreign import data DBConnection :: Type
foreign import data Statement :: Type -> Type -> Type

-- FFI imports matching our Go implementation
foreign import open :: String -> Effect DBConnection
foreign import close :: DBConnection -> Effect Unit
foreign import exec' :: String -> Array SQLParameter -> DBConnection -> Effect Unit
foreign import query' :: String -> Array SQLParameter -> DBConnection -> Effect (Array (Array SQLResult))
foreign import queryOne' :: String -> Array SQLParameter -> DBConnection -> Effect (Maybe (Array SQLResult))
foreign import lastInsertRowId :: DBConnection -> Effect Int
foreign import beginTransaction :: DBConnection -> Effect Transaction
foreign import commit :: Transaction -> Effect Unit
foreign import rollback :: Transaction -> Effect Unit

foreign import data Transaction :: Type

-- Prepare statement (we'll emulate this with our query functions)
newtype PreparedStatement i o = PreparedStatement 
  { sql :: String
  , db :: DBConnection
  }

prepare :: forall @i @o. String -> DBConnection -> Effect (Statement i o)
prepare sql db = pure $ coerce (PreparedStatement { sql, db })

-- Run a statement (INSERT/UPDATE/DELETE)
type InfoImpl = { changes :: Int, lastInsertRowid :: Nullable Int }
newtype NumberOfChanges = NumberOfChanges Int
newtype RowId = RowId Int
derive instance Newtype RowId _
type Info = { changes :: NumberOfChanges, lastInsertRowid :: Maybe RowId }

fromRowInfoImpl :: InfoImpl -> Info
fromRowInfoImpl = \{ changes, lastInsertRowid } -> 
  { changes: NumberOfChanges changes
  , lastInsertRowid: Nullable.toMaybe (coerce lastInsertRowid) 
  }

run :: forall i. Array SQLParameter -> Statement i Void -> Effect Info
run params st = do
  let PreparedStatement { sql, db } = coerce st
  exec' sql params db
  -- Get last insert row ID
  rowId <- lastInsertRowId db
  pure { changes: NumberOfChanges 1, lastInsertRowid: Just (RowId rowId) }

run_ :: Statement Void Void -> Effect Info
run_ = run []

-- Query all rows
allRaw :: forall i o. Array SQLParameter -> Statement i o -> Effect (Array (Array SQLResult))
allRaw params st = do
  let PreparedStatement { sql, db } = coerce st
  query' sql params db

all :: forall i @o. 
  Array SQLParameter -> 
  (Array SQLResult -> Either String o) -> 
  Statement i o -> 
  Effect (Either String (Array o))
all params toOutput st = do
  rows <- allRaw params st
  pure $ traverse toOutput rows

all1 :: forall i o. 
  Array SQLParameter -> 
  (Array SQLResult -> Either String o) -> 
  Statement i o -> 
  Effect (Either String o)
all1 params fn st = all params fn st <#>
  case _ of
    Right [ x ] -> Right x
    Right [] -> Left "Expected a single result, got empty array"
    Right arr -> Left $ "Expected a single result, got " <> show (length arr)
    Left err -> Left err

all_ :: forall o. Array SQLParameter -> Statement Void o -> Effect (Array (Array SQLResult))
all_ params = allRaw params

-- For compatibility with the type-safe query builder
type StatementSource = String

statementSource :: forall i o. Statement i o -> StatementSource
statementSource st = 
  let PreparedStatement { sql } = coerce st
  in sql

-- Helper function to create a database
newDB :: String -> {} -> Effect DBConnection
newDB path _ = open path

