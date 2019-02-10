# fuego - Functional Experiment in Go

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
go get gopkg.in/seborama/fuego.v6
```

## Contribute

Contributions and feedback are welcome.

As this project is still in early stages, large portions of the code get crafted and thrown away.

For contributions, you must develop in TDD fashion and ideally provide Go testable examples (if meaningful).

## The Golden rules of the game

1- Producers close their channel. In other words, when you create a channel, you are responsible for closing it. Similarly, whenever fuego creates a channel, it is responsible for closing it.

2- Consumers do not close channels.

3- Producers and consumers should be running in separate Go routines to prevent deadlocks when the channels' buffers fill up.

## Pressure

Go channels support buffering that affect the behaviour when combining channels in a pipeline.

When a consumer Stream's channel buffer is full, the producer will not be able to send more data through to it. This protects downstream operations from overloading.

Presently, a Go channel cannot dynamically change its buffer size. This prevents from adapting the stream flexibly. Constructs that use 'select' on channels on the producer side can offer opportunities for mitigation.

## Main features

The code documentation can be found on [godoc](http://godoc.org/github.com/seborama/fuego).

**The tests form the best source of documentation. Fuego comes with a good collection of unit tests and testable Go examples. Don't be shy, open them up and read them and tinker with them!**

**Note however that most tests use unbuffered channels to help detect deadlocks. On real life scenarios, it is recommended to use buffered channels for increased performance.** 

Have fun!!

### Entry

`Entry` is inspired by `hamt.Entry`. This is an elegant solution from [Yota Toyama](https://github.com/raviqqe): the type can be anything so long as it respects the simple behaviour of `hamt.Entry`. This provides an abstraction of types yet with known behaviour:

- Hash(): identifies an Entry Uniquely.
- Equal(): defines equality for a type of `Entry`. `Equal()` is expected to be based on `Hash()`.

Several Entry implementations are provided:

- EntryBool
- EntryInt
- EntryMap
- EntrySlice

#### EntryMap

This is a map of `Entry` defined as:

```go
type EntryMap map[Entry]EntrySlice
```

GroupBy methods use an `EntryMap` to return data.

It is important to remember that maps are **not** ordered.

#### EntrySlice

This is an ordered slice of `Entry` elements which signature is:

```go
type EntrySlice []Entry
```

### Maybe

A `Maybe` represents an optional value.

When the value is `nil`, `Maybe` is considered empty unless it was created with `MaybeSome()` in which case it is considered to hold the `nil` value.

`MaybeNone()` always produces an empty optional.

### Tuple

fuego provides these `Tuple`'s:

- Tuple0
- Tuple1
- Tuple2

The values of fuego `Tuples` are  of type `Entry`.

### Consumer

Consumer is a kind of side-effect function that accepts one argument and does not
return any value.

```go
type Consumer func(i Entry)
```

### Functions

See [example_function_test.go](example_function_test.go) for basic example uses of `Function` and `BiFunction` and the other tests / examples for more uses.

#### Function

A `Function` is a normal Go function which signature is:

```go
func(i Entry) Entry
```

#### BiFunction

A `BiFunction` is a normal Go function which signature is:

```go
func(i,j Entry) Entry
```

`BiFunction`'s are used with Stream.Reduce() for instance, as seen in [stream_test.go](stream_test.go).

#### ToIntFunction

This is a special case of Function used to convert a Stream to an IntStream.

```go
type ToIntFunction func(e Entry) EntryInt
```

### Stream

A Stream is a wrapper over a Go channel.

**NOTE:**

Concurrent streams are challenging to implement owing to ordering issues in parallel processing. At the moment, the view is that the most sensible approach is to delegate control to users. Multiple fuego streams can be created and data distributed across as desired. This empowers users of fuego to implement the desired behaviour of their pipelines.

#### Creation

When providing a Go channel to create a Stream, beware that until you close the channel, the Stream's internal Go function that processes the data on the channel will remain active. It will block until either new data is produced or the channel is closed by the producer. When a producer forgets to close the channel, the Go function will stray.

Streams created from a slice do not suffer from this issue because they are closed when the slice content is fully pushed to the Stream.

```go
ƒ.NewStreamFromSlice([]Entry{
    EntryInt(1),
    EntryInt(2),
    EntryInt(3),
}, 1e3)
// or if you already have a channel of Entry:
c := make(chan Entry) // you could add a buffer size as a second arg, if desired
go func() {
    defer close(c)
    c <- EntryString("one")
    c <- EntryString("two")
    c <- EntryString("three")
    // c <- ...
}()
NewStream(c)
```

#### Filter

```go
// See helpers_test.go for "newEntryIntEqualsTo()"
s := ƒ.NewStreamFromSlice([]Entry{
    EntryInt(1),
    EntryInt(2),
    EntryInt(3),
}, 0)

s.Filter(
        FunctionPredicate(entryIntEqualsTo(EntryInt(1))).
            Or(
                FunctionPredicate(entryIntEqualsTo(EntryInt(3)))),
)

// returns []EntryInt{1,3}
```

#### Reduce / LeftReduce

```go
// See helpers_test.go for "concatenateStringsBiFunc()"
ƒ.NewStreamFromSlice([]Entry{
    EntryString("four"),
    EntryString("twelve)",
    EntryString("one"),
    EntryString("six"),
    EntryString("three"),
}, 1e3).
    Reduce(concatenateStringsBiFunc)
// returns EntryString("four-twelve-one-six-three")
```

#### ForEach

```go
total := 0

computeSumTotal := func(value interface{}) {
    total += int(value.(EntryInt).Value())
}

s := ƒ.NewStreamFromSlice([]Entry{
    EntryInt(1),
    EntryInt(2),
    EntryInt(3),
}, 0).
    ForEach(calculateSumTotal)
// total == 6
```

#### Intersperse

```go
ƒ.NewStreamFromSlice([]Entry{
    EntryString("three"),
    EntryString("two"),
    EntryString("four"),
}, 1e3).
    Intersperse(EntryString(" - "))
// "three - two - four"
```

#### GroupBy

Please refer to [stream_test.go](stream_test.go) for an example that groups numbers by parity (odd / even).

#### Count

Counts the number of elements in the Stream.

#### Close

Closes the Stream. It cannot receive more data but can continue consuming buffered messages.

### IntStream

A Stream of EntryInt.

It contains all of the methods Stream exposes and additional methods that pertain to an `EntryInt` stream such as aggregate functions (`Sum()`, `Average()`, etc).

The current implementation is based on `Stream` and an intermediary channel that converts incoming `EntryInt` elements to `Entry`. This approach offers programming conciseness but the use of an intermediary channel likely decreases performance.

#### Max

Returns the greatest element in the stream.

#### Min

Returns the smallest element in the stream.

#### Sum

Returns the sum of all elements in the stream.

#### Average

Returns the average of all elements in the stream.

### Predicates

A `Predicate` is a normal Go function which signature is:

```go
type Predicate func(t Entry) bool
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
    _ = ƒ.Predicate(ƒ.False).
        ƒ.And(Predicate(ƒ.False).
            ƒ.Or(ƒ.True))(ƒ.EntryInt(1)) // returns false

res := ƒ.Predicate(intGreaterThanPredicate(50)).
        And(ƒ.True).
        Not()(ƒ.EntryInt(23)) // res = true

func intGreaterThanPredicate(rhs int) ƒ.Predicate {
    return func(lhs ƒ.Entry) bool {
        return int(lhs.(ƒ.EntryInt)) > rhs
    }
}
```

## Known limitations

- several operations may be memory intensive or poorly performing.
