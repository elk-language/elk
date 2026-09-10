package types

const VoidID = 10

// Void is the type that is incompatible with
// any other type.
type Void struct{}

func (Void) HashUint64() uint64 {
	return VoidID
}

func (Void) EqualAny(other any) bool {
	_, ok := other.(Void)
	return ok
}

func (Void) ID() ID {
	return VoidID
}

func (Void) SetID(ID) {}

func (Void) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(Void{}, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(Void{}, parent)
	}
}

func (v Void) ToNonLiteral() Type {
	return v
}

func (Void) IsLiteral() bool {
	return false
}

func IsVoid(t Type) bool {
	_, ok := t.(Void)
	return ok
}

func IsVoidID(id ID) bool {
	return id == VoidID
}

func IsVoidRef[T Type](ref Ref[T]) bool {
	return ref == VoidID
}

func (Void) inspect() string {
	return "void"
}
