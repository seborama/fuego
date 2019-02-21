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
- EntryFloat
- EntryString
- EntryMap
- EntrySlice

Check the code for additional methods each of these may provide.

#### EntryMap

This is a map of `Entry` defined as:

```go
type EntryMap map[Entry]EntrySlice
```

Stream.GroupBy uses an `EntryMap` to return data.

Note that Collector.GroupingBy offers more flexibility and can be used with `ToEntryMap` or `ToEntrySlice` for example.

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

Note that 'nil' channels are prohibited.

**NOTE:**

Concurrent streams are challenging to implement owing to ordering issues in parallel processing. At the moment, the view is that the most sensible approach is to delegate control to users. Multiple fuego streams can be created and data distributed across as desired. This empowers users of fuego to implement the desired behaviour of their pipelines.

#### Creation

When providing a Go channel to create a Stream, beware that until you close the channel, the Stream's internal Go function that processes the data on the channel will remain active. It will block until either new data is produced or the channel is closed by the producer. When a producer forgets to close the channel, the Go function will stray.

Streams created from a slice do not suffer from this issue because they are closed when the slice content is fully pushed to the Stream.

```go
ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryInt(1),
    ƒ.EntryInt(2),
    ƒ.EntryInt(3),
}, 1e3)
// or if you already have a channel of Entry:
c := make(chan ƒ.Entry) // you could add a buffer size as a second arg, if desired
go func() {
    defer close(c)
    c <- ƒ.EntryString("one")
    c <- ƒ.EntryString("two")
    c <- ƒ.EntryString("three")
    // c <- ...
}()
NewStream(c)
```

#### Filter

```go
// See helpers_test.go for "newEntryIntEqualsTo()"
s := ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryInt(1),
    ƒ.EntryInt(2),
    ƒ.EntryInt(3),
}, 0)

s.Filter(
        FunctionPredicate(entryIntEqualsTo(ƒ.EntryInt(1))).
            Or(
                FunctionPredicate(entryIntEqualsTo(ƒ.EntryInt(3)))),
)

// returns []ƒ.EntryInt{1,3}
```

#### Reduce / LeftReduce

```go
// See helpers_test.go for "concatenateStringsBiFunc()"
ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryString("four"),
    ƒ.EntryString("twelve)",
    ƒ.EntryString("one"),
    ƒ.EntryString("six"),
    ƒ.EntryString("three"),
}, 1e3).
    Reduce(concatenateStringsBiFunc)
// returns ƒ.EntryString("four-twelve-one-six-three")
```

#### ForEach

```go
total := 0

computeSumTotal := func(value ƒ.Entry) {
    total += int(value.(ƒ.EntryInt).Value())
}

s := ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryInt(1),
    ƒ.EntryInt(2),
    ƒ.EntryInt(3),
}, 0).
    ForEach(calculateSumTotal)
// total == 6
```

#### Intersperse

```go
ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryString("three"),
    ƒ.EntryString("two"),
    ƒ.EntryString("four"),
}, 1e3).
    Intersperse(ƒ.EntryString(" - "))
// "three - two - four"
```

#### GroupBy

Please refer to [stream_test.go](stream_test.go) for an example that groups numbers by parity (odd / even).

#### Count

Counts the number of elements in the Stream.

#### MapToInt

Maps this stream to an `IntStream`.

#### MapToFloat

Maps this stream to an `FloatStream`.

#### AnyMatch

Returns true if any of the elements in the stream satisfies the Predicate argument.

```go
ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryString("three"),
    ƒ.EntryString("two"),
    ƒ.EntryString("four"),
}, 1e3).
    AnyMatch(func(e ƒ.Entry) bool {
        return e.Equal(ƒ.EntryString("three"))
    })
// true
```

#### NoneMatch

Returns true if none of the elements in the stream satisfies the Predicate argument.

```go
ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryString("three"),
    ƒ.EntryString("two"),
    ƒ.EntryString("four"),
}, 1e3).
    NoneMatch(func(e ƒ.Entry) bool { return e.Equal(ƒ.EntryString("nothing like this")) })
// true
```

#### AllMatch

Returns true if every element in the stream satisfies the Predicate argument.

