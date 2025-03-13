package elk

import (
	"github.com/elk-language/elk/parser/ast/runtime"
	"github.com/elk-language/elk/vm"
)

func InitGlobalEnvironment() {
	vm.InitGlobalEnvironment()
	runtime.InitGlobalEnvironment()
}
