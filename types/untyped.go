package types

// Untyped represents no valid type. It is the bottom type and the top type.
// It is a subtype of all other types. And all types are subtypes of nothing.
// All method calls on nothing are valid.
//
// It is used internally in the typechecker as a return type for invalid expressions like
// calls to methods that do not exist. It helps with avoiding cascading errors.
type Untyped struct{}

func (Untyped) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(Untyped{}, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(Untyped{}, parent)
	}
}

func (n Untyped) ToNonLiteral(env *GlobalEnvironment) Type {
	return n
}

func (Untyped) IsLiteral() bool {
	return false
}

func IsUntyped(t Type) bool {
	_, ok := t.(Untyped)
	return ok
}

func (Untyped) inspect() string {
	return "untyped"
}
