package elk

import (
	_ "github.com/elk-language/elk/ext/std"
	lexer "github.com/elk-language/elk/lexer/runtime"
	ast "github.com/elk-language/elk/parser/ast/runtime"
	parser "github.com/elk-language/elk/parser/runtime"
	diagnostic "github.com/elk-language/elk/position/diagnostic/runtime"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func InitGlobalEnvironment() {
	value.InitGlobalEnvironment()
	vm.InitGlobalEnvironment()
	ast.InitGlobalEnvironment()
	diagnostic.InitGlobalEnvironment()
	lexer.InitGlobalEnvironment()
	parser.InitGlobalEnvironment()
}
