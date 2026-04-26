package elk

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/compiler"
	_ "github.com/elk-language/elk/ext/std"
	"github.com/elk-language/elk/info"
	"github.com/elk-language/elk/lexer"
	lexerRuntime "github.com/elk-language/elk/lexer/runtime"
	astRuntime "github.com/elk-language/elk/parser/ast/runtime"
	parserRuntime "github.com/elk-language/elk/parser/runtime"
	"github.com/elk-language/elk/position/diagnostic"
	diagnosticRuntime "github.com/elk-language/elk/position/diagnostic/runtime"
	"github.com/elk-language/elk/types/checker"
	typesRuntime "github.com/elk-language/elk/types/runtime"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func InitGlobalEnvironment() {
	value.InitGlobalEnvironment()
	vm.InitGlobalEnvironment()
	astRuntime.InitGlobalEnvironment()
	diagnosticRuntime.InitGlobalEnvironment()
	lexerRuntime.InitGlobalEnvironment()
	parserRuntime.InitGlobalEnvironment()
	typesRuntime.InitGlobalEnvironment()
}

func compileResult(buffer *bytes.Buffer, goCompiler *compiler.GoCompiler, diagnostics diagnostic.DiagnosticList) (binPath string, err error) {
	if diagnostics != nil && diagnostics.IsFailure() {
		return "", diagnostics
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	outPath := path.Join(cwd, "out")
	err = os.RemoveAll(outPath)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(outPath, 0755)
	if err != nil {
		return "", err
	}
	targetPath := path.Join(outPath, "main.go")

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return "", err
	}

	goCompiler.Flush()
	result, err := format.Source(buffer.Bytes())
	if err != nil {
		targetFile.Write(result)
		return "", fmt.Errorf("cannot format target go file: %s, %w", targetPath, err)
	}

	_, err = targetFile.Write(result)
	if err != nil {
		return "", err
	}

	var goModBuff bytes.Buffer
	fmt.Fprintf(&goModBuff, "module main\n\n")

	goVersion := strings.TrimPrefix(runtime.Version(), "go")
	fmt.Fprintf(&goModBuff, "go %s\n\n", goVersion)

	goModPath := path.Join(outPath, "go.mod")
	goModFile, err := os.Create(goModPath)
	if err != nil {
		return "", err
	}
	goModFile.Write(goModBuff.Bytes())

	if info.Version == "dev" {
		var goWorkBuff bytes.Buffer
		fmt.Fprintf(&goWorkBuff, "go %s\n\n", goVersion)
		fmt.Fprintf(&goWorkBuff, "use .\n")
		fmt.Fprintf(&goWorkBuff, "use ..\n")

		goWorkPath := path.Join(outPath, "go.work")
		goWorkFile, err := os.Create(goWorkPath)
		if err != nil {
			return "", err
		}
		goWorkFile.Write(goWorkBuff.Bytes())
	} else {
		err = sh("go", "-C", outPath, "get", fmt.Sprintf("github.com/elk-language/elk@%s", info.Version))
		if err != nil {
			return "", err
		}
	}

	err = sh("go", "-C", outPath, "mod", "tidy")
	if err != nil {
		return "", err
	}

	err = sh("go", "-C", outPath, "build", "-tags", "native", "-ldflags", fmt.Sprintf("-X 'github.com/elk-language/elk/info.Version=%s'", info.Version))
	if err != nil {
		return "", err
	}

	return path.Join(outPath, "main"), nil
}

func CompileRunSource(sourceName, source string) (err error) {
	path, err := CompileSource(sourceName, source)
	if err != nil {
		return err
	}

	return sh(path)
}

func CompileSource(sourceName, source string) (binPath string, err error) {
	var buffer bytes.Buffer
	goCompiler, diagnostics := checker.CheckSourceNative(sourceName, source, nil, &buffer, nil)
	return compileResult(
		&buffer,
		goCompiler,
		diagnostics,
	)
}

func CompileFile(fileName string) (binPath string, err error) {
	absFileName, err := filepath.Abs(fileName)
	if err != nil {
		return "", fmt.Errorf("could not find file `%s`", fileName)
	}
	_, err = os.Stat(absFileName)
	if err != nil {
		return "", fmt.Errorf("could not find file `%s`, %w", absFileName, err)
	}

	var buffer bytes.Buffer
	goCompiler, diagnostics := checker.CheckFileNative(fileName, nil, &buffer, nil)
	return compileResult(
		&buffer,
		goCompiler,
		diagnostics,
	)
}

func sh(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err == nil {
		return nil
	}

	var buff strings.Builder
	buff.WriteString(name)
	for _, arg := range args {
		buff.WriteByte(' ')
		buff.WriteString(arg)
	}
	return fmt.Errorf("error executing command: `%s`, %w", buff.String(), err)
}

// Interpret the given file.
// Panics when the file does not compile.
func Interpret(fileName string) (result value.Value, stackTrace *value.StackTrace, err value.Value) {
	bytecode, diagnostics := checker.CheckFile(fileName, nil, bitfield.BitField16{}, nil)
	if diagnostics != nil {
		fmt.Println()

		diagnosticString, err := diagnostics.HumanString(true, lexer.Colorizer{})
		if err != nil {
			panic(err)
		}
		fmt.Println(diagnosticString)
		if diagnostics.IsFailure() {
			panic("failed compilation")
		}
	}

	v := vm.New()
	result, err = v.InterpretTopLevel(bytecode)
	stackTrace = v.ErrStackTrace()
	return result, stackTrace, err
}

// Interpret the given file.
// Panics when the file does not compile or encounters a runtime error
func MustInterpret(fileName string) (result value.Value) {
	result, stackTrace, elkErr := Interpret(fileName)
	if !elkErr.IsUndefined() {
		vm.PrintError(os.Stderr, stackTrace, elkErr)
		panic("failed execution")
	}
	return result
}
