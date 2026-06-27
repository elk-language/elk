package vm

import (
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Wraps an anonymous with associated local variables
// from the outer context
type Closure interface {
	value.ValueInterface
	HasOpenUpvalues() bool
	ParameterCount() int
	Location() *position.Location
}
