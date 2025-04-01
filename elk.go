package elk

import (
	ast "github.com/elk-language/elk/parser/ast/runtime"
	diagnostic "github.com/elk-language/elk/position/diagnostic/runtime"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func InitGlobalEnvironment() {
	value.InitGlobalEnvironment()
	vm.InitGlobalEnvironment()
	ast.InitGlobalEnvironment()
	diagnostic.InitGlobalEnvironment()
}
