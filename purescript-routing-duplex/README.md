# purescript-routing-duplex

**Bidirectional routing** for PureScript - parse AND print URLs with the same definition!

## Why Bidirectional?

Traditional routers only parse:
```purescript
parse :: String -> Maybe Route  -- One direction only
```

routing-duplex does **both** with one definition:
```purescript
routes :: RouteDuplex' Route
-- Parse: "/users/42" → Right (User 42)
-- Print: User 42 → "/users/42"
```

## Features

✅ **Type-safe routing** - Invalid routes won't compile  
✅ **Bidirectional** - One definition for parsing AND printing  
✅ **Composable** - Build complex routes from simple pieces  
✅ **Zero FFI** - Pure PureScript, works everywhere  
✅ **Type-level validation** - Catch errors at compile time  

## Basic Usage

```purescript
module Main where

import Prelude
import Routing.Duplex (RouteDuplex', parse, print, root, segment)
import Routing.Duplex.Generic (sum, noArgs)
import Routing.Duplex.Generic.Syntax ((/))
import Data.Generic.Rep (class Generic)
import Data.Either (Either(..))

-- Define your routes
data Route
  = Home
  | About
  | User Int
  | UserEdit Int String

derive instance Generic Route _

-- ONE definition for both parsing and printing!
routes :: RouteDuplex' Route
routes = root $ sum
  { "Home": noArgs
  , "About": "about" / noArgs
  , "User": "users" / int segment
  , "UserEdit": "users" / int segment / "edit" / segment
  }

main = do
  -- Parse URLs
  case parse routes "/users/42" of
    Right (User 42) -> log "Parsed user 42!"
    _ -> log "Parse failed"
  
  -- Print URLs
  let url = print routes (User 42)
  log url  -- "/users/42"
  
  -- Both directions guaranteed consistent!
  let userEditUrl = print routes (UserEdit 42 "profile")
  log userEditUrl  -- "/users/42/edit/profile"
```

## Combinators

### Path Segments

```purescript
import Routing.Duplex (segment, root, path)
import Routing.Duplex.Generic.Syntax ((/))

-- Literal segment
"about" / noArgs  -- Matches "/about"

-- Capture segment as String
"users" / segment  -- Captures "/users/alice" → "alice"

-- Capture as Int
"users" / int segment  -- "/users/42" → 42

-- Multiple segments
"api" / "v1" / "users" / int segment
-- "/api/v1/users/42" → 42
```

### Query Parameters

```purescript
import Routing.Duplex.Query (param, flag, many, int) as Query

routes = root $ sum
  { "Search": "search" 
      ? Query.param "q"           -- ?q=purescript
      ? Query.int "page"          -- &page=2
      ? Query.flag "exact"        -- &exact (boolean)
  , "Filter": "items"
      ? Query.many "tag"          -- ?tag=foo&tag=bar
  }

-- Parse: "/search?q=purescript&page=2&exact"
-- Result: Search "purescript" 2 true

-- Print: Search "purescript" 2 true
-- Result: "/search?q=purescript&page=2&exact"
```

### Optional Segments

```purescript
import Routing.Duplex (optional, default)

routes = root $ sum
  { "UserProfile": "users" 
      / int segment 
      / optional segment  -- Optional tab: /users/42 or /users/42/posts
  }

-- Parse: "/users/42" → UserProfile 42 Nothing
-- Parse: "/users/42/posts" → UserProfile 42 (Just "posts")
```

### Custom Parsers

```purescript
import Routing.Duplex (RouteDuplex, prefix, suffix)
import Data.Profunctor (dimap)

-- Parse UUID format
uuid :: RouteDuplex' String String
uuid = dimap toString fromString segment
  where
    toString (UUID s) = s
    fromString s = UUID s

-- Use in routes
routes = "users" / uuid  -- /users/550e8400-e29b-41d4-a716-446655440000
```

## Integration with HTTPurple

```purescript
import HTTPurple (serve)
import Routing.Duplex (parse)

serve { port: 3000 } \{ method, path } ->
  case method, parse routes path of
    Get, Right Home -> 
      ok "Welcome home!"
    
    Get, Right (User id) ->
      ok $ "User profile for ID: " <> show id
    
    Post, Right (User id) ->
      created "User updated!"
    
    _, Left err ->
      notFound
```

