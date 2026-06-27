package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::Aborter
func initAborter() {
	// Singleton methods
	c := &value.AborterClass.SingletonClass().MethodContainer
	Def(
		c,
		"closed",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			aborter := value.NewClosedAborter()
			return aborter.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"timeout",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			span := args[1].AsTimeSpan()
			var parent *value.Aborter
			if args[2].IsNotUndefined() {
				parent = args[2].AsReference().(*value.Aborter)
			} else {
				parent = value.GLOBAL_ABORTER
			}
			aborter := value.NewTimeoutAborter(parent, span)
			return aborter.ToValue(), value.Undefined
		},
		DefWithParameters(2),
	)
	Def(
		c,
		"deadline",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			datetime := args[1].AsReference().(*value.DateTime)
			var parent *value.Aborter
			if args[2].IsNotUndefined() {
				parent = args[2].AsReference().(*value.Aborter)
			} else {
				parent = value.GLOBAL_ABORTER
			}
			aborter := value.NewDeadlineAborter(parent, datetime)
			return aborter.ToValue(), value.Undefined
		},
		DefWithParameters(2),
	)

	// Instance methods
	c = &value.AborterClass.MethodContainer
	Def(
		c,
		"#init",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			var parent *value.Aborter
			if args[1].IsNotUndefined() {
				parent = args[1].AsReference().(*value.Aborter)
			} else {
				parent = value.GLOBAL_ABORTER
			}

			aborter := value.NewCancelAborter(parent)
			return aborter.ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"is_closed",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Aborter)(args[0].Pointer())
			return value.BoolVal(self.IsClosed()), value.Undefined
		},
	)
	Def(
		c,
		"is_closable",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Aborter)(args[0].Pointer())
			return value.BoolVal(self.IsClosable()), value.Undefined
		},
	)
	Def(
		c,
		"closed",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Aborter)(args[0].Pointer())
			return self.Closed().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"close",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Aborter)(args[0].Pointer())
			self.Close()
			return value.Nil, value.Undefined
		},
	)

}
