package runtime

import (
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

// Std::Elk::Lexer
func initLexer() {
	// Singleton methods
	c := &value.ElkLexerClass.SingletonClass().MethodContainer
	vm.Def(
		c,
		"colorize",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			source := string(args[1].AsReference().(value.String))
			result := lexer.Colorize(source)
			return value.Ref(value.String(result)), value.Undefined
		},
		vm.DefWithParameters(1),
	)
	vm.Def(
		c,
		"lex",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			source := string(args[1].AsReference().(value.String))
			result := lexer.LexValue(source)
			return value.Ref(result), value.Undefined
		},
		vm.DefWithParameters(2),
	)

	// Instance methods
	c = &value.ElkLexerClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			source := string(args[1].AsReference().(value.String))

			var sourceName string
			if args[2].IsUndefined() {
				sourceName = "<main>"
			} else {
				sourceName = string(args[2].AsReference().(value.String))
			}

			lexer := lexer.NewWithName(sourceName, source)
			return value.Ref(lexer), value.Undefined
		},
		vm.DefWithParameters(2),
	)

	vm.Def(
		c,
		"next",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := (*lexer.Lexer)(args[0].Pointer())
			return value.Ref(self.Next()), value.Undefined
		},
	)
}