```go
ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryString("three"),
    ƒ.EntryString("two"),
    ƒ.EntryString("fourth"),
}, 1e3).
    AllMatch(func(e ƒ.Entry) bool {
        return strings.Contains(string(e.(ƒ.EntryString)), "t")
    })
// true
```

#### Drop

Drops the first 'n' elements of the stream and returns a new stream.

```go
ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryString("three"),
    ƒ.EntryString("two"),
    ƒ.EntryString("fourth"),
}, 1e3).
    Drop(2)
// Stream of ƒ.EntryString("fourth")
```

#### DropWhile

Drops the first elements of the stream while the predicate satisfies and returns a new stream.

```go
ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryString("three"),
    ƒ.EntryString("two"),
    ƒ.EntryString("fourth"),
}, 1e3).
    DropWhile(func(e ƒ.Entry) bool {
        return e.Equal(ƒ.EntryString("three"))
    })
// Stream of ƒ.EntryString("two") and ƒ.EntryString("fourth")
```

#### DropUntil

Drops the first elements of the stream until the predicate satisfies and returns a new stream.

```go
ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryString("three"),
    ƒ.EntryString("two"),
    ƒ.EntryString("fourth"),
}, 1e3).
    DropUntil(func(e ƒ.Entry) bool {
        return e.Equal(ƒ.EntryString("fourth"))
    })
// Stream of ƒ.EntryString("three") and ƒ.EntryString("two")
```

#### Last

Returns the last element of the stream.

This is a special case of LastN(1) which returns a single `Entry`.

```go
ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryString("three"),
    ƒ.EntryString("two"),
    ƒ.EntryString("fourth"),
}, 1e3).
    Last()
// ƒ.EntryString("fourth")
```

#### LastN

Return a slice of the last N elements of the stream.

```go
ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryString("three"),
    ƒ.EntryString("two"),
    ƒ.EntryString("fourth"),
}, 1e3).
    LastN(2)
// []ƒ.Entry{ƒ.EntryString("two"), ƒ.EntryString("fourth")}
```

#### EndsWith

Return true if the stream ends with the supplied slice of elements.

```go
ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryString("three"),
    ƒ.EntryString("two"),
    ƒ.EntryString("fourth"),
}, 1e3).
    EndsWith([]ƒ.Entry{ƒ.EntryString("two"), ƒ.EntryString("fourth")})
// true
```

#### Head

Returns the first `Entry` of the stream.

This is a special case of HeadN(1) which returns a single `Entry`.

```go
ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryString("three"),
    ƒ.EntryString("two"),
    ƒ.EntryString("fourth"),
}, 1e3).
    Head()
// ƒ.EntryString("three")
```

#### HeadN

Return a slice of the first N elements of the stream.

```go
ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryString("three"),
    ƒ.EntryString("two"),
    ƒ.EntryString("fourth"),
}, 1e3).
    HeadN(2)
// []ƒ.Entry{ƒ.EntryString("three"), ƒ.EntryString("two")}
```

#### StartsWith

Return true if the stream starts with the supplied slice of elements.

```go
ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryString("three"),
    ƒ.EntryString("two"),
    ƒ.EntryString("fourth"),
}, 1e3).
    StartsWith([]ƒ.Entry{ƒ.EntryString("three"), ƒ.EntryString("two")})
// true
```

#### Take

Takes the first 'n' elements of the stream and returns a new stream.

```go
ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryString("three"),
    ƒ.EntryString("two"),
    ƒ.EntryString("fourth"),
}, 1e3).
    Take(2)
// Stream of []ƒ.Entry{ƒ.EntryString("three"), ƒ.EntryString("two")}
```

#### TakeWhile

Takes the first elements of the stream while the predicate satisfies and returns a new stream.

```go
ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryString("three"),
    ƒ.EntryString("two"),
    ƒ.EntryString("fourth"),
}, 1e3).
    TakeWhile(func(e ƒ.Entry) bool {
        return strings.HasPrefix(string(e.(ƒ.EntryString)), "t")
    })
// Stream of []ƒ.Entry{ƒ.EntryString("three"), ƒ.EntryString("two")}
```

#### TakeUntil

Takes the first elements of the stream until the predicate satisfies and returns a new stream.

