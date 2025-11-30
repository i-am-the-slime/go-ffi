# HTTPurple Example Usage

## Simple HTTP Server

```purescript
module Main where

import Prelude
import Effect (Effect)
import HTTPurple (serve, ok, notFound, json, path, method)
import Simple.JSON (writeJSON)

main :: Effect Unit
main = serve 8080 router
  where
    router request = 
      case method request, path request of
        "GET", "/" -> 
          ok "Hello, HTTPurple!"
        
        "GET", "/api/users" ->
          json $ writeJSON { users: ["Alice", "Bob", "Charlie"] }
        
        "GET", "/about" ->
          ok "About Page"
        
        _, _ ->
          notFound

-- Server starts on http://localhost:8080
```

## With Routing

```purescript
module Main where

import Prelude
import Effect (Effect)
import HTTPurple (serve, ok, json, notFound, path)
import Routing.Match (Match, lit, int, str, end)
import Routing (match)
import Data.Either (Either(..))
import Simple.JSON (writeJSON)

-- Define routes
data Route
  = Home
  | User Int
  | Post Int String
  | NotFound

routes :: Match Route
routes = 
  Home <$ lit "home" <* end
  <|> User <$> (lit "users" *> int <* end)
  <|> Post <$> (lit "posts" *> int) <*> (lit "comments" *> str <* end)

main :: Effect Unit
main = serve 8080 router
  where
    router request = 
      case match routes (path request) of
        Right Home ->
          ok "Welcome home!"
        
        Right (User userId) ->
          json $ writeJSON { userId: userId, name: "User " <> show userId }
        
        Right (Post postId commentId) ->
          json $ writeJSON { post: postId, comment: commentId }
        
        Left _ ->
          notFound
```

## With Query Parameters & Headers

```purescript
module Main where

import Prelude
import Effect (Effect)
import HTTPurple as H
import Data.Maybe (Maybe(..))
import Simple.JSON (writeJSON)

main :: Effect Unit
main = H.serve 8080 router
  where
    router request = do
      let queryParams = H.query request
      let authHeader = H.header "authorization" request
      
      case authHeader of
        Just token -> 
          H.json $ writeJSON 
            { message: "Authenticated"
            , params: queryParams
            }
        
        Nothing ->
          H.unauthorized

-- Try: curl http://localhost:8080?foo=bar -H "Authorization: Bearer token"
```

## JSON API with Multiple Methods

```purescript
module Main where

import Prelude
import Effect (Effect)
import Effect.Class (liftEffect)
import HTTPurple as H
import Simple.JSON as JSON
import Data.Either (Either(..))

main :: Effect Unit
main = H.serve 3000 router
  where
    router request = 
      case H.method request, H.path request of
        "GET", "/api/status" ->
          H.json $ JSON.writeJSON { status: "ok", uptime: 100 }
        
        "POST", "/api/users" -> do
          bodyStr <- H.body request
          case JSON.parseJSON bodyStr of
            Right user ->
              H.created $ JSON.writeJSON 
                { id: 123, user: user }
            
            Left err ->
              H.badRequest "Invalid JSON"
        
        "PUT", "/api/users" -> do
          bodyStr <- H.body request
          H.ok "Updated"
        
        "DELETE", "/api/users" ->
          H.noContent
        
        _, _ ->
          H.notFound

-- POST example:
-- curl -X POST http://localhost:3000/api/users \
--   -H "Content-Type: application/json" \
--   -d '{"name":"Alice","email":"alice@example.com"}'
```

## With Custom Response Headers

```purescript
module Main where

import Prelude
import Effect (Effect)
import HTTPurple as H

main :: Effect Unit
main = H.serve 8080 router
  where
    router request = 
      H.ok "CORS enabled!"
        # H.withHeader "Access-Control-Allow-Origin" "*"
        # H.withHeader "Access-Control-Allow-Methods" "GET, POST, OPTIONS"
        # H.withHeader "X-Custom-Header" "my-value"
```

## HTML Server

```purescript
module Main where

import Prelude
import Effect (Effect)
import HTTPurple as H

main :: Effect Unit
main = H.serve 8080 router
  where
    router request = 
      case H.path request of
        "/" ->
          H.html """
            <!DOCTYPE html>
            <html>
              <head><title>HTTPurple</title></head>
              <body>
                <h1>Welcome to HTTPurple!</h1>
                <p>Built with PureScript → Go</p>
              </body>
            </html>
          """
        
        _ ->
          H.notFound
```

## File Upload Handler

```purescript
module Main where

import Prelude
import Effect (Effect)
import HTTPurple as H
import Node.FS.Sync as FS
import Node.Path as Path

main :: Effect Unit
main = H.serve 8080 router
  where
    router request = 
      case H.method request, H.path request of
        "POST", "/upload" -> do
          bodyContent <- H.body request
          let filename = "uploads/file.txt"
          FS.writeTextFile filename bodyContent
          H.created $ "File saved to " <> filename
        
        _, _ ->
          H.notFound
```

## Available Response Helpers

```purescript
-- Success responses
ok :: String -> Response                      -- 200
created :: String -> Response                 -- 201
accepted :: Response                          -- 202
noContent :: Response                         -- 204

-- Client error responses
badRequest :: String -> Response              -- 400
unauthorized :: Response                      -- 401
forbidden :: Response                         -- 403
notFound :: Response                          -- 404

-- Server error
internalServerError :: String -> Response     -- 500

-- Content types
json :: String -> Response                    -- JSON response
json' :: Int -> String -> Response            -- JSON with custom status
html :: String -> Response                    -- HTML response

-- Redirects
redirect :: String -> Response                -- 302 redirect
redirect' :: Int -> String -> Response        -- Custom redirect status

-- Response modifiers
withStatus :: Int -> Response -> Response
withHeader :: String -> String -> Response -> Response
withHeaders :: Object String -> Response -> Response
withBody :: String -> Response -> Response
```

## Request Accessors

```purescript
method :: Request -> Method
path :: Request -> String
query :: Request -> Object String
headers :: Request -> Object String
header :: String -> Request -> Maybe String
body :: Request -> Effect String
```

