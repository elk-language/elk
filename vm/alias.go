package vm

import "github.com/elk-language/elk/value"

// Utility method that defines an alternative name for
// an existing method.
func Alias(namespace *value.Class, newName, oldName string) {
	namespace.DefineAliasString(newName, oldName)
}
