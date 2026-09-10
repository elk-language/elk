package types

const NeverID = 6

// Never represents no value. It is the bottom type.
// For example a function that never returns
// might use the type `never`.
// It is a subtype of all other types.
// All method calls on never are valid.
//
// It is used to detect unreachable code.
type Never struct{}

func (Never) HashUint64() uint64 {
	return NeverID
}

func (Never) EqualAny(other any) bool {
	_, ok := other.(Never)
	return ok
}

func (Never) ID() ID {
	return NeverID
}

func (Never) SetID(ID) {}

func (Never) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(Never{}, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(Never{}, parent)
	}
}

func (n Never) ToNonLiteral() Type {
	return n
}

func (Never) IsLiteral() bool {
	return false
}

func IsNever(t Type) bool {
	_, ok := t.(Never)
	return ok
}

func IsNeverID(id ID) bool {
	return id == NeverID
}

func IsNeverRef[T Type](ref Ref[T]) bool {
	return ref == NeverID
}

func (Never) inspect() string {
	return "never"
}
