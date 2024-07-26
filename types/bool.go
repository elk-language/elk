package types

import "github.com/elk-language/elk/value/symbol"

type Bool struct{}

func (Bool) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Bool)
}

func (Bool) IsLiteral() bool {
	return true
}

func IsBool(t Type) bool {
	_, ok := t.(Bool)
	return ok
}

func (Bool) inspect() string {
	return "bool"
}
