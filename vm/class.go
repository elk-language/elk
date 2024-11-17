package vm

import (
	"github.com/elk-language/elk/value"
)

func initClass() {
	// Instance methods
	c := &value.ClassClass.MethodContainer
	Accessor(c, "doc")
}
