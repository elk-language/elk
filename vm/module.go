package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Module
func initModule() {
	// Instance methods
	c := &value.ModuleClass.MethodContainer
	Accessor(c, "doc")
}
