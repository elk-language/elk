package types

const AnyID = 1

// All types are subtypes of any.
// Any is not a subtype of anything other than itself.
// It is the top type.
type Any struct{}

func (a Any) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(a, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(a, parent)
	}
}

func (n Any) ToNonLiteral() Type {
	return n
}

func (Any) IsLiteral() bool {
	return false
}

func (Any) HashUint64() uint64 {
	return AnyID
}

func (Any) EqualAny(other any) bool {
	_, ok := other.(Any)
	return ok
}

func (Any) ID() ID {
	return AnyID
}

func (Any) SetID(ID) {}

func IsAny(t Type) bool {
	_, ok := t.(Any)
	return ok
}

func IsAnyID(id ID) bool {
	return id == AnyID
}

func IsAnyRef[T Type](ref Ref[T]) bool {
	return ref == AnyID
}

func (Any) inspect() string {
	return "any"
}
