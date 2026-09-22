package types

import "github.com/elk-language/elk/value/symbol"

const TrueID = 4

type True struct{}

func (True) HashUint64() uint64 {
	return TrueID
}

func (True) EqualAny(other any) bool {
	_, ok := other.(True)
	return ok
}

func (t True) CopyType() Type {
	return t
}

func (True) ID() ID {
	return TrueID
}

func (True) SetID(ID) {}

func (True) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(True{}, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(True{}, parent)
	}
}

func (True) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_Bool)
}

func (True) IsLiteral() bool {
	return true
}

func IsTrue(t Type) bool {
	return IsTrueLiteral(t) || t == Env.StdSubtype(symbol.C_True)
}

func IsTrueLiteral(t Type) bool {
	_, ok := t.(True)
	return ok
}

func (True) inspect() string {
	return "true"
}
