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

func (n *ExpressionStatementNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &ExpressionStatementNode{
		NodeBase:   NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Expression: n.Expression.Splice(loc, args, unquote).(ExpressionNode),
	}
}

func (e *ExpressionStatementNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ExpressionStatementNode)
	if !ok {
		return false
	}

	return o.loc.Equal(o.loc) &&
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

	fmt.Fprintf(&buff, "Std::Elk::AST::ExpressionStatementNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  expression: ")
	indent.IndentStringFromSecondLine(&buff, n.Expression.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (e *ExpressionStatementNode) Error() string {
	return e.Inspect()
}

// Create a new expression statement node eg. `5 * 2\n`
func NewExpressionStatementNode(loc *position.Location, expr ExpressionNode) *ExpressionStatementNode {
	return &ExpressionStatementNode{
		NodeBase:   NodeBase{loc: loc},
		Expression: expr,
	}
}

// Same as [NewExpressionStatementNode] but returns an interface
func NewExpressionStatementNodeI(loc *position.Location, expr ExpressionNode) StatementNode {
	return &ExpressionStatementNode{
		NodeBase:   NodeBase{loc: loc},
		Expression: expr,
	}
}

// Represents an empty statement eg. a statement with only a semicolon or a newline.
type EmptyStatementNode struct {
	NodeBase
}

func (n *EmptyStatementNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &EmptyStatementNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
	}
}

// Check if this node equals another node.
func (n *EmptyStatementNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*EmptyStatementNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc)
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
	return fmt.Sprintf("Std::Elk::AST::EmptyStatementNode{location: %s}", (*value.Location)(e.loc).Inspect())
}

func (e *EmptyStatementNode) Error() string {
	return e.Inspect()
}

// Create a new empty statement node.
func NewEmptyStatementNode(loc *position.Location) *EmptyStatementNode {
	return &EmptyStatementNode{
		NodeBase: NodeBase{loc: loc},
	}
}

// Represents an import statement eg. `import "./foo/bar.elk"`
type ImportStatementNode struct {
	NodeBase
	Path    StringLiteralNode
	FsPaths []string // resolved file system paths
}

func (n *ImportStatementNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &ImportStatementNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Path:     n.Path,
		FsPaths:  n.FsPaths,
	}
}

// Check if this node equals another node.
func (n *ImportStatementNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ImportStatementNode)
	if !ok {
		return false
	}

	return n.Path.Equal(value.Ref(o.Path)) && n.loc.Equal(o.loc)
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
	return fmt.Sprintf("Std::Elk::AST::ImportStatementNode{location: %s, path: %s}", (*value.Location)(e.loc).Inspect(), e.Path.Inspect())
}

func (e *ImportStatementNode) Error() string {
	return e.Inspect()
}

// Create a new import statement node eg. `import "foo"`
func NewImportStatementNode(loc *position.Location, path StringLiteralNode) *ImportStatementNode {
	return &ImportStatementNode{
		NodeBase: NodeBase{loc: loc},
		Path:     path,
	}
}

// Formal parameter optionally terminated with a newline or a semicolon.
type ParameterStatementNode struct {
	NodeBase
	Parameter ParameterNode
}

func (n *ParameterStatementNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &ParameterStatementNode{
		NodeBase:  NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Parameter: n.Parameter.Splice(loc, args, unquote).(ParameterNode),
	}
}

func (n *ParameterStatementNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ParameterStatementNode)
	if !ok {
		return false
	}

	return n.Parameter.Equal(value.Ref(o.Parameter)) &&
		n.loc.Equal(o.loc)
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

	fmt.Fprintf(&buff, "Std::Elk::AST::ParameterStatementNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  parameter: ")
	indent.IndentStringFromSecondLine(&buff, n.Parameter.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (e *ParameterStatementNode) Error() string {
	return e.Inspect()
}

// Create a new formal parameter statement node eg. `foo: Bar\n`
func NewParameterStatementNode(loc *position.Location, param ParameterNode) *ParameterStatementNode {
	return &ParameterStatementNode{
		NodeBase:  NodeBase{loc: loc},
		Parameter: param,
	}
}

// Same as [NewParameterStatementNode] but returns an interface
func NewParameterStatementNodeI(loc *position.Location, param ParameterNode) StructBodyStatementNode {
	return &ParameterStatementNode{
		NodeBase:  NodeBase{loc: loc},
		Parameter: param,
	}
}
