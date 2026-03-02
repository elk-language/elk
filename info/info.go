// Package info contains variables with useful info about the elk environment
package info

// Contains the version of Elk currently in use
var Version string

func init() {
	if Version == "" {
		Version = "dev"
	}
}
