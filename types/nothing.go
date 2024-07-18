package types

// Nothing represents no valid type. It is the bottom type and the top type.
// It is a subtype of all other types. And all types are subtypes of nothing.
// All method calls on nothing are valid.
//
// It is used internally in the typechecker as a return type for invalid expressions like
// calls to methods that do not exist. It helps with avoiding cascading errors.
type Nothing struct{}

func (n Nothing) ToNonLiteral(env *GlobalEnvironment) Type {
	return n
}

func (Nothing) IsLiteral() bool {
	return false
}

func IsNothing(t Type) bool {
	_, ok := t.(Nothing)
	return ok
}

func (Nothing) inspect() string {
	return "nothing"
}
