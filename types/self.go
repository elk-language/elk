package types

const SelfID = 8

type Self struct{}

func (Self) HashUint64() uint64 {
	return SelfID
}

func (Self) EqualAny(other any) bool {
	_, ok := other.(Self)
	return ok
}

func (Self) ID() ID {
	return SelfID
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

func (s Self) ToNonLiteral() Type {
	return s
}

func (Self) IsLiteral() bool {
	return false
}

func IsSelf(t Type) bool {
	_, ok := t.(Self)
	return ok
}

func IsSelfID(id ID) bool {
	return id == SelfID
}

func IsSelfRef[T Type](ref Ref[T]) bool {
	return ref == SelfID
}

func (Self) inspect() string {
	return "self"
}
