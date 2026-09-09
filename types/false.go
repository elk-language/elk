package types

import "github.com/elk-language/elk/value/symbol"

type False struct{}

func (False) HashUint64() uint64 {
	return 3
}

func (False) EqualAny(other any) bool {
	_, ok := other.(False)
	return ok
}

func (False) ID() ID {
	return 3
}

func (False) SetID(ID) {}

func (f False) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(f, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(f, parent)
	}
}

func (False) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_Bool)
}

func (False) IsLiteral() bool {
	return true
}

func IsFalse(t Type) bool {
	return IsFalseLiteral(t) || t == Env.StdSubtype(symbol.C_False)
}

func IsFalseLiteral(t Type) bool {
	_, ok := t.(False)
	return ok
}

func (False) inspect() string {
	return "false"
}
