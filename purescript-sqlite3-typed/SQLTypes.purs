module SQLTypes where

import Prelude

import Control.Monad.Except.Trans (ExceptT, runExceptT)
import Control.Monad.Identity.Trans (IdentityT, runIdentityT)
import Data.Array as Array
import Data.Array.NonEmpty.Internal (NonEmptyArray)
import Data.Bifunctor (lmap)
import Data.Either (Either, Either(Left), Either(Right))
import Data.Identity (Identity(Identity))
import Data.List.Types (NonEmptyList)
import Data.Maybe (Maybe, fromJust)
import Data.Newtype (class Newtype, un, wrap)
import Data.Semigroup.Foldable (intercalateMap)
import Data.Traversable (traverse)
import Data.Tuple.Nested (type (/\))
import Foreign (Foreign, ForeignError, readArray, readBoolean, readInt, readNull, readNumber, readString, renderForeignError, unsafeToForeign)
import Heterogeneous.Folding (class Folding, class FoldingWithIndex, class HFoldl, class HFoldlWithIndex, hfoldl, hfoldlWithIndex)
import Prim.Row as Row
import Prim.RowList (class RowToList)
import Prim.RowList as RowList
import Unsafe.Coerce (unsafeCoerce)
import Yoga.JSON.Error (renderHumanError)
import Data.Nullable (toNullable)
import Data.Reflectable (reflectType)
import Type.Proxy (Proxy(Proxy))
import Data.Symbol (class IsSymbol, reflectSymbol)
import Record.Studio.Keys (class Keys, keys)
import Data.Map.Internal (Map)
import Data.Map (empty, insert, lookup) as Map
import Partial.Unsafe (unsafePartial)
import Yoga.JSON (unsafeStringify)
import Data.Set (Set)
import Type.Row.Homogeneous (class Homogeneous)
import Data.Foldable (foldMap, intercalate)
import Data.FoldableWithIndex (foldMapWithIndex)
import Foreign.Object (fromHomogeneous) as Object
import Data.FunctorWithIndex (mapWithIndex)
import Record.Studio.Keys (keys) as Keys
import Record.Studio.Map (mapRecord)
import Record.Studio.MapUniform (mapUniformRecord)
import Record.Builder (Builder)
import Record.Builder (build, buildFromScratch, insert) as Builder
import Safe.Coerce (coerce)
import Prim.Coerce (class Coercible)

data SQLQuery :: Row Type -> Type
data SQLQuery rows = SQLQuery (Array String) String

sqlQueryToString :: forall r. SQLQuery r -> String
sqlQueryToString (SQLQuery _ q) = q

foreign import data SQLParameter :: Type
foreign import data SQLResult :: Type

class ToSQLParam a where
  toSQLParam :: a -> SQLParameter

instance ToSQLParam Int where
  toSQLParam = unsafeCoerce

else instance ToSQLParam String where
  toSQLParam = unsafeCoerce

else instance (ToSQLParam a) => ToSQLParam (Array a) where
  toSQLParam = unsafeCoerce

else instance (ToSQLParam a) => ToSQLParam (NonEmptyArray a) where
  toSQLParam = unsafeCoerce

else instance (ToSQLParam a) => ToSQLParam (Maybe a) where
  toSQLParam = map toSQLParam >>> toNullable >>> unsafeCoerce

else instance (ToSQLParam b, Newtype a b) => ToSQLParam a where
  toSQLParam = unsafeCoerce

else instance (Coercible a b, ToSQLParam b) => ToSQLParam b where
  toSQLParam = unsafeCoerce

fromSQLValue :: forall @a. SQLFromForeign a => SQLResult -> Either String a
fromSQLValue param =
  (un Identity $ runExceptT $ fromSQLResultImpl (unsafeToForeign param)) # lmap
    (intercalateMap ", " renderHumanError)

class SQLFromForeign a where
  fromSQLResultImpl :: Foreign -> ExceptT (NonEmptyList ForeignError) Identity a

instance SQLFromForeign Boolean where
  fromSQLResultImpl = readBoolean

instance SQLFromForeign Int where
  fromSQLResultImpl = readInt

instance SQLFromForeign String where
  fromSQLResultImpl = readString

instance SQLFromForeign Number where
  fromSQLResultImpl = readNumber

instance (SQLFromForeign a) => SQLFromForeign (Array a) where
  fromSQLResultImpl = readArray >=> traverse fromSQLResultImpl

instance (SQLFromForeign a) => SQLFromForeign (Maybe a) where
  fromSQLResultImpl = readNull >=> traverse fromSQLResultImpl

singleResult :: forall @a b. SQLFromForeign a => (a -> b) -> Array SQLResult -> Either String b
singleResult mk = case _ of
  [ x ] -> fromSQLValue @a x <#> mk
  x | len <- Array.length x -> Left $ "Expected exactly one result, got " <> show len

