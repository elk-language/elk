package types

// Void is the type that is incompatible with
// any other type.
type Void struct{}

func (Void) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(Void{}, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(Void{}, parent)
	}
}

func (v Void) ToNonLiteral(env *GlobalEnvironment) Type {
	return v
}

func (Void) IsLiteral() bool {
	return false
}

func IsVoid(t Type) bool {
	_, ok := t.(Void)
	return ok
}

func (Void) inspect() string {
	return "void"
}
