package types

import "github.com/elk-language/elk/value"

type ConstantContainer interface {
	Type
	Name() string
	Constants() map[value.Symbol]Type
	Constant(name string) Type
	DefineConstant(name string, val Type)
	Subtypes() map[value.Symbol]Type
	Subtype(name string) Type
	DefineSubtype(name string, val Type)
}