```go
ƒ.NewStreamFromSlice([]ƒ.Entry{
    ƒ.EntryString("three"),
    ƒ.EntryString("two"),
    ƒ.EntryString("fourth"),
}, 1e3).
    TakeUntil(func(e ƒ.Entry) bool {
        return e.Equal(ƒ.EntryString("fourth"))
    })
// Stream of []ƒ.Entry{ƒ.EntryString("three"), ƒ.EntryString("two")}
```

#### Collect

Applies a `Collector` to this Stream.

It should be noted that this method returns an `interface{}` which enables it to return `Entry` as well as any other Go types.

Example:

```go
    strs := EntrySlice{
        EntryString("a"),
        EntryString("bb"),
        EntryString("cc"),
        EntryString("ddd"),
    }

    NewStreamFromSlice(strs, 1e3).
        Collect(
            GroupingBy(
                stringLength,
                Mapping(
                    stringToUpper,
                    Filtering(
                        stringLengthGreaterThan(1),
                        ToEntrySlice()))))
    // map[1:[] 2:[BB CC] 3:[DDD]]
```

### IntStream

A Stream of EntryInt.

It contains all of the methods `Stream` exposes and additional methods that pertain to an `EntryInt` stream such as aggregate functions (`Sum()`, `Average()`, etc).

Note that 'nil' channels are prohibited.

**NOTE:**
The current implementation is based on `Stream` and an intermediary channel that converts incoming `EntryInt` elements to `Entry`. This approach offers programming conciseness but the use of an intermediary channel likely decreases performance. This also means that type checking is weak on methods "borrowed" from `Stream` that expect `Entry` (instead of `EntryInt`).

#### Stream methods

All methods that pertain to `Stream` are available to `IntStream`.

#### Max

Returns the greatest element in the stream.

#### Min

Returns the smallest element in the stream.

#### Sum

Returns the sum of all elements in the stream.

#### Average

Returns the average of all elements in the stream.

### FloatStream

A Stream of `EntryFloat`. It is akin to `IntStream` but for `EntryFloat`'s.

Note that 'nil' channels are prohibited.

It contains all of the methods `Stream` exposes and additional methods that pertain to an `EntryFloat` stream such as aggregate functions (`Sum()`, `Average()`, etc).

**NOTE:**
The current implementation is based on `Stream` and an intermediary channel that converts incoming `EntryFloat` elements to `Entry`. This approach offers programming conciseness but the use of an intermediary channel likely decreases performance. This also means that type checking is weak on methods "borrowed" from `Stream` that expect `Entry` (instead of `EntryFloat`).

#### Stream methods

All methods that pertain to `Stream` are available to `FloatStream`.

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
- Xor

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
        Negate()(ƒ.EntryInt(23)) // res = true

func intGreaterThanPredicate(rhs int) ƒ.Predicate {
    return func(lhs ƒ.Entry) bool {
        return int(lhs.(ƒ.EntryInt)) > rhs
    }
}
```

### Collector

A `Collector` is a mutable reduction operation, optionally transforming the accumulated result.

Collectors can be combined to express complex operations in a concise manner.

It should be noted that the `finisher` function is optional (i.e. it may acceptably be `nil`).

Example:

```go
    strs := EntrySlice{
        EntryString("a"),
        EntryString("bb"),
        EntryString("cc"),
        EntryString("ddd"),
    }

    NewStreamFromSlice(strs, 1e3).
        Collect(
            GroupingBy(
                stringLength,
                Mapping(
                    stringToUpper,
                    ToEntryMap())))
// map[1:[A] 2:[BB CC] 3:[DDD]]
```

#### Available collectors

- GroupingBy:

  ```go
  GroupingBy(classifier Function, collector Collector) Collector
  ```

- Mapping:

  ```go
  Mapping(mapper Function, collector Collector) Collector
  ```

- Filtering:

  ```go
  Mapping(mapper Function, collector Collector) Collector
  ```

- ToEntrySlice:

  ```go
  ToEntrySlice() Collector
  ```

#### Available finishers

- IdentityFinisher

  This is a basic finisher that returns its input unchanged.

## Known limitations

- several operations may be memory intensive or poorly performing.