twoResults :: forall @a @b c. SQLFromForeign a => SQLFromForeign b => (a -> b -> c) -> Array SQLResult -> Either String c
twoResults mk = case _ of
  [ _1, _2 ] -> mk <$> fromSQLValue @a _1 <*> fromSQLValue @b _2
  x | len <- Array.length x -> Left $ "Expected exactly two results, got " <> show len

threeResults :: forall @a @b @c d. SQLFromForeign a => SQLFromForeign b => SQLFromForeign c => (a -> b -> c -> d) -> Array SQLResult -> Either String d
threeResults mk = case _ of
  [ _1, _2, _3 ] -> mk <$> fromSQLValue @a _1 <*> fromSQLValue @b _2 <*> fromSQLValue @c _3
  x | len <- Array.length x -> Left $ "Expected exactly three results, got " <> show len

fourResults :: forall @a @b @c @d e. SQLFromForeign a => SQLFromForeign b => SQLFromForeign c => SQLFromForeign d => (a -> b -> c -> d -> e) -> Array SQLResult -> Either String e
fourResults mk = case _ of
  [ _1, _2, _3, _4 ] -> mk <$> fromSQLValue @a _1 <*> fromSQLValue @b _2 <*> fromSQLValue @c _3 <*> fromSQLValue @d _4
  x | len <- Array.length x -> Left $ "Expected exactly four results, got " <> show len

fiveResults :: forall @a @b @c @d @e f. SQLFromForeign a => SQLFromForeign b => SQLFromForeign c => SQLFromForeign d => SQLFromForeign e => (a -> b -> c -> d -> e -> f) -> Array SQLResult -> Either String f
fiveResults mk = case _ of
  [ _1, _2, _3, _4, _5 ] -> mk <$> fromSQLValue @a _1 <*> fromSQLValue @b _2 <*> fromSQLValue @c _3 <*> fromSQLValue @d _4 <*> fromSQLValue @e _5
  x | len <- Array.length x -> Left $ "Expected exactly five results, got " <> show len

class FromResultArray fn a | fn -> a where
  fromResultArray :: fn -> Array SQLResult -> Either String a

instance (SQLFromForeign a, SQLFromForeign b, SQLFromForeign c, SQLFromForeign d, SQLFromForeign e) => FromResultArray (a -> b -> c -> d -> e -> f) f where
  fromResultArray = fiveResults

else instance (SQLFromForeign a, SQLFromForeign b, SQLFromForeign c, SQLFromForeign d) => FromResultArray (a -> b -> c -> d -> e) e where
  fromResultArray = fourResults

else instance (SQLFromForeign a, SQLFromForeign b, SQLFromForeign c) => FromResultArray (a -> b -> c -> d) d where
  fromResultArray = threeResults

else instance (SQLFromForeign a, SQLFromForeign b) => FromResultArray (a -> b -> c) c where
  fromResultArray = twoResults

else instance SQLFromForeign a => FromResultArray (a -> b) b where
  fromResultArray = singleResult

data TurnIntoSQLParam = TurnIntoSQLParam

instance (ToSQLParam n) => Folding TurnIntoSQLParam (Array SQLParameter) n (Array SQLParameter) where
  folding TurnIntoSQLParam acc a = Array.snoc acc (toSQLParam a)

instance
  ( ToSQLParam n
  , IsSymbol sym
  ) =>
  FoldingWithIndex TurnIntoSQLParam (Proxy sym) (Map String SQLParameter) n (Map String SQLParameter) where
  foldingWithIndex TurnIntoSQLParam sym acc a = Map.insert (reflectSymbol sym) (toSQLParam a) acc

newtype SQLBuilder r1 r2 = SQLBuilder (SQLQuery r1 -> SQLQuery r2)

class ToBuilder s r1 r2 | r1 -> r2, r2 -> r1 where
  toBuilder :: s -> SQLBuilder r1 r2

instance ToBuilder (SQLBuilder r1 r2) r1 r2 where
  toBuilder = identity
else instance ToBuilder String r r where
  toBuilder = nonArg

combineBuilders :: forall a b r1 r2 r3. ToBuilder a r2 r3 => ToBuilder b r1 r2 => a -> b -> SQLBuilder r1 r3
combineBuilders a b = SQLBuilder (f <<< g)
  where
  SQLBuilder f = toBuilder a
  SQLBuilder g = toBuilder b

infixl 8 combineBuilders as ^

--arg :: forall @a @b @sym r1 r2. IsSymbol sym => Coercible a b => ToSQLParam b => Row.Cons sym a r1 r2 => SQLBuilder r1 r2
sql :: forall r. SQLBuilder () r -> SQLQuery r
sql (SQLBuilder builder) = builder (SQLQuery [] "")

arg :: forall @a @sym r1 r2. IsSymbol sym => ToSQLParam a => Row.Cons sym a r1 r2 => SQLBuilder r1 r2
arg = SQLBuilder \(SQLQuery oldKeys oldQuery) -> SQLQuery (Array.cons (reflectSymbol (Proxy :: Proxy sym)) oldKeys) ("?" <> oldQuery)

