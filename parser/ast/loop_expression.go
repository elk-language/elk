package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a `loop` expression.
type LoopExpressionNode struct {
	TypedNodeBase
	ThenBody []StatementNode // then expression body
}

func (n *LoopExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &LoopExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		ThenBody:      SpliceSlice(n.ThenBody, loc, args, unquote),
	}
}

func (n *LoopExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::LoopExpressionNode", env)
}

func (n *LoopExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, stmt := range n.ThenBody {
		if stmt.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *LoopExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*LoopExpressionNode)
	if !ok {
		return false
	}

	if len(n.ThenBody) != len(o.ThenBody) {
		return false
	}

	for i, stmt := range n.ThenBody {
		if !stmt.Equal(value.Ref(o.ThenBody[i])) {
			return false
		}
	}

	return n.loc.Equal(o.loc)
}

func (n *LoopExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("loop\n")

	for _, stmt := range n.ThenBody {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteString("\n")
	}

	buff.WriteString("end")

	return buff.String()
}

func (*LoopExpressionNode) IsStatic() bool {
	return false
}

// Create a new `loop` expression node eg. `loop println('elk is awesome')`
func NewLoopExpressionNode(loc *position.Location, then []StatementNode) *LoopExpressionNode {
	return &LoopExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		ThenBody:      then,
	}
}

func (*LoopExpressionNode) Class() *value.Class {
	return value.LoopExpressionNodeClass
}

func (*LoopExpressionNode) DirectClass() *value.Class {
	return value.LoopExpressionNodeClass
}

func (n *LoopExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::LoopExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  then_body: %[\n")
	for i, stmt := range n.ThenBody {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *LoopExpressionNode) Error() string {
	return n.Inspect()
}
