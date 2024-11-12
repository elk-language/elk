package vm

import (
	"github.com/elk-language/elk/value"
)

func init() {
	// Instance methods
	c := &value.MixinClass.MethodContainer
	Accessor(c, "doc")
}
