package compiler_test

import (
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func set[T any](variable *T, value T) T {
	*variable = value
	return value
}

func nativeMethodStr(class *value.Class, name string) *vm.NativeMethod {
	return nativeMethod(class, value.ToSymbol(name))
}

func nativeMethod(class *value.Class, name value.Symbol) *vm.NativeMethod {
	return class.LookupMethod(name).(*vm.NativeMethod)
}

// Create a new position in tests
var P = position.New

// Create a new span in tests
var S = position.NewSpan
var T = token.New
var V = token.NewWithValue

// Create a new source location in tests.
// Create a new location in tests
func L(startPos, endPos *position.Position) *position.Location {
	return position.NewLocation(testFileName, position.NewSpan(startPos, endPos))
}
