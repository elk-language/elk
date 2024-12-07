package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Timezone
func initTimezone() {
	// Class methods
	c := &value.TimezoneClass.SingletonClass().MethodContainer
	Def(
		c,
		"get",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			name := args[1].MustReference().(value.String)
			return value.RefErr(value.LoadTimezone(string(name)))
		},
		DefWithParameters(1),
	)
	Alias(c, "[]", "get")

	// Instance methods
	c = &value.TimezoneClass.MethodContainer
	Def(
		c,
		"name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Timezone)
			return value.Ref(value.String(self.Name())), value.Nil
		},
	)
	Alias(c, "to_string", "name")

	Def(
		c,
		"is_utc",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Timezone)
			return value.ToElkBool(self.IsUTC()), value.Nil
		},
	)
	Def(
		c,
		"is_local",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Timezone)
			return value.ToElkBool(self.IsLocal()), value.Nil
		},
	)

}
