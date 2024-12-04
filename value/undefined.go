package value

type UndefinedType struct{}

// Elk's internal undefined value
// Serves as a sentinel value that indicates
// that no "real" value is present.
var Undefined = UndefinedType{}
var UndefinedClass *Class

func initUndefined() {
	UndefinedClass = NewClassWithOptions(
		ClassWithName("Undefined"),
	)
}

func (UndefinedType) Class() *Class {
	return UndefinedClass
}

func (UndefinedType) DirectClass() *Class {
	return UndefinedClass
}

func (UndefinedType) SingletonClass() *Class {
	return nil
}

func (UndefinedType) Inspect() string {
	return "undefined"
}

func (u UndefinedType) Error() string {
	return u.Inspect()
}

func (UndefinedType) InstanceVariables() SymbolMap {
	return nil
}
