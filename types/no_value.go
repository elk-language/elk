package types

const NoValueID = 7

// Represents no value.
// It is used to mark constants of pure types, that do not have a runtime value.
type NoValue struct{}

func (NoValue) HashUint64() uint64 {
	return NoValueID
}

func (NoValue) EqualAny(other any) bool {
	_, ok := other.(NoValue)
	return ok
}

func (n NoValue) CopyType() Type {
	return n
}

func (NoValue) ID() ID {
	return NoValueID
}

func (NoValue) SetID(ID) {}

func (NoValue) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(NoValue{}, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(NoValue{}, parent)
	}
}

func (n NoValue) ToNonLiteral() Type {
	return n
}

func (NoValue) IsLiteral() bool {
	return false
}

func IsNoValue(t Type) bool {
	_, ok := t.(NoValue)
	return ok
}

func IsNoValueID(id ID) bool {
	return id == NoValueID
}

func IsNoValueRef[T Type](ref Ref[T]) bool {
	return ref == NoValueID
}

func (NoValue) inspect() string {
	return "NoValue"
}
