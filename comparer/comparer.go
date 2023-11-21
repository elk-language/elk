// Package comparer implements
// comparer functions for Elk values.
package comparer

import (
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

var Comparer cmp.Options

func init() {
	bigFloatComparer := cmp.Comparer(func(x, y *value.BigFloat) bool {
		if x.IsNaN() || y.IsNaN() {
			return x.IsNaN() && y.IsNaN()
		}
		return x.AsGoBigFloat().Cmp(y.AsGoBigFloat()) == 0 &&
			(x.IsInf(0) || y.IsInf(0) || x.Precision() == y.Precision())
	})
	floatComparer := cmp.Comparer(func(x, y value.Float) bool {
		if x.IsNaN() || y.IsNaN() {
			return x.IsNaN() && y.IsNaN()
		}
		return x == y
	})
	float64Comparer := cmp.Comparer(func(x, y value.Float64) bool {
		if x.IsNaN() || y.IsNaN() {
			return x.IsNaN() && y.IsNaN()
		}
		return x == y
	})
	float32Comparer := cmp.Comparer(func(x, y value.Float32) bool {
		if x.IsNaN() || y.IsNaN() {
			return x.IsNaN() && y.IsNaN()
		}
		return x == y
	})

	Comparer = make(cmp.Options, 13)

	Comparer[0] = cmp.AllowUnexported(
		value.Error{},
		value.BigInt{},
		value.Class{},
		bitfield.Bitfield8{},
	)
	Comparer[1] = cmp.AllowUnexported(vm.BytecodeMethod{})
	Comparer[2] = floatComparer
	Comparer[3] = bigFloatComparer
	Comparer[4] = float32Comparer
	Comparer[5] = float64Comparer
	Comparer[6] = value.NewSymbolTableComparer()
	Comparer[7] = vm.NewNativeMethodComparer()
	Comparer[8] = value.NewObjectComparer(Comparer)
	Comparer[9] = value.NewErrorComparer(Comparer)
	Comparer[10] = value.NewClassComparer(Comparer)
	Comparer[11] = value.NewMixinComparer(Comparer)
	Comparer[12] = value.NewModuleComparer(Comparer)
}
