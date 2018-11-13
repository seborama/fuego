package fuego

// A Maybe is a maybe monad
type Maybe struct {
	//TODO: implement! NEED MORE RESEARCH - SEARCH VAVR.OPTION
}

func NewMaybe(i interface{}) Maybe {
	return Maybe{} // TODO: review implementation when Maybe struct is defined
}

func Empty() Maybe      {}
func Of() Maybe         {}
func OfNullable() Maybe {}

func Get() Maybe {}
