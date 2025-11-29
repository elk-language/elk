package types

import "github.com/elk-language/elk/value/symbol"

type Nil struct{}

func (Nil) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(Nil{}, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(Nil{}, parent)
	}
}

func (v Nil) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Nil)
}

func (Nil) IsLiteral() bool {
	return true
}

func IsNilLiteral(t Type) bool {
	_, ok := t.(Nil)
	return ok
}

func IsNil(t Type, env *GlobalEnvironment) bool {
	return IsNilLiteral(t) || t == env.StdSubtype(symbol.Nil)
}

func (Nil) inspect() string {
	return "nil"
}
