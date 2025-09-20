package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Time
func initTime() {
	// Class methods
	c := &value.TimeClass.SingletonClass().MethodContainer
	Def(
		c,
		"now",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return value.Ref(value.TimeNow()), value.Undefined
		},
	)

	// Instance methods
	c = &value.TimeClass.MethodContainer
	Def(
		c,
		"format",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			arg := args[1]
			switch a := arg.SafeAsReference().(type) {
			case value.String:
				result, err := self.Format(string(a))
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.Ref(value.String(result)), value.Undefined
			default:
				return value.Undefined, value.Ref(value.Errorf(
					value.ArgumentErrorClass,
					"expected a format string, got: %s",
					arg.Inspect(),
				))
			}
		},
		DefWithParameters(1),
	)
	Alias(c, "strftime", "format")

	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			return value.Ref(self.ToString()), value.Undefined
		},
	)
	Def(
		c,
		"to_time_span",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			return self.ToTimeSpan().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"hour",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			return value.SmallInt(self.Hour()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"minute",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			return value.SmallInt(self.Minute()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			return value.SmallInt(self.Second()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"millisecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			return value.SmallInt(self.Millisecond()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"microsecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			return value.SmallInt(self.Microsecond()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"microseconds_in_second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			return value.SmallInt(self.MicrosecondsInSecond()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"nanosecond",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			return value.SmallInt(self.Nanosecond()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"nanoseconds_in_second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			return value.SmallInt(self.NanosecondsInSecond()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"picoseconds_in_second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			return value.ToElkInt(self.PicosecondsInSecond()), value.Undefined
		},
	)
	Def(
		c,
		"femtoseconds_in_second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			return value.ToElkInt(self.FemtosecondsInSecond()), value.Undefined
		},
	)
	Def(
		c,
		"attoseconds_in_second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			return value.ToElkInt(self.AttosecondsInSecond()), value.Undefined
		},
	)
	Def(
		c,
		"zeptoseconds_in_second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			return value.Ref(self.ZeptosecondsInSecond()), value.Undefined
		},
	)
	Def(
		c,
		"yoctoseconds_in_second",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			return value.Ref(self.YoctosecondsInSecond()), value.Undefined
		},
	)

	Def(
		c,
		"is_am",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			return value.ToElkBool(self.IsAM()), value.Undefined
		},
	)
	Def(
		c,
		"is_pm",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			return value.ToElkBool(self.IsPM()), value.Undefined
		},
	)
	Def(
		c,
		"meridiem",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			return value.Ref(value.String(self.Meridiem())), value.Undefined
		},
	)

	Def(
		c,
		"hour12",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsTime()
			return value.SmallInt(self.Hour12()).ToValue(), value.Undefined
		},
	)
}
