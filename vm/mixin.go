package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Mixin
func initMixin() {
	// Instance methods
	c := &value.MixinClass.MethodContainer
	Accessor(c, "doc")
}
