package vm

import (
	"time"

	"github.com/elk-language/elk/value"
)

// Std::Date
func initDate() {
	// Singleton methods
	c := &value.DateClass.SingletonClass().MethodContainer
	Def(
		c,
		"now",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			year, month, day := time.Now().Date()
			return value.MakeDate(year, int(month), day).ToValue(), value.Undefined
		},
	)

	// Instance methods
	c = &value.DateClass.MethodContainer
	Def(
		c,
		"#init",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			var argYear int
			if !args[1].IsUndefined() {
				if args[1].IsSmallInt() {
					argYear = int(args[1].AsSmallInt())
				} else {
					argYear = value.MaxSmallInt
				}
			}

			var argMonth int
			if !args[2].IsUndefined() {
				if args[2].IsSmallInt() {
					argMonth = int(args[2].AsSmallInt())
				} else {
					argMonth = value.MaxSmallInt
				}
			} else {
				argMonth = 1
			}

			var argDay int
			if !args[3].IsUndefined() {
				if args[3].IsSmallInt() {
					argDay = int(args[3].AsSmallInt())
				} else {
					argDay = value.MaxSmallInt
				}
			} else {
				argDay = 1
			}

			self, err := value.MakeValidatedDate(
				argYear,
				argMonth,
				argDay,
			)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return self.ToValue(), value.Undefined
		},
		DefWithParameters(3),
	)

	Def(
		c,
		"year",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsDate()
			return value.SmallInt(self.Year()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"month",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDate()
			return value.SmallInt(self.Month()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"day",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDate()
			return value.SmallInt(self.Day()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustDate()
			return value.Ref(self.ToString()), value.Undefined
		},
	)
}
