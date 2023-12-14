package vm

import (
	"github.com/elk-language/elk/value"
)

func init() {
	// Instance methods
	c := &value.ModuleClass.MethodContainer
	Accessor(c, "doc", false)
}
