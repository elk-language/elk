package types

import "github.com/elk-language/elk/value/symbol"

type Nil struct{}

func (Nil) HashUint64() uint64 {
	return 5
}

func (Nil) EqualAny(other any) bool {
	_, ok := other.(Nil)
	return ok
}

func (Nil) ID() ID {
	return 5
}

func (Nil) SetID(ID) {}

func (Nil) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(Nil{}, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(Nil{}, parent)
	}
}

func (v Nil) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.C_Nil)
}

func (Nil) IsLiteral() bool {
	return true
}

func IsNilLiteral(t Type) bool {
	_, ok := t.(Nil)
	return ok
}

func IsNil(t Type, env *GlobalEnvironment) bool {
	return IsNilLiteral(t) || t == env.StdSubtype(symbol.C_Nil)
}

func (Nil) inspect() string {
	return "nil"
}
