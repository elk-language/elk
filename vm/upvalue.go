package vm

import (
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value"
)

const (
	UpvalueLongIndexFlag bitfield.BitFlag8 = 1 << iota
	UpvalueLocalFlag
)

// Represents a captured variable from an outer context
type Upvalue struct {
	location *value.Value
}

func NewUpvalue(loc *value.Value) *Upvalue {
	return &Upvalue{
		location: loc,
	}
}
