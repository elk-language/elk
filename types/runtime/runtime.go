package runtime

import "github.com/elk-language/elk/types/checker/runtime"

func InitGlobalEnvironment() {
	runtime.InitGlobalEnvironment()
	// TODO: expose Elk types
}

func init() {
	InitGlobalEnvironment()
}
