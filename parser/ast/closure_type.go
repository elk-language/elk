package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a closure type eg. `|i: Int|: String`
type ClosureTypeNode struct {
	TypedNodeBase
	Parameters []ParameterNode // formal parameters of the closure separated by semicolons
	ReturnType TypeNode
	ThrowType  TypeNode
}

func (n *ClosureTypeNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	params := SpliceSlice(n.Parameters, loc, args, unquote)

	var returnType TypeNode
	if n.ReturnType != nil {
		returnType = n.ReturnType.splice(loc, args, unquote).(TypeNode)
	}

	var throwType TypeNode
	if n.ThrowType != nil {
		throwType = n.ThrowType.splice(loc, args, unquote).(TypeNode)
	}

	return &ClosureTypeNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Parameters:    params,
		ReturnType:    returnType,
		ThrowType:     throwType,
	}
}

func (n *ClosureTypeNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, param := range n.Parameters {
		if param.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	if n.ReturnType != nil {
		if n.ReturnType.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	if n.ThrowType != nil {
		if n.ThrowType.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

// Check if this node equals another node.
func (n *ClosureTypeNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ClosureTypeNode)
	if !ok {
		return false
	}

	if n.ReturnType == o.ReturnType {
	} else if n.ReturnType == nil || o.ReturnType == nil {
		return false
	} else if !n.ReturnType.Equal(value.Ref(o.ReturnType)) {
		return false
	}

	if n.ThrowType == o.ThrowType {
	} else if n.ThrowType == nil || o.ThrowType == nil {
		return false
	} else if !n.ThrowType.Equal(value.Ref(o.ThrowType)) {
		return false
	}

	if len(n.Parameters) != len(o.Parameters) {
		return false
	}

	for i, param := range n.Parameters {
		if !param.Equal(value.Ref(o.Parameters[i])) {
			return false
		}
	}

	return n.loc.Equal(o.loc)
}

// Return a string representation of the node.
func (n *ClosureTypeNode) String() string {
	var buff strings.Builder

	buff.WriteString("|")
	for i, param := range n.Parameters {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(param.String())
	}
	buff.WriteString("|")

	if n.ReturnType != nil {
		buff.WriteString(": ")

		parens := TypePrecedence(n) > TypePrecedence(n.ReturnType)
		if parens {
			buff.WriteRune('(')
		}
		buff.WriteString(n.ReturnType.String())
		if parens {
			buff.WriteRune(')')
		}
	}

	if n.ThrowType != nil {
		buff.WriteString(" ! ")

		parens := TypePrecedence(n) > TypePrecedence(n.ReturnType)
		if parens {
			buff.WriteRune('(')
		}
		buff.WriteString(n.ThrowType.String())
		if parens {
			buff.WriteRune(')')
		}
	}

	return buff.String()
}

func (*ClosureTypeNode) Class() *value.Class {
	return value.ClosureTypeNodeClass
}

func (*ClosureTypeNode) DirectClass() *value.Class {
	return value.ClosureTypeNodeClass
}

func (n *ClosureTypeNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ClosureTypeNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  return_type: ")
	indent.IndentStringFromSecondLine(&buff, n.ReturnType.Inspect(), 1)

	buff.WriteString(",\n  throw_type: ")
	indent.IndentStringFromSecondLine(&buff, n.ThrowType.Inspect(), 1)

	buff.WriteString(",\n  parameters: %[\n")
	for i, stmt := range n.Parameters {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ClosureTypeNode) Error() string {
	return n.Inspect()
}

func (*ClosureTypeNode) IsStatic() bool {
	return false
}

// Create a new closure type node eg. `|i: Int|: String`
func NewClosureTypeNode(loc *position.Location, params []ParameterNode, retType TypeNode, throwType TypeNode) *ClosureTypeNode {
	return &ClosureTypeNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Parameters:    params,
		ReturnType:    retType,
		ThrowType:     throwType,
	}
}
