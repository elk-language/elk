package sh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Sh executes a shell program with the given arguments.
// stdout and stderr are shared with the current process.
func Sh(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err == nil {
		return nil
	}

	var cmdBuff strings.Builder
	cmdBuff.WriteString(name)
	for _, arg := range args {
		cmdBuff.WriteByte(' ')
		cmdBuff.WriteString(arg)
	}
	return fmt.Errorf("error executing command: `%s`, %w", cmdBuff.String(), err)
}

// Sh executes a shell program with the given arguments.
// stdout and stderr are set to new buffers, their content
// is returned if an error occurs.
func SilentSh(name string, args ...string) error {
	var cmdStdout strings.Builder
	var cmdStderr strings.Builder
	cmd := exec.Command(name, args...)
	cmd.Stdout = &cmdStdout
	cmd.Stderr = &cmdStderr
	err := cmd.Run()
	if err == nil {
		return nil
	}

	fmt.Println(cmdStdout.String())
	fmt.Println(cmdStderr.String())
	var cmdBuff strings.Builder
	cmdBuff.WriteString(name)
	for _, arg := range args {
		cmdBuff.WriteByte(' ')
		cmdBuff.WriteString(arg)
	}
	return fmt.Errorf("error executing command: `%s`, %w", cmdBuff.String(), err)
}
