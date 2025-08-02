package parser

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/value"
)

type Result struct {
	AST         *ast.ProgramNode
	Diagnostics diagnostic.DiagnosticList
}

func (*Result) Class() *value.Class {
	return value.ElkParserResultClass
}

func (*Result) DirectClass() *value.Class {
	return value.ElkParserResultClass
}

func (r *Result) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::Parser::Result{\n  &: %p", r)

	buff.WriteString(",\n  ast: ")
	indent.IndentStringFromSecondLine(&buff, r.AST.Inspect(), 1)

	buff.WriteString(",\n  diagnostics: ")
	indent.IndentStringFromSecondLine(&buff, (*value.DiagnosticList)(&r.Diagnostics).Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (r *Result) Error() string {
	return r.Inspect()
}

func (r *Result) SingletonClass() *value.Class {
	return nil
}

func (r *Result) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (r *Result) Copy() value.Reference {
	return &Result{
		AST:         r.AST,
		Diagnostics: r.Diagnostics,
	}
}
