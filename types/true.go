package types

import "github.com/elk-language/elk/value/symbol"

type True struct{}

func (True) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Bool)
}

func (True) IsLiteral() bool {
	return true
}

func IsTrue(t Type) bool {
	_, ok := t.(True)
	return ok
}

func (True) inspect() string {
	return "true"
}
