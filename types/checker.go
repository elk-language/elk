package types

import (
	"github.com/elk-language/elk/value"
)

type Checker interface {
	Env() *GlobalEnvironment
	IsSubtype(a, b Type) bool
	IsNilable(typ Type) bool
	IsNotNilable(typ Type) bool
	IsNil(typ Type) bool
	IsFalsy(typ Type) bool
	IsTruthy(typ Type) bool
	CanBeFalsy(typ Type) bool
	CanBeTruthy(typ Type) bool
	IsTheSameType(a, b Type) bool
	TypesIntersect(a, b Type) bool
	NewNormalisedUnion(elements ...Type) Type
	NewNormalisedIntersection(elements ...Type) Type
	NormaliseType(typ Type) Type
	StdInt() *Class
	StdFloat() *Class
	StdBigFloat() *Class
	StdClass() *Class
	Std(name value.Symbol) Type
	StdString() *Class
	StdBool() *Class
	StdNil() *Class
	StdTrue() *Class
	StdFalse() *Class
	GetIteratorElementType(Type) (Type, Type)
	GetIteratorType(Type) Type
}
