go-func
=======

[![Build Status](https://travis-ci.org/benashford/go-func.png?branch=master)](https://travis-ci.org/benashford/go-func)

An experiment in both higher-order functions, and how to implement them in Go.

Starting with map, this package provides three functions, vaguely inspired by how map works in Clojure:

* Map(slice, function) - will return a channel which is lazily populated (fed by a Go routine with no buffer) with the results of the function against every item in the slice.
* Maps(slice, function) - same as above, but will eagerly consume the channel and create a result slice.
* MapChan(channel, function) - will return a channel which is lazily populated (fed by a Go routine with no buffer) with the results of the function against every item in the source channel.

Still to do: other higher-order functions (map, reduce, etc.)

For example uses see go_test.go.  Typical use-case:

```go
a := []int{1, 2, 3, 4, 5}
f := func(x int) int {
  return x + 1
}
b := Map(a, f).([]int) // will be the same as []int{2, 3, 4, 5, 6}
```

(Since Go lacks generics, the type assertion is required for the results.)

# Performance note

This is implemented using reflection, so is not to be expected to be as fast as low-level operations.
