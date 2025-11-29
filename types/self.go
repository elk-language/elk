package types

type Self struct{}

func (Self) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(Self{}, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(Self{}, parent)
	}
}

func (s Self) ToNonLiteral(env *GlobalEnvironment) Type {
	return s
}

func (Self) IsLiteral() bool {
	return false
}

func IsSelf(t Type) bool {
	_, ok := t.(Self)
	return ok
}

func (Self) inspect() string {
	return "self"
}
