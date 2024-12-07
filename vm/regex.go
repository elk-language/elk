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
			self := args[0].MustReference().(*value.Regex)
			return self.Matches(args[1])
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Regex)
			other := args[1]
			return self.Concat(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Regex)
			other := args[1]
			return self.Repeat(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Regex)
			withFlags := args[1]
			if !withFlags.IsUndefined() && value.Truthy(withFlags) {
				return value.Ref(self.ToStringWithFlags()), value.Undefined
			}
			return value.Ref(self.ToString()), value.Undefined
		},
		DefWithParameters(1),
		DefWithOptionalParameters(1),
	)

}
