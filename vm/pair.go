package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::Pair
func init() {
	// Instance methods
	c := &value.PairClass.MethodContainer
	Def(
		c,
		"#init",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Pair)
			self.Key = args[1]
			self.Value = args[2]
			return self, nil
		},
		DefWithParameters("key", "value"),
	)
	Def(
		c,
		"key",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Pair)
			return self.Key, nil
		},
	)
	Def(
		c,
		"value",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Pair)
			return self.Value, nil
		},
	)

}
