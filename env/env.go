// Package env contains environment variables used by Elk
package env

import (
	"fmt"
	"os"
	"path/filepath"
)

// Path to the directory with the global Elk version
var ELKROOT string

// Path to the directory with the currently used Elk version
var ELKPATH string

// Path to the current Elk executable
var ELKEXEC string

func init() {
	elkExecEnv := os.Getenv("ELKEXEC")
	if elkExecEnv == "" {
		exec, err := os.Executable()
		if err != nil {
			panic(fmt.Sprintf("could not get the elk executable path: %s", err))
		}

		ELKEXEC = exec
	} else {
		ELKEXEC = elkExecEnv
	}

	elkRootEnv := os.Getenv("ELKROOT")
	if elkRootEnv == "" {
		ELKROOT = filepath.Dir(ELKEXEC)
	} else {
		ELKROOT = elkRootEnv
	}

	elkPathEnv := os.Getenv("ELKPATH")
	if elkPathEnv == "" {
		path, err := os.Getwd()
		if err != nil {
			panic(fmt.Sprintf("could not get the working directory: %s", err))
		}
		ELKPATH = path
	} else {
		ELKPATH = elkPathEnv
	}
}
