package vm

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type Colorizer func(string) (string, error)

func (f Colorizer) Colorize(s string) (string, error) {
	return f(s)
}

func MakeColorizer(vm *VM, val value.Value) Colorizer {
	return func(s string) (string, error) {
		result, err := vm.CallMethodByName(symbol.L_colorize, value.Ref(value.String(s)))
		if !err.IsUndefined() {
			return "", err
		}
		return string(result.AsReference().(value.String)), nil
	}
}
