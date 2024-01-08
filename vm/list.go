package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::List
func init() {
	// Instance methods
	c := &value.ListClass.MethodContainer
	Def(
		c,
		"iterator",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.List)
			iterator := value.NewListIterator(self)
			return iterator, nil
		},
	)

}

// ::Std::ListIterator
func init() {
	// Instance methods
	c := &value.ListIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ListIterator)
			return self.Next()
		},
	)

}
