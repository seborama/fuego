package fuego

// Useful constants that represent native go types.
const (
	Bool       = false
	Int        = int(0)
	Int8       = int8(0)
	Int16      = int16(0)
	Int32      = int32(0)
	Int64      = int64(0)
	Uint       = uint(0)
	Uint8      = uint8(0)
	Uint16     = uint16(0)
	Uint32     = uint32(0)
	Uint64     = uint64(0)
	Float32    = float32(0)
	Float64    = float64(0)
	Complex64  = complex64(0)
	Complex128 = complex128(0)
	String     = ""
)

// Useful variables that represent native go types.
var (
	SBool       = []bool{}
	SInt        = []int{}
	SInt8       = []int8{}
	SInt16      = []int16{}
	SInt32      = []int32{}
	SInt64      = []int64{}
	SUint       = []uint{}
	SUint8      = []uint8{}
	SUint16     = []uint16{}
	SUint32     = []uint32{}
	SUint64     = []uint64{}
	SFloat32    = []float32{}
	SFloat64    = []float64{}
	SComplex64  = []complex64{}
	SComplex128 = []complex128{}
	SString     = []string{}
)

// Useful variables that represent native go types.
var (
	BoolPtr       = ptr(Bool)
	IntPtr        = ptr(Int)
	Int8Ptr       = ptr(Int8)
	Int16Ptr      = ptr(Int16)
	Int32Ptr      = ptr(Int32)
	Int64Ptr      = ptr(Int64)
	UintPtr       = ptr(Uint)
	Uint8Ptr      = ptr(Uint8)
	Uint16Ptr     = ptr(Uint16)
	Uint32Ptr     = ptr(Uint32)
	Uint64Ptr     = ptr(Uint64)
	Float32Ptr    = ptr(Float32)
	Float64Ptr    = ptr(Float64)
	Complex64Ptr  = ptr(Complex64)
	Complex128Ptr = ptr(Complex128)
	StringPtr     = ptr(String)
)

// Useful variables that represent native go types.
var (
	SBoolPtr       = []*bool{}
	SIntPtr        = []*int{}
	SInt8Ptr       = []*int8{}
	SInt16Ptr      = []*int16{}
	SInt32Ptr      = []*int32{}
	SInt64Ptr      = []*int64{}
	SUintPtr       = []*uint{}
	SUint8Ptr      = []*uint8{}
	SUint16Ptr     = []*uint16{}
	SUint32Ptr     = []*uint32{}
	SUint64Ptr     = []*uint64{}
	SFloat32Ptr    = []*float32{}
	SFloat64Ptr    = []*float64{}
	SComplex64Ptr  = []*complex64{}
	SComplex128Ptr = []*complex128{}
	SStringPtr     = []*string{}
)

// Comparable is a constraint that matches any type that supports the operators:
// >= <= > < == !=
type Comparable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		string
}

// Addable is a constraint that matches any type that supports the operator '+'.
type Addable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		~complex64 | ~complex128 |
		string
}

// Max is a BiFunction that returns the greatest of 2 values.
func Max[T Comparable](n1, n2 T) T {
	if n1 > n2 {
		return n1
	}
	return n2
}

// Min is a BiFunction that returns the smallest of 2 values.
func Min[T Comparable](n1, n2 T) T {
	if n1 < n2 {
		return n1
	}
	return n2
}

// Sum is a BiFunction that returns the sum of 2 values.
func Sum[T Addable](n1, n2 T) T {
	return n1 + n2
}

func ptr[T any](t T) *T { return &t }
