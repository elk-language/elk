package macros

import (
	"github.com/elk-language/elk/types"
)

// Create a new global environment for type checking
// with native macros.
func NewGlobalEnvironment() *types.GlobalEnvironment {
	env := types.NewGlobalEnvironment()

	initResult(env)

	return env
}
