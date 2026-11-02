package macros

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Macro execution context created by the typechecker
// and passed to the macro during expansion.
type Context struct {
	TypeChecker types.Checker
	Location    *position.Location // location of the macro call
	Namespace   types.Namespace    // namespace the macro call is nested under
}

var _ value.Reference = &Context{}

func NewContext(checker types.Checker, loc *position.Location, namespace types.Namespace) *Context {
	return &Context{
		TypeChecker: checker,
		Location:    loc,
		Namespace:   namespace,
	}
}

func (*Context) Class() *value.Class {
	return value.MacroContextClass
}

func (*Context) DirectClass() *value.Class {
	return value.MacroContextClass
}

func (*Context) SingletonClass() *value.Class {
	return nil
}

func (c *Context) Error() string {
	return c.Inspect()
}

func (c *Context) Inspect() string {
	return fmt.Sprintf("Std::Macro::Context{&: %p}", c)
}

func (*Context) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (c *Context) ToValue() value.Value {
	return value.Ref(c)
}

func (c *Context) Copy() value.Reference {
	return &Context{
		TypeChecker: c.TypeChecker,
		Namespace:   c.Namespace,
		Location:    c.Location,
	}
}
