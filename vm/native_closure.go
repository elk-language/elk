package vm

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Wraps a native function with associated local variables
// from the outer context
type NativeClosure struct {
	Function       NativeFunction
	loc            *position.Location
	parameterCount int
}

var _ Closure = &NativeClosure{}

func NewNativeClosure(fn NativeFunction, parameterCount int, loc *position.Location) *NativeClosure {
	return &NativeClosure{
		Function:       fn,
		parameterCount: parameterCount,
		loc:            loc,
	}
}

func (c *NativeClosure) ParameterCount() int {
	return c.parameterCount
}

func (c *NativeClosure) Location() *position.Location {
	return c.loc
}

func (c *NativeClosure) HasOpenUpvalues() bool {
	return false
}

func (*NativeClosure) Class() *value.Class {
	return value.ClosureClass
}

func (*NativeClosure) DirectClass() *value.Class {
	return value.ClosureClass
}

func (*NativeClosure) SingletonClass() *value.Class {
	return nil
}

func (c *NativeClosure) Copy() value.Reference {
	return c
}

func (c *NativeClosure) ToValue() value.Value {
	return value.Ref(c)
}

func (c *NativeClosure) Inspect() string {
	return fmt.Sprintf("Std::Closure{location: %s, type: :native}", c.loc.String())
}

func (c *NativeClosure) Error() string {
	return c.Inspect()
}

func (*NativeClosure) InstanceVariables() *value.InstanceVariables {
	return nil
}
