package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a single statement, so for example
// a single valid "line" of Elk code.
// Usually its an expression optionally terminated with a newline or a semicolon.
type StatementNode interface {
	Node
	statementNode()
}

func (*InvalidNode) statementNode()             {}
func (*ExpressionStatementNode) statementNode() {}
func (*EmptyStatementNode) statementNode()      {}
func (*ImportStatementNode) statementNode()     {}
func (*ParameterStatementNode) statementNode()  {}

// Expression optionally terminated with a newline or a semicolon.
type ExpressionStatementNode struct {
	NodeBase
	Expression ExpressionNode
}

func (e *ExpressionStatementNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ExpressionStatementNode)
	if !ok {
		return false
	}

	return o.span.Equal(o.span) &&
		e.Expression.Equal(value.Ref(o.Expression))
}

// Return a string representation of the node.
func (e *ExpressionStatementNode) String() string {
	return e.Expression.String()
}

func (e *ExpressionStatementNode) IsStatic() bool {
	return e.Expression.IsStatic()
}

func (e *ExpressionStatementNode) Class() *value.Class {
	return value.ExpressionStatementNodeClass
}

func (e *ExpressionStatementNode) DirectClass() *value.Class {
	return value.ExpressionStatementNodeClass
}

func (n *ExpressionStatementNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ExpressionStatementNode{\n  &: %p", n)

	buff.WriteString(",\n  expression: ")
	indent.IndentStringFromSecondLine(&buff, n.Expression.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (e *ExpressionStatementNode) Error() string {
	return e.Inspect()
}

// Create a new expression statement node eg. `5 * 2\n`
func NewExpressionStatementNode(span *position.Span, expr ExpressionNode) *ExpressionStatementNode {
	return &ExpressionStatementNode{
		NodeBase:   NodeBase{span: span},
		Expression: expr,
	}
}

// Same as [NewExpressionStatementNode] but returns an interface
func NewExpressionStatementNodeI(span *position.Span, expr ExpressionNode) StatementNode {
	return &ExpressionStatementNode{
		NodeBase:   NodeBase{span: span},
		Expression: expr,
	}
}

// Represents an empty statement eg. a statement with only a semicolon or a newline.
type EmptyStatementNode struct {
	NodeBase
}

// Check if this node equals another node.
func (n *EmptyStatementNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*EmptyStatementNode)
	if !ok {
		return false
	}

	return n.span.Equal(o.span)
}

// Return a string representation of the node.
func (n *EmptyStatementNode) String() string {
	return ""
}

func (*EmptyStatementNode) IsStatic() bool {
	return false
}

func (e *EmptyStatementNode) Class() *value.Class {
	return value.EmptyStatementNodeClass
}

func (e *EmptyStatementNode) DirectClass() *value.Class {
	return value.EmptyStatementNodeClass
}

func (e *EmptyStatementNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::EmptyStatementNode{&: %p}", e)
}

func (e *EmptyStatementNode) Error() string {
	return e.Inspect()
}

// Create a new empty statement node.
func NewEmptyStatementNode(span *position.Span) *EmptyStatementNode {
	return &EmptyStatementNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents an import statement eg. `import "./foo/bar.elk"`
type ImportStatementNode struct {
	NodeBase
	Path    StringLiteralNode
	FsPaths []string // resolved file system paths
}

// Check if this node equals another node.
func (n *ImportStatementNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ImportStatementNode)
	if !ok {
		return false
	}

	return n.Path.Equal(value.Ref(o.Path)) && n.span.Equal(o.span)
}

// Return a string representation of the node.
func (n *ImportStatementNode) String() string {
	var buff strings.Builder

	buff.WriteString("import ")
	buff.WriteString(n.Path.String())

	return buff.String()
}

func (i *ImportStatementNode) IsStatic() bool {
	return false
}

func (e *ImportStatementNode) Class() *value.Class {
	return value.ImportStatementNodeClass
}

func (e *ImportStatementNode) DirectClass() *value.Class {
	return value.ImportStatementNodeClass
}

func (e *ImportStatementNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::ImportStatementNode{&: %p, path: %s}", e, e.Path.Inspect())
}

func (e *ImportStatementNode) Error() string {
	return e.Inspect()
}

// Create a new import statement node eg. `import "foo"`
func NewImportStatementNode(span *position.Span, path StringLiteralNode) *ImportStatementNode {
	return &ImportStatementNode{
		NodeBase: NodeBase{span: span},
		Path:     path,
	}
}

// Formal parameter optionally terminated with a newline or a semicolon.
type ParameterStatementNode struct {
	NodeBase
	Parameter ParameterNode
}

func (n *ParameterStatementNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ParameterStatementNode)
	if !ok {
		return false
	}

	return n.Parameter.Equal(value.Ref(o.Parameter)) &&
		n.span.Equal(o.span)
}

// Return a string representation of the node.
func (n *ParameterStatementNode) String() string {
	return n.Parameter.String()
}

func (*ParameterStatementNode) IsStatic() bool {
	return false
}

func (*ParameterStatementNode) Class() *value.Class {
	return value.ParameterStatementNodeClass
}

func (*ParameterStatementNode) DirectClass() *value.Class {
	return value.ParameterStatementNodeClass
}

func (n *ParameterStatementNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ParameterStatementNode{\n  &: %p", n)

	buff.WriteString(",\n  parameter: ")
	indent.IndentStringFromSecondLine(&buff, n.Parameter.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (e *ParameterStatementNode) Error() string {
	return e.Inspect()
}

// Create a new formal parameter statement node eg. `foo: Bar\n`
func NewParameterStatementNode(span *position.Span, param ParameterNode) *ParameterStatementNode {
	return &ParameterStatementNode{
		NodeBase:  NodeBase{span: span},
		Parameter: param,
	}
}

// Same as [NewParameterStatementNode] but returns an interface
func NewParameterStatementNodeI(span *position.Span, param ParameterNode) StructBodyStatementNode {
	return &ParameterStatementNode{
		NodeBase:  NodeBase{span: span},
		Parameter: param,
	}
}
