//go:build !native

// Package vm contains the Elk Virtual Machine.
// It interprets Elk Bytecode produced by
// the Elk compiler.
package vm

import (
	"fmt"
	"io"
	"sync/atomic"

	"github.com/elk-language/elk/config"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/value"
	"github.com/fatih/color"
)

var INIT_VALUE_STACK_SIZE int
var MAX_VALUE_STACK_SIZE int
var DefaultThreadPool = &ThreadPool{}
var DEFAULT_THREAD_POOL_SIZE int
var DEFAULT_THREAD_POOL_QUEUE_SIZE int

// Global counter of VM IDs
var currentID atomic.Int64

func init() {
	val, ok := config.IntFromEnvVar("ELK_INIT_VALUE_STACK_SIZE")
	if ok {
		INIT_VALUE_STACK_SIZE = val / int(value.ValueSize)
	} else {
		INIT_VALUE_STACK_SIZE = 24_000 / int(value.ValueSize) // 24KB by default
	}

	val, ok = config.IntFromEnvVar("ELK_MAX_VALUE_STACK_SIZE")
	if ok {
		MAX_VALUE_STACK_SIZE = val / int(value.ValueSize)
	} else {
		MAX_VALUE_STACK_SIZE = 100_000_000 / int(value.ValueSize) // 100MB by default
	}

	val, ok = config.IntFromEnvVar("ELK_DEFAULT_THREAD_POOL_SIZE")
	if ok {
		DEFAULT_THREAD_POOL_SIZE = val
	} else {
		DEFAULT_THREAD_POOL_SIZE = 4
	}
	val, ok = config.IntFromEnvVar("ELK_DEFAULT_THREAD_POOL_QUEUE_SIZE")
	if ok {
		DEFAULT_THREAD_POOL_QUEUE_SIZE = val
	} else {
		DEFAULT_THREAD_POOL_QUEUE_SIZE = 256
	}

	DefaultThreadPool.initThreadPool(DEFAULT_THREAD_POOL_SIZE, DEFAULT_THREAD_POOL_QUEUE_SIZE)
}

type Option func(*VM) // constructor option function

// Assign the given io.Reader as the Stdin of the VM.
func WithStdin(stdin io.Reader) Option {
	return func(vm *VM) {
		vm.Stdin = stdin
	}
}

// Assign the given io.Writer as the Stdout of the VM.
func WithStdout(stdout io.Writer) Option {
	return func(vm *VM) {
		vm.Stdout = stdout
	}
}

// Assign the given io.Writer as the Stderr of the VM.
func WithStderr(stderr io.Writer) Option {
	return func(vm *VM) {
		vm.Stderr = stderr
	}
}

func WithThreadPool(tp *ThreadPool) Option {
	return func(vm *VM) {
		vm.threadPool = tp
	}
}

func PrintError(stderr io.Writer, stackTrace *value.StackTrace, err value.Value) {
	fmt.Fprint(stderr, stackTrace.String())
	c := color.New(color.FgRed, color.Bold)
	if value.IsA(err, value.ErrorClass) {
		errObj := (*value.Object)(err.Pointer())
		c.Fprint(stderr, "Error! Uncaught error ")
		fmt.Fprint(stderr, lexer.Colorize(errObj.Class().Name))
		fmt.Fprint(stderr, ": ")
		fmt.Fprintln(stderr, lexer.ColorizeEmbellishedText(errObj.Message().AsString().String()))
	} else {
		c.Fprint(stderr, "Error! Uncaught thrown value:")
		fmt.Fprint(stderr, " ")
		fmt.Fprintln(stderr, lexer.Colorize(err.Inspect()))
	}

	fmt.Fprintln(stderr)
}

// Get the stored error stack trace.
func (vm *VM) ErrStackTrace() *value.StackTrace {
	if vm.state == errorState {
		return vm.errStackTrace
	}

	return nil
}

func (vm *VM) populateMissingParameters(args []value.Value, paramCount, argumentCount int) []value.Value {
	// populate missing optional arguments with undefined
	missingParams := uintptr(paramCount - argumentCount)
	if missingParams > 0 {
		newArgs := make([]value.Value, paramCount)
		copy(newArgs, args)
		return newArgs
	}

	return args
}

var callSymbol = value.ToSymbol("call")
var toStringSymbol = value.ToSymbol("to_string")
