package types

import "github.com/elk-language/elk/value/symbol"

type Nil struct{}

func (v Nil) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Nil)
}

func (Nil) IsLiteral() bool {
	return true
}

func IsNil(t Type) bool {
	_, ok := t.(Nil)
	return ok
}

func (Nil) inspect() string {
	return "nil"
}
