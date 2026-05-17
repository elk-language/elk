package vm

import (
	"reflect"
	"strings"

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

func (s *Select) Inspect() string {
	var buff strings.Builder
	buff.WriteString("Select{")
	for i, selectCase := range s.Cases {
		if i != 0 {
			buff.WriteString(", ")
		}

		switch selectCase.Direction {
		case reflect.SelectSend:
			buff.WriteString("send")
		case reflect.SelectRecv:
			buff.WriteString("receive")
		case reflect.SelectDefault:
			buff.WriteString("else")
		default:
			buff.WriteString("unknown")
		}
	}
	buff.WriteString("}")

	return buff.String()
}

func (s *Select) String() string {
	return s.Inspect()
}

func (s *Select) Error() string {
	return s.Inspect()
}

type SelectCase struct {
	Direction reflect.SelectDir
}
