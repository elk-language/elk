package vm

import (
	"github.com/elk-language/elk/value"
)

func init() {
	// Instance methods
	c := &value.ClassClass.MethodContainer
	Accessor(c, "doc")
}
