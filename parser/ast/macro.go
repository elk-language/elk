package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a macro definition eg. `macro foo(a: Elk::AST::StringLiteralNode); a; end`
type MacroDefinitionNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	sealed     bool
	Name       IdentifierNode
	Parameters []ParameterNode // formal parameters
	Body       []StatementNode // body of the method
}

func (n *MacroDefinitionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	params := SpliceSlice(n.Parameters, loc, args, unquote)
	body := SpliceSlice(n.Body, loc, args, unquote)

	return &MacroDefinitionNode{
		TypedNodeBase:          TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		DocCommentableNodeBase: n.DocCommentableNodeBase,
		Name:                   n.Name.splice(loc, args, unquote).(IdentifierNode),
		Parameters:             params,
		Body:                   body,
		sealed:                 n.sealed,
	}
}

func (n *MacroDefinitionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::MacroDefinitionNode", env)
}

func (n *MacroDefinitionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Name.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	for _, param := range n.Parameters {
		if param.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	for _, stmt := range n.Body {
		if stmt.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

// Check if this method definition is equal to another value.
func (n *MacroDefinitionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*MacroDefinitionNode)
	if !ok {
		return false
	}

	if len(n.Parameters) != len(o.Parameters) ||
		len(n.Body) != len(o.Body) {
		return false
	}

	for i, param := range n.Parameters {
		if !param.Equal(value.Ref(o.Parameters[i])) {
			return false
		}
	}

	for i, stmt := range n.Body {
		if !stmt.Equal(value.Ref(o.Body[i])) {
			return false
		}
	}

	return n.loc.Equal(o.loc) &&
		n.Name.Equal(value.Ref(o.Name)) &&
		n.sealed == o.sealed
}

// Return a string representation of this method definition.
func (n *MacroDefinitionNode) String() string {
	var buff strings.Builder

	if n.IsSealed() {
		buff.WriteString("sealed ")
	}
	buff.WriteString("macro ")
	buff.WriteString(n.Name.String())

	buff.WriteString("(")
	for i, param := range n.Parameters {
		if i > 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(param.String())
	}
	buff.WriteString(")")

	buff.WriteRune('\n')
	for _, stmt := range n.Body {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteRune('\n')
	}
	buff.WriteString("end")

	return buff.String()
}

func (*MacroDefinitionNode) IsStatic() bool {
	return false
}

func (m *MacroDefinitionNode) IsSealed() bool {
	return m.sealed
}

func (m *MacroDefinitionNode) SetSealed() {
	m.sealed = true
}

func (*MacroDefinitionNode) Class() *value.Class {
	return value.MacroDefinitionNodeClass
}

func (*MacroDefinitionNode) DirectClass() *value.Class {
	return value.MacroDefinitionNodeClass
}

func (n *MacroDefinitionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::MacroDefinitionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  doc_comment: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	fmt.Fprintf(&buff, ",\n  sealed: %t", n.IsSealed())

	buff.WriteString(",\n  name: ")
	indent.IndentStringFromSecondLine(&buff, n.Name.Inspect(), 1)

	buff.WriteString(",\n  parameters: %[\n")
	for i, element := range n.Parameters {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  body: %[\n")
	for i, element := range n.Body {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (p *MacroDefinitionNode) Error() string {
	return p.Inspect()
}

// Create a method definition node eg. `def foo: String then 'hello world'`
func NewMacroDefinitionNode(
	loc *position.Location,
	docComment string,
	sealed bool,
	name IdentifierNode,
	params []ParameterNode,
	body []StatementNode,
) *MacroDefinitionNode {
	return &MacroDefinitionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		sealed:     sealed,
		Name:       name,
		Parameters: params,
		Body:       body,
	}
}

// Represents a method call eg. `Foo::bar!(5)`
type ScopedMacroCallNode struct {
	TypedNodeBase
	Receiver            ExpressionNode
	MacroName           IdentifierNode
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (n *ScopedMacroCallNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &ScopedMacroCallNode{
		TypedNodeBase:       TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Receiver:            n.Receiver.splice(loc, args, unquote).(ExpressionNode),
		MacroName:           n.MacroName.splice(loc, args, unquote).(IdentifierNode),
		PositionalArguments: SpliceSlice(n.PositionalArguments, loc, args, unquote),
		NamedArguments:      SpliceSlice(n.NamedArguments, loc, args, unquote),
	}
}

func (n *ScopedMacroCallNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::ScopedMacroCallNode", env)
}

func (n *ScopedMacroCallNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.MacroName.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	if n.Receiver.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	for _, arg := range n.PositionalArguments {
		if arg.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	for _, arg := range n.NamedArguments {
		if arg.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *ScopedMacroCallNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ScopedMacroCallNode)
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
		n.MacroName.Equal(value.Ref(o.MacroName))
}

func (n *ScopedMacroCallNode) String() string {
	var buff strings.Builder

	receiverParen := ExpressionPrecedence(n) > ExpressionPrecedence(n.Receiver)

	if receiverParen {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Receiver.String())
	if receiverParen {
		buff.WriteRune(')')
	}

	buff.WriteString("::")
	buff.WriteString(n.MacroName.String())
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

func (*ScopedMacroCallNode) IsStatic() bool {
	return false
}

func (*ScopedMacroCallNode) Class() *value.Class {
	return value.ScopedMacroCallNodeClass
}

func (*ScopedMacroCallNode) DirectClass() *value.Class {
	return value.ScopedMacroCallNodeClass
}

func (n *ScopedMacroCallNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ScopedMacroCallNode{\n  loc: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  receiver: ")
	indent.IndentStringFromSecondLine(&buff, n.Receiver.Inspect(), 1)

	buff.WriteString(",\n  macro_name: ")
	indent.IndentStringFromSecondLine(&buff, n.MacroName.Inspect(), 1)

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

func (n *ScopedMacroCallNode) Error() string {
	return n.Inspect()
}

// Create a macro call node eg. `'123'.to_int!()`
func NewScopedMacroCallNode(loc *position.Location, recv ExpressionNode, macroName IdentifierNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *ScopedMacroCallNode {
	return &ScopedMacroCallNode{
		TypedNodeBase:       TypedNodeBase{loc: loc},
		Receiver:            recv,
		MacroName:           macroName,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a method call eg. `'123'.to_int!()`
type MacroCallNode struct {
	TypedNodeBase
	Receiver            ExpressionNode
	MacroName           IdentifierNode
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (n *MacroCallNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &MacroCallNode{
		TypedNodeBase:       TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Receiver:            n.Receiver.splice(loc, args, unquote).(ExpressionNode),
		MacroName:           n.MacroName.splice(loc, args, unquote).(IdentifierNode),
		PositionalArguments: SpliceSlice(n.PositionalArguments, loc, args, unquote),
		NamedArguments:      SpliceSlice(n.NamedArguments, loc, args, unquote),
	}
}

func (n *MacroCallNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::MacroCallNode", env)
}

func (n *MacroCallNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.MacroName.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	if n.Receiver.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	for _, arg := range n.PositionalArguments {
		if arg.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	for _, arg := range n.NamedArguments {
		if arg.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
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
		n.MacroName.Equal(value.Ref(o.MacroName))
}

func (n *MacroCallNode) String() string {
	var buff strings.Builder

	receiverParen := ExpressionPrecedence(n) > ExpressionPrecedence(n.Receiver)

	if receiverParen {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Receiver.String())
	if receiverParen {
		buff.WriteRune(')')
	}

	buff.WriteRune('.')
	buff.WriteString(n.MacroName.String())
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

	buff.WriteString(",\n  macro_name: ")
	indent.IndentStringFromSecondLine(&buff, n.MacroName.Inspect(), 1)

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
func NewMacroCallNode(loc *position.Location, recv ExpressionNode, macroName IdentifierNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *MacroCallNode {
	return &MacroCallNode{
		TypedNodeBase:       TypedNodeBase{loc: loc},
		Receiver:            recv,
		MacroName:           macroName,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a function-like macro call eg. `foo!(123)`
type ReceiverlessMacroCallNode struct {
	TypedNodeBase
	MacroName           IdentifierNode
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (n *ReceiverlessMacroCallNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &ReceiverlessMacroCallNode{
		TypedNodeBase:       TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		MacroName:           n.MacroName.splice(loc, args, unquote).(IdentifierNode),
		PositionalArguments: SpliceSlice(n.PositionalArguments, loc, args, unquote),
		NamedArguments:      SpliceSlice(n.NamedArguments, loc, args, unquote),
	}
}

func (n *ReceiverlessMacroCallNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::ReceiverlessMacroCallNode", env)
}

func (n *ReceiverlessMacroCallNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, arg := range n.PositionalArguments {
		if arg.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	for _, arg := range n.NamedArguments {
		if arg.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
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

	buff.WriteString(n.MacroName.String())
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
	indent.IndentStringFromSecondLine(&buff, n.MacroName.Inspect(), 1)

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
func NewReceiverlessMacroCallNode(loc *position.Location, macroName IdentifierNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *ReceiverlessMacroCallNode {
	return &ReceiverlessMacroCallNode{
		TypedNodeBase:       TypedNodeBase{loc: loc},
		MacroName:           macroName,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}
