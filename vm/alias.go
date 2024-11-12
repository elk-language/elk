package vm

import "github.com/elk-language/elk/value"

// Utility method that defines an alternative name for
// an existing method.
func Alias(container *value.MethodContainer, newName, oldName string) {
	container.DefineAliasString(newName, oldName)
}
