package types

import "github.com/elk-language/elk/value"

type ConstantContainer interface {
	Constants() map[value.Symbol]Type
	Constant(name string) Type
	DefineConstant(name string, val Type)
}
