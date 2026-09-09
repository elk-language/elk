package types

type Self struct{}

func (Self) HashUint64() uint64 {
	return 8
}

func (Self) EqualAny(other any) bool {
	_, ok := other.(Self)
	return ok
}

func (Self) ID() ID {
	return 8
}

func (Self) SetID(ID) {}

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
