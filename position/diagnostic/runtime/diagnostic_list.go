package runtime

import (
	"slices"

	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/lexer/colorizer"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

// Std::DiagnosticList
func initDiagnosticList() {
	c := &value.DiagnosticListClass.MethodContainer

	vm.Def(
		c,
		"iter",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DiagnosticList)
			iterator := value.NewDiagnosticListIterator(self)
			return value.Ref(iterator), value.Undefined
		},
	)
	vm.Def(
		c,
		"capacity",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DiagnosticList)
			return value.SmallInt(self.Capacity()).ToValue(), value.Undefined
		},
	)
	vm.Def(
		c,
		"length",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DiagnosticList)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	vm.Def(
		c,
		"left_capacity",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DiagnosticList)
			return value.SmallInt(self.LeftCapacity()).ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"[]",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DiagnosticList)(args[0].Pointer())
			other := args[1]
			return self.Subscript(other)
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"[]=",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DiagnosticList)
			key := args[1]
			val := (*value.Diagnostic)(args[2].Pointer())
			err := self.SubscriptSet(key, val)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.Ref(val), value.Undefined
		},
		vm.DefWithParameters(2),
	)
	vm.Def(
		c,
		"append",
		func(vm *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DiagnosticList)
			values := args[1].MustReference().(*value.ArrayTuple)
			*self = slices.Grow(*self, values.Length())
			for _, element := range *values {
				*self = append(*self, (*diagnostic.Diagnostic)(element.Pointer()))
			}
			return value.Ref(self), value.Undefined
		},
		vm.DefWithParameters(1),
	)
	vm.Def(
		c,
		"<<",
		func(vm *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DiagnosticList)
			val := (*value.Diagnostic)(args[1].Pointer())
			self.Append(val)
			return value.Ref(self), value.Undefined
		},
		vm.DefWithParameters(1),
	)
	vm.Alias(c, "push", "<<")

	vm.Def(
		c,
		"==",
		func(v *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DiagnosticList)(args[0].Pointer())
			switch other := args[1].SafeAsReference().(type) {
			case *value.DiagnosticList:
				equal, err := DiagnosticListEqual(v, self, other)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.ToElkBool(equal), value.Undefined
			default:
				return value.False, value.Undefined
			}
		},
		vm.DefWithParameters(1),
	)
	vm.Alias(c, "=~", "==")

	vm.Def(
		c,
		"contains",
		func(v *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DiagnosticList)(args[0].Pointer())
			contains, err := DiagnosticListContains(v, self, args[1])
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.ToElkBool(contains), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"+",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DiagnosticList)
			other := args[1]
			return value.RefErr(self.Concat(other))
		},
		vm.DefWithParameters(1),
	)
	vm.Def(
		c,
		"*",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.DiagnosticList)
			other := args[1]
			return value.RefErr(self.Repeat(other))
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_human_string",
		func(v *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := (*diagnostic.DiagnosticList)(args[0].Pointer())

			style := true
			if !args[1].IsUndefined() {
				style = value.Truthy(args[1])
			}

			var colorizer colorizer.Colorizer
			if args[2].IsUndefined() {
				colorizer = lexer.Colorizer{}
			} else if !args[2].IsNil() {
				colorizer = vm.MakeColorizer(v, args[2])
			}

			var result string
			var err error
			if args[3].IsUndefined() {
				result, err = self.HumanString(style, colorizer)
			} else {
				result, err = self.HumanStringWithSource(
					string(args[3].AsReference().(value.String)),
					style,
					colorizer,
				)
			}

			if err != nil {
				return value.Undefined, value.Ref(value.NewError(value.ColorizerErrorClass, err.Error()))
			}
			return value.Ref(value.String(result)), value.Undefined
		},
		vm.DefWithParameters(3),
	)
	vm.Def(
		c,
		"is_failure",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := (*diagnostic.DiagnosticList)(args[0].Pointer())
			return value.ToElkBool(self.IsFailure()), value.Undefined
		},
	)
}

// ::Std::DiagnosticList::Iterator
func initDiagnosticListIterator() {
	// Instance methods
	c := &value.DiagnosticListIteratorClass.MethodContainer
	vm.Def(
		c,
		"next",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DiagnosticListIterator)(args[0].Pointer())
			return self.Next()
		},
	)
	vm.Def(
		c,
		"iter",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	vm.Def(
		c,
		"reset",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.DiagnosticListIterator)(args[0].Pointer())
			self.Reset()
			return args[0], value.Undefined
		},
	)
}

func DiagnosticListEqual(v *vm.VM, x, y *value.DiagnosticList) (bool, value.Value) {
	xLen := x.Length()
	if xLen != y.Length() {
		return false, value.Undefined
	}

	for i := 0; i < xLen; i++ {
		equal, err := v.CallMethodByName(
			symbol.OpEqual,
			value.Ref(x.At(i)),
			value.Ref(y.At(i)),
		)
		if !err.IsUndefined() {
			return false, err
		}
		if value.Falsy(equal) {
			return false, value.Undefined
		}
	}
	return true, value.Undefined
}

func DiagnosticListContains(v *vm.VM, list *value.DiagnosticList, val value.Value) (bool, value.Value) {
	for _, element := range *list {
		equal, err := v.CallMethodByName(
			symbol.OpEqual,
			value.Ref((*value.Diagnostic)(element)),
			val,
		)
		if !err.IsUndefined() {
			return false, err
		}
		if value.Truthy(equal) {
			return true, value.Undefined
		}
	}
	return false, value.Undefined
}
