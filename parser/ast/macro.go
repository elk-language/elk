package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Represents a method call eg. `'123'.to_int!()`
type MacroCallNode struct {
	TypedNodeBase
	Receiver            ExpressionNode
	Op                  *token.Token
	MacroName           string
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (n *MacroCallNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &MacroCallNode{
		TypedNodeBase:       TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Receiver:            n.Receiver.Splice(loc, args, unquote).(ExpressionNode),
		Op:                  n.Op.Splice(loc, unquote),
		MacroName:           n.MacroName,
		PositionalArguments: SpliceSlice(n.PositionalArguments, loc, args, unquote),
		NamedArguments:      SpliceSlice(n.NamedArguments, loc, args, unquote),
	}
}

func (n *MacroCallNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*MacroCallNode)
	if !ok {
		return false
	}

	if len(n.PositionalArguments) != len(o.PositionalArguments) ||
		len(n.NamedArguments) != len(o.NamedArguments) {
		return false
	}

	for i, arg := range n.PositionalArguments {
		if !arg.Equal(value.Ref(o.PositionalArguments[i])) {
			return false
		}
	}

	for i, arg := range n.NamedArguments {
		if !arg.Equal(value.Ref(o.NamedArguments[i])) {
			return false
		}
	}

	return n.loc.Equal(o.loc) &&
		n.Receiver.Equal(value.Ref(o.Receiver)) &&
		n.Op.Equal(o.Op) &&
		n.MacroName != o.MacroName
}

func (n *MacroCallNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.Receiver.String())
	buff.WriteString(n.Op.String())
	buff.WriteString(n.MacroName)
	buff.WriteString("!(")

	var hasMultilineArgs bool
	argCount := len(n.PositionalArguments) + len(n.NamedArguments)
	argStrings := make([]string, 0, argCount)

	for _, arg := range n.PositionalArguments {
		argString := arg.String()
		argStrings = append(argStrings, argString)
		if strings.ContainsRune(argString, '\n') {
			hasMultilineArgs = true
		}
	}
	for _, arg := range n.NamedArguments {
		argString := arg.String()
		argStrings = append(argStrings, argString)
		if strings.ContainsRune(argString, '\n') {
			hasMultilineArgs = true
		}
	}

	if hasMultilineArgs || argCount > 6 {
		buff.WriteRune('\n')
		for i, argStr := range argStrings {
			if i != 0 {
				buff.WriteString(",\n")
			}
			indent.IndentString(&buff, argStr, 1)
		}
		buff.WriteRune('\n')
	} else {
		for i, argStr := range argStrings {
			if i != 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(argStr)
		}
	}
	buff.WriteString(")")

	return buff.String()
}

func (*MacroCallNode) IsStatic() bool {
	return false
}

func (*MacroCallNode) Class() *value.Class {
	return value.MacroCallNodeClass
}

func (*MacroCallNode) DirectClass() *value.Class {
	return value.MacroCallNodeClass
}

func (n *MacroCallNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::MacroCallNode{\n  loc: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  receiver: ")
	indent.IndentStringFromSecondLine(&buff, n.Receiver.Inspect(), 1)

	buff.WriteString(",\n  op: ")
	indent.IndentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString(",\n  macro_name: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.MacroName).Inspect(), 1)

	buff.WriteString(",\n  positional_arguments: %[\n")
	for i, element := range n.PositionalArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  named_arguments: %[\n")
	for i, element := range n.NamedArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *MacroCallNode) Error() string {
	return n.Inspect()
}

// Create a macro call node eg. `'123'.to_int!()`
func NewMacroCallNode(loc *position.Location, recv ExpressionNode, op *token.Token, macroName string, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *MacroCallNode {
	return &MacroCallNode{
		TypedNodeBase:       TypedNodeBase{loc: loc},
		Receiver:            recv,
		Op:                  op,
		MacroName:           macroName,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a function-like macro call eg. `foo!(123)`
type ReceiverlessMacroCallNode struct {
	TypedNodeBase
	MacroName           string
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (n *ReceiverlessMacroCallNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &ReceiverlessMacroCallNode{
		TypedNodeBase:       TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		MacroName:           n.MacroName,
		PositionalArguments: SpliceSlice(n.PositionalArguments, loc, args, unquote),
		NamedArguments:      SpliceSlice(n.NamedArguments, loc, args, unquote),
	}
}

func (n *ReceiverlessMacroCallNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ReceiverlessMacroCallNode)
	if !ok {
		return false
	}

	if n.MacroName != o.MacroName ||
		len(n.PositionalArguments) != len(o.PositionalArguments) ||
		len(n.NamedArguments) != len(o.NamedArguments) ||
		!n.loc.Equal(o.loc) {
		return false
	}

	for i, arg := range n.PositionalArguments {
		if !arg.Equal(value.Ref(o.PositionalArguments[i])) {
			return false
		}
	}

	for i, arg := range n.NamedArguments {
		if !arg.Equal(value.Ref(o.NamedArguments[i])) {
			return false
		}
	}

	return true
}

func (n *ReceiverlessMacroCallNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.MacroName)
	buff.WriteString("!(")

	var hasMultilineArgs bool
	argCount := len(n.PositionalArguments) + len(n.NamedArguments)
	argStrings := make([]string, 0, argCount)

	for _, arg := range n.PositionalArguments {
		argString := arg.String()
		argStrings = append(argStrings, argString)
		if strings.ContainsRune(argString, '\n') {
			hasMultilineArgs = true
		}
	}
	for _, arg := range n.NamedArguments {
		argString := arg.String()
		argStrings = append(argStrings, argString)
		if strings.ContainsRune(argString, '\n') {
			hasMultilineArgs = true
		}
	}

	if hasMultilineArgs || argCount > 6 {
		buff.WriteRune('\n')
		for i, argStr := range argStrings {
			if i != 0 {
				buff.WriteString(",\n")
			}
			indent.IndentString(&buff, argStr, 1)
		}
		buff.WriteRune('\n')
	} else {
		for i, argStr := range argStrings {
			if i != 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(argStr)
		}
	}
	buff.WriteString(")")

	return buff.String()
}

func (*ReceiverlessMacroCallNode) IsStatic() bool {
	return false
}

func (*ReceiverlessMacroCallNode) Class() *value.Class {
	return value.ReceiverlessMacroCallNodeClass
}

func (*ReceiverlessMacroCallNode) DirectClass() *value.Class {
	return value.ReceiverlessMacroCallNodeClass
}

func (n *ReceiverlessMacroCallNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ReceiverlessMacroCallNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  macro_name: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.MacroName).Inspect(), 1)

	buff.WriteString(",\n  positional_arguments: %[\n")
	for i, element := range n.PositionalArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  named_arguments: %[\n")
	for i, element := range n.NamedArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ReceiverlessMacroCallNode) Error() string {
	return n.Inspect()
}

// Create a function call node eg. `to_string(123)`
func NewReceiverlessMacroCallNode(loc *position.Location, macroName string, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *ReceiverlessMacroCallNode {
	return &ReceiverlessMacroCallNode{
		TypedNodeBase:       TypedNodeBase{loc: loc},
		MacroName:           macroName,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}
