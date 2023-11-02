package value

type UndefinedType struct{}

// Elk's internal undefined value
// Serves as a sentinel value that indicates
// that no "real" value is present.
var Undefined = UndefinedType{}

func (UndefinedType) Class() *Class {
	return nil
}

func (UndefinedType) IsFrozen() bool {
	return true
}

func (UndefinedType) SetFrozen() {}

func (UndefinedType) Inspect() string {
	return "undefined"
}

func (UndefinedType) InstanceVariables() SimpleSymbolMap {
	return nil
}