int :: forall @sym r1 r2. IsSymbol sym => Row.Cons sym Int r1 r2 => SQLBuilder r1 r2
int = arg @Int @sym

str :: forall @sym r1 r2. IsSymbol sym => Row.Cons sym String r1 r2 => SQLBuilder r1 r2
str = arg @String @sym

nonArg :: forall r. String -> SQLBuilder r r
nonArg newQuery = SQLBuilder \(SQLQuery keys oldQuery) -> SQLQuery keys (newQuery <> oldQuery)

-- example = sql $
--   "SELECT 1 FROM x WHERE"
--     ! "x="
--     ! int @"foo"
--     ! "hi"
--     ! (string @"bar")

-- x = argsFor { foo: 4, bar: "heinz" } example

argsFor :: forall @params. HFoldlWithIndex TurnIntoSQLParam (Map String SQLParameter) { | params } (Map String SQLParameter) => SQLQuery params -> { | params } -> Array SQLParameter
argsFor (SQLQuery args _) params = do
  let theMap = hfoldlWithIndex TurnIntoSQLParam (Map.empty :: Map String SQLParameter) params
  args <#> \positionedArg -> (unsafePartial $ fromJust $ Map.lookup positionedArg theMap)

--argsFor :: forall @params. HFoldl TurnIntoSQLParam (Array SQLParameter) { | params } (Array SQLParameter) => { | params } -> Array SQLParameter
--argsFor params = hfoldl TurnIntoSQLParam ([] :: Array SQLParameter) params

newtype DatabaseColumns row = DatabaseColumns row

data SQLiteBaseType
  = TextColumn
  | IntColumn
  | RealColumn
  | BooleanColumn
  | BlobColumn
  | JsonBColumn
  | IntegerColumn
  | NumericColumn
  | NullColumn

renderBaseType :: SQLiteBaseType -> String
renderBaseType baseType = case baseType of
  TextColumn -> "TEXT"
  IntColumn -> "INT"
  RealColumn -> "REAL"
  BooleanColumn -> "BOOLEAN"
  BlobColumn -> "BLOB"
  JsonBColumn -> "JSONB"
  IntegerColumn -> "INTEGER"
  NumericColumn -> "NUMERIC"
  NullColumn -> "NULL"

data Constraint = Unique | PrimaryKey | NotNull | ForeignKey String | Check String | Default SQLParameter

renderConstraint :: Constraint -> String
renderConstraint constraint = case constraint of
  Unique -> "UNIQUE"
  PrimaryKey -> "PRIMARY KEY"
  NotNull -> "NOT NULL"
  ForeignKey ref -> "FOREIGN KEY (" <> ref <> ")"
  Check condition -> "CHECK (" <> condition <> ")"
  Default value -> "DEFAULT " <> unsafeStringify value

data SQLColumn = SQLColumn SQLiteBaseType (Array Constraint)

newtype TableName = TableName String

newtype ColumnName = ColumnName String

data Table columns = Table TableName { | columns }

-- Helper for type inference
data MapRecord a b = MapRecord (String -> a -> b)

instance
  ( IsSymbol sym
  , Row.Lacks sym rb
  , Row.Cons sym b rb rc
  ) =>
  FoldingWithIndex
    (MapRecord a b)
    (Proxy sym)
    (Builder { | ra } { | rb })
    a
    (Builder { | ra } { | rc }) where
  foldingWithIndex (MapRecord f) prop rin a = (rin >>> Builder.insert prop (f (reflectSymbol prop) a))

mapRecordWithIndex
  :: forall @a @b @rin @rout
   . HFoldlWithIndex (MapRecord a b) (Builder {} {}) { | rin } (Builder {} { | rout })
  => (String -> a -> b)
  -> { | rin }
  -> { | rout }
mapRecordWithIndex f =
  Builder.buildFromScratch
    <<< hfoldlWithIndex
      (MapRecord f :: MapRecord a b)
      (identity :: Builder {} {})

table :: forall @cols. Homogeneous cols SQLColumn => TableName -> { | cols } -> Table cols
table = Table

columnNamesOf
  :: forall cols colsRL out
   . RowToList cols colsRL
  => HFoldlWithIndex (MapRecord SQLColumn ColumnName)
       (Builder {} {})
       { | cols }
       (Builder {} { | out })
  => Table cols
  -> { | out }
columnNamesOf (Table _ cols) = cols # mapRecordWithIndex \key (_ :: SQLColumn) -> ColumnName key

newtype CreateTableStatement = CreateTableStatement String

createTable :: forall cols. Homogeneous cols SQLColumn => Table cols -> CreateTableStatement
createTable (Table (TableName tableName) tab) = CreateTableStatement $ "CREATE TABLE " <> tableName <> " (" <> cols <> ")"
  where
  toConstraint :: String -> SQLColumn -> String
  toConstraint key (SQLColumn _ constraints) = key <> " " <> (constraints <#> renderConstraint # intercalate " ")
  cols = (Object.fromHomogeneous tab) # mapWithIndex toConstraint # intercalate ", "
