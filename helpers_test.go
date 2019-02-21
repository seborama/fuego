package fuego_test

import (
	ƒ "github.com/seborama/fuego"
)

// concatenateStringsBiFunc returns a ƒ.BiFunction that
// joins 'i' and 'j' together with a '-' in between.
func concatenateStringsBiFunc(i, j ƒ.Entry) ƒ.Entry {
	iStr := i.(ƒ.EntryString)
	jStr := j.(ƒ.EntryString)
	return iStr + "-" + jStr
}

// timesTwo returns a ƒ.Function than multiplies 'i' by 2.
func timesTwo() ƒ.Function {
	return func(i ƒ.Entry) ƒ.Entry {
		return 2 * i.(ƒ.EntryInt)
	}
}

// isEvenNumberFunction is a Function that returns true when 'i' is
// an even number.
func isEvenNumberFunction(i ƒ.Entry) ƒ.Entry {
	return ƒ.EntryBool(i.(ƒ.EntryInt)&1 == 0)
}

// intGreaterThanPredicate is a Predicate for numbers greater
// than 'rhs'.
func intGreaterThanPredicate(rhs int) ƒ.Predicate {
	return func(lhs ƒ.Entry) bool {
		return int(lhs.(ƒ.EntryInt)) > rhs
	}
}
