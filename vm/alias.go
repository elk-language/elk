package vm

import "github.com/elk-language/elk/value"

// Utility method that defines an alternative name for
// an existing method.
// Panics when the method can't be defined.
func Alias(container *value.MethodContainer, newName, oldName string) {
	err := container.DefineAliasString(newName, oldName)
	if err != nil {
		panic(err)
	}
}
