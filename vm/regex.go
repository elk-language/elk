package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::Regex
func init() {
	// Instance methods
	c := &value.RegexClass.MethodContainer
	Def(
		c,
		"match",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Regex)
			return value.ToValueErr(self.Match(args[1]))
		},
		DefWithParameters("str"),
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Regex)
			other := args[1]
			return value.ToValueErr(self.Concat(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.Regex)
			other := args[1]
			return value.ToValueErr(self.Repeat(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)

}
