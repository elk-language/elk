package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Weak
func initWeak() {
	// Instance methods
	c := &value.WeakClass.MethodContainer
	Def(
		c,
		"#init",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			switch box := args[1].AsReference().(type) {
			case *value.Box:
				return value.MakeWeak(box).ToValue(), value.Undefined
			case *value.ImmutableBox:
				return value.MakeWeak(box.ToBox()).ToValue(), value.Undefined
			case *LocalBox:
				return value.MakeWeak(box.ToBox()).ToValue(), value.Undefined
			default:
				return value.Undefined, value.Ref(
					value.Errorf(
						value.ArgumentErrorClass,
						"expected a box, got: %s",
						args[1].Inspect(),
					),
				)
			}
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"to_immutable_box",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsWeak()
			return self.ToImmutableBoxValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_box",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsWeak()
			return self.ToBoxValue(), value.Undefined
		},
	)
}
