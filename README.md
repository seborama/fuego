# fuego - Functional Experiment in Go.

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/seborama/fuego) [![goreportcard](https://img.shields.io/badge/go%20report-A%2B-brightgreen.svg)](http://goreportcard.com/report/seborama/fuego) [![cover.run](https://cover.run/go/github.com/seborama/fuego.svg?style=flat&tag=golang-1.9)](https://cover.run/go?tag=golang-1.9&repo=github.com%2Fseborama%2Ffuego) 

## Overview

_Making Go come to functional programming._

This is a research project in functional programming which I hope will prove useful!

Fuego brings a few functional paradigms to Go. The intent is to save development time while promoting code readability and reduce the risk of complex bugs.

## Install

```bash
go get github.com/seborama/fuego
```

Or for a specific version:

```bash
go get gopkg.in/seborama/fuego.v5
```

## Contribute

Contributions and feedback are welcome.

For contributions, you must develop in TDD fashion and ideally provide Go testable examples (if meaningful).

## Main features

The code documentation can be found on [godoc](http://godoc.org/github.com/seborama/fuego).

**The tests form the best source of documentation. Fuego comes with a good collection of unit tests and testable Go examples. Don't be shy, open them up and read them and tinker with them!**

Have fun!!

### Entry

`Entry` is inspired by `hamt.Entry`. This is an elegant solution from [Yota Toyama](https://github.com/raviqqe): the type can be anything so long as it respects the simple behaviour of `hamt.Entry`. This provides an abstraction of types yet with known behaviour:

- Hash(): identifies an Entry Uniquely.
- Equal(): defines equality for a type of `Entry`. `Equal()` is expected to be based on `Hash()`.

Several Entry implementations are provided:

- EntryBool
- EntryMap
- EntrySlice

#### EntryMap

TODO: TBC

#### EntrySlice

TODO: TBC

### Maybe

A `Maybe` represents an optional value.

When the value is `nil`, `Maybe` is considered empty unless it was created with `MaybeSome()` in which case it is considered to hold the `nil` value.

`MaybeNone()` always produces an empty optional.

### Tuple

fuego provides these `Tuple`'s:

- Tuple0
- Tuple1
- Tuple2

The values of fuego `Tuples` is of type `Entry` but can represent any type (see EntryInt and EntryString examples).

### Consumer

TODO: TBC

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

### Stream

A Stream is a wrapper over a Go channel.

**NOTE**
At present, the Go channel is bufferred. This poses ordering issues in parallel processing. It is likely that in a future release this will change. One option is a slice of channels. This improves parallelism but still poses issues with some scenarios where the continuity of values matters (e.g. calculating Fibonacci sequences) as with Reduce, ...

#### Creation

When providing a Go channel to create a Stream, beware that until you close the channel, the Stream's internal Go function that processes the Stream will remain active. This can lead to a stray Go function.

```go
ƒ.NewStreamFromSlice([]int{1, 2, 3})
// or if you already have a channel of Entry:
c := make(chan Entry, 1e3)
defer close(c)
c <- EntryString("one")
// c <- ...
NewStream(c)
```

#### Filter

```go
// See helpers_test.go for "newEntryIntEqualsTo()"
s := ƒ.NewStreamFromSlice([]int{1, 2, 3})
s.Filter(FunctionPredicate(entryIntEqualsTo(EntryInt(1))).
    Or(FunctionPredicate(entryIntEqualsTo(EntryInt(3)))))
// returns []EntryInt{1,3}
```

#### Reduce / LeftReduce

```go
// See helpers_test.go for "concatenateStringsBiFunc()"
ƒ.NewStreamFromSlice(string{
    "four",
    "twelve",
    "one",
    "six",
    "three",
}).
Reduce(concatenateStringsBiFunc)
// returns EntryString("four-twelve-one-six-three")
```

#### ForEach

```go
total := 0

computeSumTotal := func(value interface{}) {
    total += int(value.(EntryInt).Value())
}

ƒ.NewStreamFromSlice([]int{1, 2, 3}).
ForEach(calculateSumTotal)
// total == 6
```

#### Intersperse

```go
ƒ.NewStreamFromSlice([]string{
    "three",
    "two",
    "four",
}).
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

- several operations may be memory intensive or poorly performing.
