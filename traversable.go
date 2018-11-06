package fuego

import (
	"github.com/raviqqe/hamt"
)

// Traversable is an interface for inherently recursive, multi-valued data structures.
// The order of elements is determined by Iterable.iterator(), which may vary each time
// it is called.
type Traversable interface {
	//     Traversable<T> filter(Predicate<? super T> var1);
	//     default Option<T> find(Predicate<? super T> predicate)
	Get() hamt.Entry
	Head() hamt.Entry
	HashCode() uint32
	// Traversable<T> init();
	IsEmpty() bool
	Last() hamt.Entry
	Length() int
	NonEmpty() bool
	Size() int
	Tail() Traversable
}
