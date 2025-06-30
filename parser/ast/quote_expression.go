package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a quoted piecie of AST
type QuoteExpressionNode struct {
	TypedNodeBase
	Body []StatementNode
}

// Check if the quote block consists of a single expression
// statement, if so retrieves it and returns the expression.
// Otherwise returns nil.
func (n *QuoteExpressionNode) SingleExpression() ExpressionNode {
	var stmtCount int
	var exprStmt *ExpressionStatementNode

	for _, stmt := range n.Body {
		switch stmt := stmt.(type) {
		case *ExpressionStatementNode:
			exprStmt = stmt
			if stmtCount > 0 {
				return nil
			}
			stmtCount++
		case *EmptyStatementNode:
		default:
			if stmtCount > 0 {
				return nil
			}
			stmtCount++
		}
	}

	if stmtCount == 1 && exprStmt != nil {
		return exprStmt.Expression
	}

	return nil
}

func (n *QuoteExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &QuoteExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Body:          SpliceSlice(n.Body, loc, args, unquote),
	}
}

func (n *QuoteExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::QuoteExpressionNode", env)
}

func (n *QuoteExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, stmt := range n.Body {
		if stmt.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

// Check if this node equals another node.
func (n *QuoteExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*QuoteExpressionNode)
	if !ok {
		return false
	}

	if len(n.Body) != len(o.Body) {
		return false
	}

	for i, stmt := range n.Body {
		if !stmt.Equal(value.Ref(o.Body[i])) {
			return false
		}
	}

	return n.loc.Equal(o.loc)
}

// Return a string representation of the node.
func (n *QuoteExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("quote\n")

	for _, stmt := range n.Body {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteString("\n")
	}

	buff.WriteString("end")

	return buff.String()
}

func (*QuoteExpressionNode) IsStatic() bool {
	return false
}

func (*QuoteExpressionNode) Class() *value.Class {
	return value.QuoteExpressionNodeClass
}

func (*QuoteExpressionNode) DirectClass() *value.Class {
	return value.QuoteExpressionNodeClass
}

func (n *QuoteExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::QuoteExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  body: %[\n")
	for i, stmt := range n.Body {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *QuoteExpressionNode) Error() string {
	return n.Inspect()
}

// Create a quote expression node eg.
//
//	quote
//		print("awesome!")
//	end
func NewQuoteExpressionNode(loc *position.Location, body []StatementNode) *QuoteExpressionNode {
	return &QuoteExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Body:          body,
	}
}
