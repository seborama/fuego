// package fuego provides various functional facilities.
//
//////////////////
// Important note:
//////////////////
//
// Go does not yet support parameterised methods:
// https://go.googlesource.com/proposal/+/master/design/43651-type-parameters.md#no-parameterized-methods
//
// The below construct is not currently possible:
// func (s Stream[T]) Map[R any](mapper Function[T, R]) Stream[R] {...}
//                       ^^^^^^^ no!
//
// One option would be to make `Map` a function rather than a method but constructs would be chained
// right-to-left rather than left-to-right, which I think is awkward.
// Example: "Map(Map(stream,f1),f2)" instead of "stream.Map(f1).Map(f2)".
//
// A syntactically lighter approach is provided with `SC`` and `C``.
// See functions `SC`` and `C `for casting Stream[R] to a typed Stream[T any].
package fuego
