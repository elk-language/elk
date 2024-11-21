package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::Regex
func initRegex() {
	// Instance methods
	c := &value.RegexClass.MethodContainer
	Def(
		c,
		"matches",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Regex)
			return self.Matches(args[1])
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Regex)
			other := args[1]
			return self.Concat(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Regex)
			other := args[1]
			return self.Repeat(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Regex)
			withFlags := args[1]
			if withFlags != value.Undefined && value.Truthy(withFlags) {
				return self.ToStringWithFlags(), nil
			}
			return self.ToString(), nil
		},
		DefWithParameters(1),
		DefWithOptionalParameters(1),
	)

}
