<p align="center">
  <img src="doc/fuego_logo.png" alt="ƒuego logo" height="300">
</p>

<h3 align="center">
  <a href="#">ƒuego - Functional Experiment in Go</a>
</h3>

<p align="center">
  <a href="https://twitter.com/intent/tweet?text=Fuego%20Go%20language%20functional%20experiment&url=https://www.github.com/seborama/fuego&hashtags=golang,functional,programming,developers">
    <img src="https://img.shields.io/twitter/url/http/shields.io.svg?style=social" alt="Tweet">
  </a>
</p>

<p align="center">
  <a href="https://godoc.org/github.com/seborama/fuego">
    <img src="https://img.shields.io/badge/godoc-reference-blue.svg" alt="fuego">
  </a>
  <a href="http://goreportcard.com/report/seborama/fuego">
    <img src="https://img.shields.io/badge/go%20report-A%2B-brightgreen.svg" alt="goreportcard">
  </a>
  <a href="https://cover.run/go?tag=golang-1.9&repo=github.com%2Fseborama%2Ffuego">
    <img src="https://cover.run/go/github.com/seborama/fuego.svg?style=flat&tag=golang-1.9" alt="cover.run">
  </a>
</p>

<p align="center">
<a href="https://buymeacoff.ee/seborama/" target="_blank">
  <img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" />
  </a>
</p>

<p align="left">
  <img src="doc/fuego_code.png" alt="ƒuego example" width="654">
</p>

<p align="right">
  <img src="doc/fuego_code_employees.png" alt="ƒuego example" width="654">
</p>

## Table of content

- [Overview](#overview)
- [Documentation](#documentation)
- [Installation](#installation)
- [Example Stream](#example-stream)
- [Contributions](#contributions)
- [The Golden rules of the game](#the-golden-rules-of-the-game)
- [Pressure](#pressure)
- [Concept: Entry](#concept-entry)
- [Features summary](#features-summary)
  - [Concurrency](#concurrency)
- [Collectors](#collectors)
- [Known limitations](#known-limitations)

## [Overview](#overview)

**_Making Go come to functional programming._**

This is a research project in functional programming which I hope will prove useful!

___ƒuego___ brings a few functional paradigms to Go. The intent is to save development time while promoting code readability and reduce the risk of complex bugs.

I hope you will find it useful!

Have fun!!

[(toc)](#table-of-content)

## [Documentation](#documentation)

The code documentation and some examples can be found on [godoc](http://godoc.org/github.com/seborama/fuego).

The tests form the best source of documentation. ___ƒuego___ comes with a good collection of unit tests and testable Go examples. Don't be shy, open them up and read them and tinker with them!

> **Note:**
> <br/>
> Most tests use unbuffered channels to help detect deadlocks. In real life scenarios, it is recommended to use buffered channels for increased performance.

[(toc)](#table-of-content)

## [Installation](#installation)

### Download

```bash
go get github.com/seborama/fuego
```

Or for a specific version:

```bash
go get gopkg.in/seborama/fuego.v8
```

### Import in your code

You can import the package in the usual Go fashion.

To simplify usage, you can use an alias:

```go
package sample

import ƒ "gopkg.in/seborama/fuego.v8"
```

...or import as an unqualified dot import:

```go
package sample

import . "gopkg.in/seborama/fuego.v8"
```

[(toc)](#table-of-content)

## [Example Stream](#example-stream)

```go
    strs := EntrySlice{
        EntryString("a"),
        EntryString("bb"),
        EntryString("cc"),
        EntryString("ddd"),
    }
    
    NewStreamFromSlice(strs, 500).
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
    }

    // result: map[1:[] 2:[BB CC] 3:[DDD]]
```

[(toc)](#table-of-content)

## [Contributions](#contributions)

Contributions and feedback are welcome.

For contributions, you must develop in TDD fashion and ideally provide Go testable examples (if meaningful).

If you have an idea to improve ___ƒuego___, please share it via an issue. And if you like ___ƒuego___ give it a star to show your support for the project - it will put a smile on my face! :blush:

Thanks!!

[(toc)](#table-of-content)

## [The Golden rules of the game](#the-golden-rules-of-the-game)

1. Producers close their channel. In other words, when you create a channel, you are responsible for closing it. Similarly, whenever ___ƒuego___ creates a channel, it is responsible for closing it.

1. Consumers do not close channels.

1. Producers and consumers should be running in separate Go routines to prevent deadlocks when the channels' buffers fill up.

[(toc)](#table-of-content)

## [Pressure](#pressure)

Go channels support buffering that affects the behaviour when combining channels in a pipeline.

When the buffer of a Stream's channel of a consumer  is full, the producer will not be able to send more data through to it. This protects downstream operations from overloading.

Presently, a Go channel cannot dynamically change its buffer size. This prevents from adapting the stream flexibly. Constructs that use 'select' on channels on the producer side can offer opportunities for mitigation.

[(toc)](#table-of-content)

## [Concept: Entry](#concept-entry)

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

## [Features summary](#features-summary)

Streams:

- Stream
- IntStream
- FloatStream
- CStream - concurrent implementation of Stream

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

Check the [godoc](http://godoc.org/github.com/seborama/fuego) for full details.

[(toc)](#table-of-content)

### Concurrency

As of v8.0.0, a new concurrent model offers to process a stream concurrently while preserving order.

This is not possible yet with all Stream methods but is available with e.g. `Stream.Map`.

#### Notes on concurrency

Concurrent streams are challenging to implement owing to ordering issues in parallel processing. At the moment, the view is that the most sensible approach is to delegate control to users. Multiple ___ƒuego___ streams can be created and data distributed across as desired. This empowers users of ___ƒuego___ to implement the desired behaviour of their pipelines.

`Stream` has some methods that fan out (e.g. `ForEachC`). See the [godoc](http://godoc.org/github.com/seborama/fuego) for further information and limitations.

I recommend Rob Pike's slides on Go concurrency patterns:

- [Go Concurrency Patterns, Rob Pike, 2012](https://talks.golang.org/2012/concurrency.slide#1)

As a proof of concept and for facilitation, ___ƒuego___ has a `CStream` implementation to manage concurrently a collection of Streams.

[(toc)](#table-of-content)

## [Collectors](#collectors)

A `Collector` is a mutable reduction operation, optionally transforming the accumulated result.

Collectors can be combined to express complex operations in a concise manner.
<br/>
Simply put, a collector allows creating custom actions on a Stream.

___ƒuego___ exposes a number of functional methods such as `MapToInt`, `Head`, `LastN`, `Filter`, etc...
<br/>
Collectors also provide a few functional methods.

But... what if you need something else? And it is not straighforward or readable when combining the existing methods ___ƒuego___ offers?

Enters `Collector`: implement you own requirement functionally!
<br/>
Focus on _**what**_ needs doing in your streams (and delegate the details of the _**how**_ to the implementation of your `Collector`).

[(toc)](#table-of-content)

## [Known limitations](#known-limitations)

- several operations may be memory intensive or poorly performing.

[(toc)](#table-of-content)

<p align="center">
<a href="https://buymeacoff.ee/seborama/" target="_blank">
  <img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" />
  </a>
</p>
