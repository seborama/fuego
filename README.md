# fuego - Functional Experiment in Go.

## Overview

_Making Go come._

This is a research project.

## Install

```bash
go get github.com/seborama/fuego
```

## Main features

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

See [example_map_test.go](example_map_test.go) for more details of an example of Map with Stream and Filter combined together to extract entries which key are an even number.

### Function

```go
type functionTimesTwo int

func newFunctionTimesTwo() functionTimesTwo {
	return *(new(functionTimesTwo))
}

func (f functionTimesTwo) Apply(i interface{}) interface{} {
	num := i.(EntryInt).Value()
	return interface{}(2 * num)
}

f := newFunctionTimesTwo()
f.Apply(7) // returns EntryInt 7
```

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
// See stream_test.go for "NewFunctionTimesTwo()"
NewSet().
    Insert(EntryInt(1)).
    Insert(EntryInt(2)).
    Insert(EntryInt(3)).
    Stream().
    Map(NewFunctionTimesTwo())
// returns EntryInt's {2,4,6}
```

#### Filter

```go
// See stream_test.go for "NewFunctionTimesTwo()"
NewSet().
    Insert(EntryInt(1)).
    Insert(EntryInt(2)).
    Insert(EntryInt(3)).
    Stream().
    Filter(Or(
        NewFunctionPredicate(newEntryIntEqualsTo(EntryInt(1))),
        NewFunctionPredicate(newEntryIntEqualsTo(EntryInt(3))),
    ))
// returns EntryInt's {1,3}
```

### Predicates

You can create you own `Predicate`'s.

For convenience, several pre-defined `Predicate`'s are supplied:
- True
- False
- Or
- And
- Not
- Function - a Predicate that wraps over a Function

```go
type intGreaterThanPredicate struct {
	number int
}

func newIntGreaterThanPredicate(number int) intGreaterThanPredicate {
	return intGreaterThanPredicate{
		number: number,
	}
}

func (p intGreaterThanPredicate) Test(t interface{}) bool {
	return t.(int) > p.number
}

newIntGreaterThanPredicate(2).Test(7) // returns false
Not(newIntGreaterThanPredicate(2)).Test(7) // returns true
```
