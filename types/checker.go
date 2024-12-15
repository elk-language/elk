package types

import (
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

type Checker interface {
	Env() *GlobalEnvironment
	IsSubtype(a, b Type, errSpan *position.Span) bool
	IsNilable(typ Type) bool
	IsNotNilable(typ Type) bool
	IsNil(typ Type) bool
	IsFalsy(typ Type) bool
	IsTruthy(typ Type) bool
	CanBeFalsy(typ Type) bool
	CanBeTruthy(typ Type) bool
	IsTheSameType(a, b Type, errSpan *position.Span) bool
	TypesIntersect(a, b Type) bool
	NewNormalisedUnion(elements ...Type) Type
	NewNormalisedIntersection(elements ...Type) Type
	NormaliseType(typ Type) Type
	StdInt() *Class
	StdFloat() *Class
	StdBigFloat() *Class
	StdClass() *Class
	Std(name value.Symbol) Type
	StdString() Type
	StdBool() *Class
	StdNil() *Class
	StdTrue() *Class
	StdFalse() *Class
}
