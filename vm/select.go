package vm

import (
	"reflect"

	"github.com/elk-language/elk/value"
)

// Wraps data for `select` expressions
type Select struct {
	Cases []SelectCase
	value.ValueBase
}

var _ value.Reference = &Select{}

func NewSelect(cases []SelectCase) *Select {
	return &Select{
		Cases: cases,
	}
}

func (s *Select) ToValue() value.Value {
	return value.Ref(s)
}

func (s *Select) Copy() value.Reference {
	return &Select{
		Cases: s.Cases,
	}
}

type SelectCase struct {
	Direction reflect.SelectDir
}
