package golang_args

/*
Entity combines name, pseudonyms, necessariness and callback functions
*/
type argEntity struct {
	ArgName            string
	AssociatedFunction func(...string)
	Pseudonyms         []string
	NargsChecker       func(uint) bool
	IsNecessary        bool
}

func constructEntityChecker(function func(...string), name string, nargsChecker func(uint) bool, isNec bool, pseudonyms ...string) argEntity {
	return argEntity{name, function, pseudonyms, nargsChecker, isNec}
}
