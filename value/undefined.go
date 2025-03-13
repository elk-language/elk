package value

// Elk's internal undefined value.
// Serves as a sentinel value that indicates
// that no "real" value is present.
//
// It is the zero value of `Value` and maps to Go's `nil`
type UndefinedType struct{}

var Undefined = UndefinedType{}.ToValue()
var UndefinedClass *Class

func initUndefined() {
	UndefinedClass = NewClassWithOptions(
		ClassWithName("Undefined"),
		ClassWithParent(ValueClass),
	)
}

func (u UndefinedType) ToValue() Value {
	return Value{
		flag: UNDEFINED_FLAG,
	}
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
