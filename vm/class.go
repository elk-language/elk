package vm

import (
	"github.com/elk-language/elk/value"
)

func init() {
	// DefineMethodWithOptions(
	// 	value.FloatClass.Methods,
	// 	"<=>",
	// 	func(_ *VM, args []value.Value) (value.Value, value.Value) {
	// 		self := args[0].(value.Float)
	// 		other := args[1]
	// 		return value.ToValueErr(self.Compare(other))
	// 	},
	// 	NativeMethodWithStringParameters("other"),
	// 	NativeMethodWithFrozen(),
	// )
	// DefineMethodWithOptions(
	// 	value.FloatClass.Methods,
	// 	">",
	// 	func(_ *VM, args []value.Value) (value.Value, value.Value) {
	// 		self := args[0].(value.Float)
	// 		other := args[1]
	// 		return value.ToValueErr(self.GreaterThan(other))
	// 	},
	// 	NativeMethodWithStringParameters("other"),
	// 	NativeMethodWithFrozen(),
	// )
	// DefineMethodWithOptions(
	// 	value.FloatClass.Methods,
	// 	">=",
	// 	func(_ *VM, args []value.Value) (value.Value, value.Value) {
	// 		self := args[0].(value.Float)
	// 		other := args[1]
	// 		return value.ToValueErr(self.GreaterThanEqual(other))
	// 	},
	// 	NativeMethodWithStringParameters("other"),
	// 	NativeMethodWithFrozen(),
	// )
	// DefineMethodWithOptions(
	// 	value.FloatClass.Methods,
	// 	"<",
	// 	func(_ *VM, args []value.Value) (value.Value, value.Value) {
	// 		self := args[0].(value.Float)
	// 		other := args[1]
	// 		return value.ToValueErr(self.LessThan(other))
	// 	},
	// 	NativeMethodWithStringParameters("other"),
	// 	NativeMethodWithFrozen(),
	// )
	// DefineMethodWithOptions(
	// 	value.FloatClass.Methods,
	// 	"<=",
	// 	func(_ *VM, args []value.Value) (value.Value, value.Value) {
	// 		self := args[0].(value.Float)
	// 		other := args[1]
	// 		return value.ToValueErr(self.LessThanEqual(other))
	// 	},
	// 	NativeMethodWithStringParameters("other"),
	// 	NativeMethodWithFrozen(),
	// )

	Accessor(&value.ClassClass.MethodContainer, "doc", false)

	Def(
		&value.ClassClass.MethodContainer,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0]
			return value.String(self.Inspect()), nil
		},
	)

}
