# fuego - Functional Experiment in Go.

## Overview

This is a research project.

## Install

```bash
go get github.com/seborama/fuego
```

## Main features

### Set (partial implementation)

Set is based on hamt.Set and entries must implement interface `hamt.Entry`.

This is an elegant solution from [Yota Toyama](https://github.com/raviqqe) that somewhat mimics generics: the type can be anything so long as it respects the simple behaviour of `hamt.Entry`.

An example `hamt.Entry` implementation is provided in [entry_test.go](entry_test.go).

```go
// See entry_test.go for "EntryInt"
NewSet().
    Insert(EntryInt(1)).
    Insert(EntryInt(2)).
    Insert(EntryInt(3)).
    Stream()
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
```

#### Filter

This is not yet implemented. It will make use of `Predicate`'s.

## Predicates

You can create you own `Predicate`'s.

For convenience, several pre-defined `Predicate`'s are supplied:
- True
- False
- Or
- And
- Not

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
