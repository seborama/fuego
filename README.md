# fuego - Functional Experiment in Go

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/seborama/fuego) [![goreportcard](https://img.shields.io/badge/go%20report-A%2B-brightgreen.svg)](http://goreportcard.com/report/seborama/fuego) [![cover.run](https://cover.run/go/github.com/seborama/fuego.svg?style=flat&tag=golang-1.9)](https://cover.run/go?tag=golang-1.9&repo=github.com%2Fseborama%2Ffuego)

<p align="center">
  <img src="doc/fuego_logo.png" alt="ƒuego logo" height="300">
</p>

## Table of content

- [Overview](#overview)
- [Documentation](#documentation)
- [Installation](#installation)
- [Contributions](#contributions)
- [The Golden rules of the game](#the-golden-rules-of-the-game)
- [Pressure](#pressure)
- [Concept: Entry](#concept-entry)
- [Example Stream](#example-stream)
- [Features summary](#features-summary)
  - [Concurrency](#concurrency)
- [Collectors](#collectors)
- [Known limitations](#known-limitations)

## Overview

_Making Go come to functional programming._

<p align="left">
  <img src="doc/fuego_code.png" alt="ƒuego example" width="654">
</p>


This is a research project in functional programming which I hope will prove useful!

Fuego brings a few functional paradigms to Go. The intent is to save development time while promoting code readability and reduce the risk of complex bugs.

Have fun!!

[(toc)](#table-of-content)

## Documentation

The code documentation and some examples can be found on [godoc](http://godoc.org/github.com/seborama/fuego).

**The tests form the best source of documentation. Fuego comes with a good collection of unit tests and testable Go examples. Don't be shy, open them up and read them and tinker with them!**

**Note however that most tests use unbuffered channels to help detect deadlocks. In real life scenarios, it is recommended to use buffered channels for increased performance.**

[(toc)](#table-of-content)

## Installation

```bash
go get github.com/seborama/fuego
```

Or for a specific version:

```bash
go get gopkg.in/seborama/fuego.v7
```

[(toc)](#table-of-content)

## Contributions

Contributions and feedback are welcome.

For contributions, you must develop in TDD fashion and ideally provide Go testable examples (if meaningful).

If you have ideas to improve **fuego**, please share them via an issue. And if you like **fuego** give it a star to show your support for the project - it is my greatest reward! :blush:

Thanks! 

[(toc)](#table-of-content)

## The Golden rules of the game

1. Producers close their channel. In other words, when you create a channel, you are responsible for closing it. Similarly, whenever **fuego** creates a channel, it is responsible for closing it.

1. Consumers do not close channels.

1. Producers and consumers should be running in separate Go routines to prevent deadlocks when the channels' buffers fill up.

[(toc)](#table-of-content)

## Pressure

Go channels support buffering that affects the behaviour when combining channels in a pipeline.

When the buffer of a Stream's channel of a consumer  is full, the producer will not be able to send more data through to it. This protects downstream operations from overloading.

Presently, a Go channel cannot dynamically change its buffer size. This prevents from adapting the stream flexibly. Constructs that use 'select' on channels on the producer side can offer opportunities for mitigation.

[(toc)](#table-of-content)

## Concept: Entry

`Entry` is inspired by `hamt.Entry`. This is an elegant solution from [Yota Toyama](https://github.com/raviqqe): the type can be anything so long as it respects the simple behaviour of the`Entry` interface. This provides an abstraction of types yet with known behaviour:

- Hash(): identifies an Entry Uniquely.
- Equal(): defines equality for a concrete type of `Entry`. `Equal()` is expected to be based on `Hash()` for non-basic types. Equal should ensure the compared Entry is of the same type as the reference Entry. For instance, `EntryBool(false)` and `EntryInt(0)` both have a Hash of `0`, yet they aren't equal.

Several Entry implementations are provided:

- EntryBool
- EntryInt
- EntryFloat
- EntryString
- EntryMap
- EntrySlice
- Tuples

Check the [godoc](http://godoc.org/github.com/seborama/fuego) for additional methods each of these may provide.

[(toc)](#table-of-content)

## Example Stream

```go
    strs := EntrySlice{
        EntryString("a"),
        EntryString("bb"),
        EntryString("cc"),
        EntryString("ddd"),
    }
    got := NewStreamFromSlice(strs, 500).
        Filter(isEntryString).
        Distinct().
        Collect(
            GroupingBy(
                stringLength,
                Mapping(
                    stringToUpper,
                    Filtering(
                        stringLengthGreaterThan(1),
                        ToEntrySlice()))))

    // result: map[1:[] 2:[BB CC] 3:[DDD]]
```

[(toc)](#table-of-content)

## Features summary

Streams:

- Stream
- IntStream
- FloatStream

Functional Types:

- Maybe
- Tuple
- Predicate:
  - True
  - False
  - FunctionPredicate

Functions:

- Consumer
- Function:
  - ToIntFunction
  - ToFloatFunction
- BiFunction
- StreamFunction:
  - FlattenEntrySliceToEntry
- Predicate:
  - Or
  - Xor
  - And
  - Not / Negate

Collectors:

- GroupingBy
- Mapping
- FlatMapping
- Filtering
- Reducing
- ToEntrySlice

[(toc)](#table-of-content)

### Concurrency

Concurrent streams are challenging to implement owing to ordering issues in parallel processing. At the moment, the view is that the most sensible approach is to delegate control to users. Multiple **fuego** streams can be created and data distributed across as desired. This empowers users of **fuego** to implement the desired behaviour of their pipelines.

`Stream` has some methods that fan out (e.g. `ForEachC`). See the [godoc](http://godoc.org/github.com/seborama/fuego) for further information and limitations.

I recommend Rob Pike's slides on Go concurrency patterns:

- [Go Concurrency Patterns, Rob Pike, 2012](https://talks.golang.org/2012/concurrency.slide#1)

[(toc)](#table-of-content)

## Collectors

A `Collector` is a mutable reduction operation, optionally transforming the accumulated result.

Collectors can be combined to express complex operations in a concise manner.

In other words, a collector allows creating custom actions on a Stream. **fuego** comes shipped with a number of methods such as `MapToInt`, `Head`, `LastN`, `Filter`, etc, and Collectors also provide a few additional methods. But what if you need something else? And it is not straighforward or readable when combining the existing methods **fuego** offers? Enters `Collector`: implement you own requirement functionally! Focus on _**what**_ needs to be done in your streams (and delegate the details of the _**how**_ to the implementation of your `Collector`).

It should be noted that the `finisher` function is optional (i.e. it may acceptably be `nil`).

[(toc)](#table-of-content)

## Known limitations

- several operations may be memory intensive or poorly performing.

[(toc)](#table-of-content)
