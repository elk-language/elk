package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::StackTrace
func initStackTrace() {
	// Instance methods
	c := &value.StackTraceClass.MethodContainer
	Def(
		c,
		"[]",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.StackTrace)(args[0].Pointer())
			nVal := args[1]
			n, ok := value.ToGoInt(nVal)
			if !ok {
				if n == -1 {
					return value.Undefined, value.Ref(value.NewIndexOutOfRangeError(nVal.Inspect(), self.Length()))
				}
				return value.Undefined, value.Ref(value.NewCoerceError(value.IntClass, nVal.Class()))
			}
			cf, err := self.Get(n)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.Ref(cf), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"length",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.StackTrace)(args[0].Pointer())
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.StackTrace)(args[0].Pointer())
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.StackTrace)(args[0].Pointer())
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.StackTrace)(args[0].Pointer())
			return value.Ref(value.NewStackTraceIterator(self)), value.Undefined
		},
	)

}

// ::Std::StackTrace::Iterator
func initStackTraceIterator() {
	// Instance methods
	c := &value.StackTraceIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.StackTraceIterator)(args[0].Pointer())
			return self.Next()
		},
	)
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"reset",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.StackTraceIterator)(args[0].Pointer())
			self.Reset()
			return args[0], value.Undefined
		},
	)
}