## Type Safety Examples

### Wrong Route - Won't Compile

```purescript
-- ❌ Constructor "NonExistent" is not in data type Route
routes = root $ sum
  { "NonExistent": "wrong" / noArgs  -- Error!
  }
```

### Type Mismatch - Won't Compile

```purescript
data Route = User String  -- Takes String

routes = root $ sum
  { "User": "users" / int segment  -- ❌ Error! Returns Int, expects String
  }
```

### Missing Case - Won't Compile

```purescript
data Route = Home | About

routes = root $ sum
  { "Home": noArgs
  -- ❌ Error! Missing "About" case
  }
```

## Advanced: Nested Routes

```purescript
data AdminRoute = AdminHome | AdminUsers | AdminSettings
data Route = Home | Admin AdminRoute | Api ApiRoute

derive instance Generic AdminRoute _
derive instance Generic Route _

adminRoutes :: RouteDuplex' AdminRoute  
adminRoutes = sum
  { "AdminHome": noArgs
  , "AdminUsers": "users" / noArgs
  , "AdminSettings": "settings" / noArgs
  }

routes :: RouteDuplex' Route
routes = root $ sum
  { "Home": noArgs
  , "Admin": "admin" / adminRoutes  -- Nest admin routes!
  , "Api": "api" / apiRoutes
  }

-- Parse: "/admin/users" → Admin AdminUsers
-- Print: Admin AdminSettings → "/admin/settings"
```

## Performance

- ✅ **Zero runtime overhead** - Types erased at compile time
- ✅ **Fast parsing** - Optimized string operations
- ✅ **Fast printing** - Direct string concatenation
- ✅ **No reflection** - Pure functions, no magic

## Comparison with Other Routers

### purescript-routing (one-way)

```purescript
-- Can only parse
match :: RouteDuplex' Route -> String -> Maybe Route
-- Cannot print back to URLs!
```

### routing-duplex (bidirectional)

```purescript
-- Both directions
parse :: RouteDuplex' Route -> String -> Either Error Route
print :: RouteDuplex' Route -> Route -> String
-- One definition, two operations!
```

## How It Works

### Profunctor Magic

```purescript
newtype RouteDuplex i o = RouteDuplex
  { print :: i -> String      -- i → String (print)
  , parse :: String -> o      -- String → o (parse)
  }

-- i and o are usually the same type (Route)
-- but can differ for transformations
```

### Generic Deriving

```purescript
derive instance Generic Route _

-- Automatically generates:
-- - Constructor names for URL segments
-- - Parsing/printing for sum types
-- - Composition for product types
```

## Dependencies

All pure PureScript, no FFI:
- prelude
- profunctor
- record
- strings
- generic-rep (for automatic deriving)

## Real-World Example

```purescript
module Blog.Routes where

import Routing.Duplex
import Routing.Duplex.Generic
import Routing.Duplex.Query as Query

data Route
  = Home
  | BlogPost String
  | BlogTag String (Array String)
  | Author String (Maybe String)
  | Search String Int

derive instance Generic Route _

routes :: RouteDuplex' Route
routes = root $ sum
  { "Home": noArgs
  , "BlogPost": "posts" / segment
  , "BlogTag": "tags" / segment ? Query.many "filter"
  , "Author": "authors" / segment / optional segment
  , "Search": "search" ? Query.param "q" ? Query.int "page"
  }

-- Usage
case parse routes path of
  Right (BlogPost slug) -> renderPost slug
  Right (Search query page) -> searchPosts query page  
  Right Home -> renderHomepage
  Left err -> show404
  
-- Generate links
linkToPost = print routes (BlogPost "my-awesome-post")
-- "/posts/my-awesome-post"

linkToSearch = print routes (Search "purescript" 2)
-- "/search?q=purescript&page=2"
```

## Credits

Ported from [purescript-routing-duplex](https://github.com/natefaubion/purescript-routing-duplex) by @natefaubion

Pure PureScript - works with any backend! 🎉

