# fuego - Functional Experiment in Go.

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/seborama/fuego) [![goreportcard](https://img.shields.io/badge/go%20report-A%2B-brightgreen.svg)](http://goreportcard.com/report/seborama/fuego) [![cover.run](https://cover.run/go/github.com/seborama/fuego.svg?style=flat&tag=golang-1.9)](https://cover.run/go?tag=golang-1.9&repo=github.com%2Fseborama%2Ffuego) 

## Overview

_Making Go come._

This is a research project in functional programming which I hope will prove useful!

Fuego brings a few functional paradigms to Go. The intent is to save development time while promoting code readability and reduce the risk of complex bugs.



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

**The tests form the best source of documentation. Fuego comes with a good collection of unit tests and testable Go examples. Don't be shy, open them up and read them and tinker with them!**

Have fun!!

### Entry

`Entry` is based on `hamt.Entry`. This is an elegant solution from [Yota Toyama](https://github.com/raviqqe): the type can be anything so long as it respects the simple behaviour of `hamt.Entry`. This provides an abstraction of types yet with known behaviour.

### Maybe

A `Maybe` represents an optional value.

When the value is `nil`, `Maybe` is considered empty unless it was created with `MaybeSome()` in which case it is considered to hold the `nil` value.

`MaybeNone()` always produces an empty optional.

### HamtSet

`HamtSet` is based on hamt.Set and entries must implement interface `Entry`.

It is an unordered collection, yet unnaturally sorted (i.e. sorting is based on the hash of the `hamt.Entry`).

An example `Entry` implementation called `EntryInt` is provided in [entry_test.go](entry_test.go).

```go
// See entry_test.go for "EntryInt"
NewHamtSet().
    Insert(EntryInt(1)).
    Insert(EntryInt(2)).
    Delete(EntryInt(1)).
    Insert(EntryInt(3)).
    Stream()
```

Uses of streams with Sets are also available in [example_map_test.go](example_map_test.go).

### OrderedSet

This is an **ordered** implementation of a `Set`.

### HamtMap

As with `HamtSet`, `HamtMap` is based on hamt.Map and entries must implement interface `Entry` for its keys but values can be anything (`interface{}`).

It is an unordered collection, yet unnaturally sorted (i.e. sorting is based on the hash of the `hamt.Entry` keys).

See [example_map_test.go](example_map_test.go) for more details of an example of `HamtMap` with Stream and Filter combined together to extract entries which keys are an even number.

### OrderedMap

This is an **ordered** implementation of a `Map`.

### Tuple

fuego provides these `Tuple`'s:
- Tuple0
- Tuple1
- Tuple2

The values of fuego `Tuples` is of type `Entry` but can represent any type (see EntryInt and EntryString examples).

### Functions

See [example_function_test.go](example_function_test.go) for basic example uses of `Function` and `BiFunction` and the other tests / examples for more uses.

#### Function
A `Function` is a normal Go function which signature is

```go
func(i Entry) Entry
```

#### BiFunction

A `BiFunction` is a normal Go function which signature is

```go
func(i,j Entry) Entry
```

`BiFunction`'s are used with Stream.Reduce() for instance, as seen in [stream_test.go](stream_test.go).

### Iterator

You can create you own `Iterator`'s.

See [iterator.go](iterator.go) for several convenience implementations of iterators:
- NewSliceIterator
- NewHamtSetIterator
- NewOrderedSetIterator

Iterator supports:
- Forward
- Value
- Reverse
- Size

```go
NewSliceIterator([]Entry{EntryInt(2), EntryInt(3)}) // returns an Iterator over []interface{2, 3}

NewHamtSetIterator(NewHamtSet().
    Insert(EntryInt(2))), // returns an Iterator over a HamtSet that contains a single EntryInt(2)
```

### Stream

#### Creation

```go
someGoSlice := []int{1, 2, 3}
NewStream(
    NewSliceIterator(someGoSlice)),
```

```go
NewStream(
    NewSetIterator(
        NewOrderedSet().
            Insert(EntryInt(1)).
            Insert(EntryInt(2))))
```

#### Map

```go
// See in this README and in helpers_test.go for "functionTimesTwo()"
NewOrderedSet().
    Insert(EntryInt(1)).
    Insert(EntryInt(2)).
    Insert(EntryInt(3)).
    Stream().
    Map(functionTimesTwo())
// returns EntryInt's {2,4,6}
```

#### Filter

```go
// See helpers_test.go for "newEntryIntEqualsTo()"
NewOrderedSet().
    Insert(EntryInt(1)).
    Insert(EntryInt(2)).
    Insert(EntryInt(3)).
    Stream().
    Filter(FunctionPredicate(entryIntEqualsTo(EntryInt(1))).
        Or(FunctionPredicate(entryIntEqualsTo(EntryInt(3)))))
// returns EntryInt's {1,3}
```

#### Reduce / LeftReduce / RightReduce

```go
// See helpers_test.go for "concatenateStringsBiFunc()"
NewOrderedSet().
    Insert(EntryString("four")).
    Insert(EntryString("twelve")).
    Insert(EntryString("one")).
    Insert(EntryString("six")).
    Insert(EntryString("three"))
    Stream().
    Reduce(concatenateStringsBiFunc)
// returns EntryString("one-three-twelve-six-four")
```

#### ForEach

```go
total := 0

computeSumTotal := func(value interface{}) {
    total += int(value.(EntryInt).Value())
}

NewOrderedSet().
    Insert(EntryInt(1)).
    Insert(EntryInt(2)).
    Insert(EntryInt(3)).
    Stream().
    ForEach(calculateSumTotal)
// total == 6
```

#### Intersperse

```go
NewOrderedSet().
    Insert(EntryString("three")).
    Insert(EntryString("two")).
    Insert(EntryString("four")).
    Stream().
    Intersperse(EntryString(" - "))
// "three - two - four"
```

#### GroupBy

Please refer to [stream_test.go](stream_test.go) for an example that groups numbers by parity (odd / even).

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
// ƒ is ALT+f on Mac. For other OSes, search the internet,  for instance,  this page: https://en.wikipedia.org/wiki/%C6%91#Appearance_in_computer_fonts
_ = ƒ.Predicate(ƒ.False).And(ƒ.Predicate(ƒ.False).Or(ƒ.True))(1) // returns false

res := ƒ.Predicate(intGreaterThanPredicate(50)).And(ƒ.True).Not()(23) // res = true

func intGreaterThanPredicate(rhs int) ƒ.Predicate {
	return func(lhs interface{}) bool {
		return lhs.(int) > rhs
	}
}
```

## Known limitations

- hamt.Set and hamt.Map are not ordered as per their initialisation but rather following their Hash. Use OrderedSet as an alternative.
- several operations may be memory intensive or poorly performing, notably - but not limited to - in OrderedSet.
