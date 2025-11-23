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
	Def(
		c,
		"from_offset",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			offset := args[1].AsTimeSpan()
			return value.Ref(value.NewTimezoneFromOffset(offset)), value.Undefined
		},
		DefWithParameters(1),
	)

	// Instance methods
	c = &value.TimezoneClass.MethodContainer
	Def(
		c,
		"name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Timezone)
			return value.Ref(value.String(self.Name())), value.Undefined
		},
	)
	Alias(c, "to_string", "name")

	Def(
		c,
		"offset",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Timezone)
			return self.StandardOffset().ToValue(), value.Undefined
		},
	)
	Alias(c, "standard_offset", "offset")

	Def(
		c,
		"dst_offset",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Timezone)
			return self.DSTOffset().ToValue(), value.Undefined
		},
	)

	Def(
		c,
		"is_utc",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Timezone)
			return value.ToElkBool(self.IsUTC()), value.Undefined
		},
	)
	Def(
		c,
		"is_local",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Timezone)
			return value.ToElkBool(self.IsLocal()), value.Undefined
		},
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.Timezone)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		DefWithParameters(1),
	)

}
