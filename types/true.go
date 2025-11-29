package types

import "github.com/elk-language/elk/value/symbol"

type True struct{}

func (True) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(True{}, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(True{}, parent)
	}
}

func (True) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Bool)
}

func (True) IsLiteral() bool {
	return true
}

func IsTrue(t Type, env *GlobalEnvironment) bool {
	return IsTrueLiteral(t) || t == env.StdSubtype(symbol.True)
}

func IsTrueLiteral(t Type) bool {
	_, ok := t.(True)
	return ok
}

func (True) inspect() string {
	return "true"
}
