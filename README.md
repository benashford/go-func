go-func
=======

[![Build Status](https://travis-ci.org/benashford/go-func.png?branch=master)](https://travis-ci.org/benashford/go-func)

For example uses see [func_test.go](func_test.go).  Typical use-case:

```go
a := []int{1, 2, 3, 4, 5}
f := func(x int) int {
  return x + 1
}
b := Maps(a, f).([]int) // will be the same as []int{2, 3, 4, 5, 6}
```

The functions implemented are vaguely inspired by the way they're implemented in Clojure.  For example in Clojure ```map``` returns a lazy sequence, but ```mapv``` returns a vector containing the results; in this package ```Map``` returns a channel from which the results can be drawn, and ```Maps``` returns a slice that contains all the results.

# Implementation

Requires Go 1.1 or above and is implemented using reflection.  The types of parameters and the return type are defined as ```interface{}``` so a type assertion (as seen in the above example) is required.

# Functions implemented
## map

1. ```channel := Map(data, func)```: Takes two parameters, the first can be either a channel or a slice; the second is a function which must have one parameter of the same type as the contents of the channel or slice.  The return type is a channel of the same type as the return type of the function.  If the type rules are not met this function will panic.  The returned channel is unbuffered, which is otherwise unimportant if the mapping function does not depend on mutible state.
2. ```slice := Maps(data, func)```: Same as above, but a slice is returned. 

## Parallel map

1. ```channel := PMap(data, func)```: The parameters are the same as for ```Map```, but the function will be executed in parallel.  Multiple goroutines are used, by default the number used will be the same as logical CPUs on the machine.  Please note that you may need to set the ```GOMAXPROCS``` environment variable in order to achieve performance.  Unlike the non-parallel map, the channel is buffered.  Order is preserved.
2. ```slice := PMaps(data, func)```: Same as above, but a slice is returned.

## filter

1. ```channel := Map(data, func)```: Takes two parameters, the first can be either a channel or a slice; the second is a function which must have one parameter of the same type as the contents of the channel or slice.  The return type must be boolean.  If type rules are not met this function will panic.  The returned channel is unbuffered, which is otherwise unimportant if the mapping function does not depend on mutible state.
2. ```slice := Filters(data, func)```: Same as above, but a slice is returned.

## reduce

1. ```result := Reduce(data, func)```: Takes two parameters, the first can be either a channel or a slice; the second is a function which must have two parameters, the second of which must have the same type as the contents of the channel or slice, the first must have the same type as the return type of the function.  The type of the result will be the same as the return type of the function.

## group-by

Examples to follow

## index-by

Examples to follow

# Examples
```go
a := []int{1, 2, 3, 4, 5, 6}
b := PMap(a, func(x int) int { return x * 3})
c := Filter(b, func(x int) bool { return x % 2 == 0})
d := Reduce(c, func(a, b int) int { return a + b}).(int) 
# d will equal 36
```

# Performance
TBC