package types

import "github.com/elk-language/elk/value/symbol"

type False struct{}

func (False) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.False)
}

func (False) IsLiteral() bool {
	return true
}

func IsFalse(t Type) bool {
	_, ok := t.(False)
	return ok
}

func (False) inspect() string {
	return "false"
}
