### Go foreign export implementations for the standard library

This is a fork of [andyarvanitis/purescript-native-go-ffi](https://github.com/andyarvanitis/purescript-native-go-ffi) to support Golang modules.

#### Some basic conventions used in this code (not absolutely required, but please try to follow them)
* Run your code through `gofmt`
* Use provided type alias `Any` instead of `interface{}`
* Use provided type alias `Dict` instead of `map[string]Any` (intended for when you're dealing with ps records)
* Use provided type `[]Any` for arrays
* For input parameters used as-is (i.e. when no type assertion from `Any` is necessary), simply use a name without a trailing underscore.
  * For example:
    ```go
    exports["foo"] = func(bar Any) Any {
      return bar
    }
    ```
* If you need to create a local variable from an argument type assertion (for readability, performance, etc.), name the parameter with a trailing underscore. Then create a type-inferred variable declaration with assignment, with the trailing underscore removed from the name.
  * For example:
  ```go
  exports["foo"] = func(bar_ Any) Any {
    bar := bar_.(int)
    ...
  }
