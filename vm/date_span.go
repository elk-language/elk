package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Date::Span
func initDateSpan() {
	// Instance methods
	c := &value.DateSpanClass.MethodContainer
	Def(
		c,
		"#init",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			var argYear int
			if !args[1].IsUndefined() {
				if args[1].IsSmallInt() {
					argYear = int(args[1].AsSmallInt())
				} else {
					argYear = int(args[1].AsBigInt().ToSmallInt())
				}
			}

			var argMonth int
			if !args[2].IsUndefined() {
				if args[2].IsSmallInt() {
					argMonth = int(args[2].AsSmallInt())
				} else {
					argMonth = int(args[2].AsBigInt().ToSmallInt())
				}
			}

			var argDay int
			if !args[3].IsUndefined() {
				if args[3].IsSmallInt() {
					argDay = int(args[3].AsSmallInt())
				} else {
					argDay = int(args[3].AsBigInt().ToSmallInt())
				}
			}

			self := value.MakeDateSpan(
				argYear,
				argMonth,
				argDay,
			)
			return self.ToValue(), value.Undefined
		},
		DefWithParameters(3),
	)

	Def(
		c,
		"years",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return value.SmallInt(self.Years()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"months",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return value.SmallInt(self.Months()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"months_mod",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return value.SmallInt(self.MonthsMod()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"days",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return value.SmallInt(self.Days()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDateSpan()
			return value.Ref(self.ToString()), value.Undefined
		},
	)
}
