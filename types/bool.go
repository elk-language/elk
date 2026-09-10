package types

import "github.com/elk-language/elk/value/symbol"

const BoolID = 2

type Bool struct{}

func (Bool) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_Bool)
}

func (b Bool) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(b, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(b, parent)
	}
}

func (Bool) IsLiteral() bool {
	return true
}

func IsBool(t Type) bool {
	_, ok := t.(Bool)
	return ok
}

func IsBoolID(id ID) bool {
	return id == BoolID
}

func IsBoolRef[T Type](ref Ref[T]) bool {
	return ref == BoolID
}

func (Bool) HashUint64() uint64 {
	return BoolID
}

func (Bool) EqualAny(other any) bool {
	_, ok := other.(Bool)
	return ok
}

func (Bool) ID() ID {
	return BoolID
}

func (Bool) SetID(ID) {}

func (Bool) inspect() string {
	return "bool"
}
