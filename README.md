# fuego - Functional Experiment in Go.

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/seborama/fuego) [![Build Status](https://api.travis-ci.org/seborama/fuego.svg?branch=master)](https://travis-ci.org/seborama/fuego) [![goreportcard](https://img.shields.io/badge/go%20report-A%2B-brightgreen.svg)](http://goreportcard.com/report/seborama/fuego) [![License](https://img.shields.io/:license-apache2-blue.svg)](https://github.com/seborama/fuego/blob/master/LICENSE) [![cover.run](https://cover.run/go/github.com/seborama/fuego.svg?style=flat&tag=golang-1.10)](https://cover.run/go?tag=golang-1.10&repo=github.com%2Fseborama%2Ffuego) 

## Overview

_Making Go come._

This is a research project in functional programming which I hope will prove useful!

Fuego brings a few functional paradigms to Go.

## Install

```bash
go get github.com/seborama/fuego
```

Or for a specific version:

```bash
go get gopkg.in/seborama/govcr.v3
```

## Contribute

Contributions and feedback are welcome.

For contributions, you must develop in TDD fashion and ideally provide Go testable examples (if meaningful).

## Main features

The code documentation can be found on [godoc](http://godoc.org/github.com/seborama/fuego).

### Set (partial implementation)

Set is based on hamt.Set and entries must implement interface `hamt.Entry`.

This is an elegant solution from [Yota Toyama](https://github.com/raviqqe) that somewhat mimics generics: the type can be anything so long as it respects the simple behaviour of `hamt.Entry`.

An example `hamt.Entry` implementation called `EntryInt` is provided in [entry_test.go](entry_test.go).

```go
// See entry_test.go for "EntryInt"
NewSet().
    Insert(EntryInt(1)).
    Insert(EntryInt(2)).
    Delete(EntryInt(1)).
    Insert(EntryInt(3)).
    Stream()
```

### Map (partial implementation)

As with Set, Map is based on hamt.Map and entries must implement interface `hamt.Entry` for its keys but values can be anything (`interface{}`).

See [example_map_test.go](example_map_test.go) for more details of an example of Map with Stream and Filter combined together to extract entries which keys are an even number.

### Function

A `Function` is a normal Go function which signature is

```go
func(i interface{}) interface{}
```

See [example_function_test.go](example_function_test.go) for a basic example and the other tests / examples for more uses.

### Iterator

You can create you own `Iterator`'s.

See [iterator.go](iterator.go) for examples of convenience `Iterator`'s:
- NewSliceIterator
- NewSetIterator

```go
NewSliceIterator([]interface{}{2, 3}) // returns an Iterator over []interface{2, 3}

NewSetIterator(NewSet().
    Insert(EntryInt(2))), // returns an Iterator over a Set that contains a single EntryInt(2)
```

### Stream (partial implementation)

#### Creation

```go
someGoSlice := []int{1, 2, 3}
NewStream(
    NewBaseIterable(
        NewSliceIterator(someGoSlice))),
```

```go
NewStream(
    NewBaseIterable(
        NewSetIterator(
            NewSet().
                Insert(EntryInt(1)).
                Insert(EntryInt(2)))))
```

#### Map

```go
// See in this README and in stream_test.go for "functionTimesTwo()"
NewSet().
    Insert(EntryInt(1)).
    Insert(EntryInt(2)).
    Insert(EntryInt(3)).
    Stream().
    Map(functionTimesTwo())
// returns EntryInt's {2,4,6}
```

#### Filter

```go
// See stream_test.go for "newEntryIntEqualsTo()"
NewSet().
    Insert(EntryInt(1)).
    Insert(EntryInt(2)).
    Insert(EntryInt(3)).
    Stream().
    Filter(FunctionPredicate(entryIntEqualsTo(EntryInt(1))).
        Or(FunctionPredicate(entryIntEqualsTo(EntryInt(3)))))
// returns EntryInt's {1,3}
```

#### ForEach

```go
total := 0

computeSumTotal := func(value interface{}) {
    total += int(value.(EntryInt).Value())
}

NewSet().
    Insert(EntryInt(1)).
    Insert(EntryInt(2)).
    Insert(EntryInt(3)).
    Stream().
    ForEach(calculateSumTotal)
// total == 6
```

### Predicates

A `Predicate` is a normal Go function which signature is

```go
func(t interface{}) bool
```

A `Predicate` has convenient pre-defined methods:
- Or
- And
- Not

Several pre-defined `Predicate`'s exist too:
- True
- False
- FunctionPredicate - a Predicate that wraps over a Function

See [example_predicate_test.go](example_predicate_test.go) for some examples.

```go
_ = ƒ.Predicate(ƒ.False).And(ƒ.Predicate(ƒ.False).Or(ƒ.True))(1) // returns false

res := ƒ.Predicate(intGreaterThanPredicate(50)).And(ƒ.True).Not()(23) // res = true

func intGreaterThanPredicate(rhs int) ƒ.Predicate {
	return func(lhs interface{}) bool {
		return lhs.(int) > rhs
	}
}
```
