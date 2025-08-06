package types

import "github.com/elk-language/elk/value/symbol"

type False struct{}

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
