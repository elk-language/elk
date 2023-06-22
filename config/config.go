// Package config provides utilities
// for reading configuration options
// of the interpreter.
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Read an int value from an ENV variable.
// Panic if the value is not a valid int.
// Return false when the variable is not present.
func IntFromEnvVar(varName string) (int, bool) {
	val, ok := os.LookupEnv(varName)
	if !ok {
		return 0, false
	}

	valInt, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Sprintf("invalid value for %s, expected int, got %v", varName, val))
	}

	return valInt, true
}
