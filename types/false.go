package types

import "github.com/elk-language/elk/value/symbol"

type False struct{}

func (f False) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(f, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(f, parent)
	}
}

func (False) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Bool)
}

func (False) IsLiteral() bool {
	return true
}

func IsFalse(t Type, env *GlobalEnvironment) bool {
	return IsFalseLiteral(t) || t == env.StdSubtype(symbol.False)
}

func IsFalseLiteral(t Type) bool {
	_, ok := t.(False)
	return ok
}

func (False) inspect() string {
	return "false"
}
