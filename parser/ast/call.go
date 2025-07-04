package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a new expression eg. `new(123)`
type NewExpressionNode struct {
	TypedNodeBase
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (n *NewExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &NewExpressionNode{
		TypedNodeBase:       TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		PositionalArguments: SpliceSlice(n.PositionalArguments, loc, args, unquote),
		NamedArguments:      SpliceSlice(n.NamedArguments, loc, args, unquote),
	}
}

func (n *NewExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::NewExpressionNode", env)
}

func (n *NewExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
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

func (n *NewExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*NewExpressionNode)
	if !ok {
		return false
	}

	if !n.loc.Equal(o.loc) {
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

	return true
}

func (n *NewExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("new(")

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

func (*NewExpressionNode) IsStatic() bool {
	return false
}

// Create a new expression node eg. `new(123)`
func NewNewExpressionNode(loc *position.Location, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *NewExpressionNode {
	return &NewExpressionNode{
		TypedNodeBase:       TypedNodeBase{loc: loc},
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

func (*NewExpressionNode) Class() *value.Class {
	return value.NewExpressionNodeClass
}

func (*NewExpressionNode) DirectClass() *value.Class {
	return value.NewExpressionNodeClass
}

func (n *NewExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::NewExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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

func (n *NewExpressionNode) Error() string {
	return n.Inspect()
}

// Represents a constructor call eg. `ArrayList::[Int](1, 2, 3)`
type GenericConstructorCallNode struct {
	TypedNodeBase
	ClassNode           ComplexConstantNode // class that is being instantiated
	TypeArguments       []TypeNode
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (n *GenericConstructorCallNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &GenericConstructorCallNode{
		TypedNodeBase:       TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		ClassNode:           n.ClassNode.splice(loc, args, unquote).(ComplexConstantNode),
		TypeArguments:       SpliceSlice(n.TypeArguments, loc, args, unquote),
		PositionalArguments: SpliceSlice(n.PositionalArguments, loc, args, unquote),
		NamedArguments:      SpliceSlice(n.NamedArguments, loc, args, unquote),
	}
}

func (n *GenericConstructorCallNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::GenericConstructorCallNode", env)
}

func (n *GenericConstructorCallNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.ClassNode.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	for _, arg := range n.TypeArguments {
		if arg.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
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

// Equal checks if the given GenericConstructorCallNode is equal to another value.
func (n *GenericConstructorCallNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*GenericConstructorCallNode)
	if !ok {
		return false
	}

	if !n.ClassNode.Equal(value.Ref(o.ClassNode)) {
		return false
	}

	if len(n.TypeArguments) != len(o.TypeArguments) ||
		len(n.NamedArguments) != len(o.NamedArguments) ||
		len(n.PositionalArguments) != len(o.PositionalArguments) {
		return false
	}

	for i, arg := range n.TypeArguments {
		if !arg.Equal(value.Ref(o.TypeArguments[i])) {
			return false
		}
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

	return n.loc.Equal(o.loc)
}

// String returns a string representation of the GenericConstructorCallNode.
func (n *GenericConstructorCallNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.ClassNode.String())
	buff.WriteString("::[")

	for i, arg := range n.TypeArguments {
		if i > 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(arg.String())
	}

	buff.WriteString("](")

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

func (*GenericConstructorCallNode) IsStatic() bool {
	return false
}

func (*GenericConstructorCallNode) Class() *value.Class {
	return value.GenericConstructorCallNodeClass
}

func (*GenericConstructorCallNode) DirectClass() *value.Class {
	return value.GenericConstructorCallNodeClass
}

func (n *GenericConstructorCallNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::GenericConstructorCallNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  class_node: ")
	indent.IndentStringFromSecondLine(&buff, n.ClassNode.Inspect(), 1)

	buff.WriteString(",\n  type_arguments: %[\n")
	for i, element := range n.TypeArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

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

func (n *GenericConstructorCallNode) Error() string {
	return n.Inspect()
}

// Create a constructor call node eg. `ArrayList::[Int](1, 2, 3)`
func NewGenericConstructorCallNode(loc *position.Location, class ComplexConstantNode, typeArgs []TypeNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *GenericConstructorCallNode {
	return &GenericConstructorCallNode{
		TypedNodeBase:       TypedNodeBase{loc: loc},
		ClassNode:           class,
		TypeArguments:       typeArgs,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a constructor call eg. `String(123)`
type ConstructorCallNode struct {
	TypedNodeBase
	ClassNode           ComplexConstantNode // class that is being instantiated
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (n *ConstructorCallNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &ConstructorCallNode{
		TypedNodeBase:       TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		ClassNode:           n.ClassNode.splice(loc, args, unquote).(ComplexConstantNode),
		PositionalArguments: SpliceSlice(n.PositionalArguments, loc, args, unquote),
		NamedArguments:      SpliceSlice(n.NamedArguments, loc, args, unquote),
	}
}

func (n *ConstructorCallNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::ConstructorCallNode", env)
}

func (n *ConstructorCallNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.ClassNode.traverse(n, enter, leave) == TraverseBreak {
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

// Check if this node equals another node.
func (n *ConstructorCallNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ConstructorCallNode)
	if !ok {
		return false
	}

	if !n.ClassNode.Equal(value.Ref(o.ClassNode)) {
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

	return n.loc.Equal(o.loc)
}

// Return a string representation of the node.
func (n *ConstructorCallNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.ClassNode.String())
	buff.WriteString("(")

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

func (*ConstructorCallNode) IsStatic() bool {
	return false
}

func (*ConstructorCallNode) Class() *value.Class {
	return value.ConstructorCallNodeClass
}

func (*ConstructorCallNode) DirectClass() *value.Class {
	return value.ConstructorCallNodeClass
}

func (n *ConstructorCallNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ConstructorCallNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  class_node: ")
	indent.IndentStringFromSecondLine(&buff, n.ClassNode.Inspect(), 1)

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

func (n *ConstructorCallNode) Error() string {
	return n.Inspect()
}

// Create a constructor call node eg. `String(123)`
func NewConstructorCallNode(loc *position.Location, class ComplexConstantNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *ConstructorCallNode {
	return &ConstructorCallNode{
		TypedNodeBase:       TypedNodeBase{loc: loc},
		ClassNode:           class,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents attribute access eg. `foo.bar`
type AttributeAccessNode struct {
	TypedNodeBase
	Receiver      ExpressionNode
	AttributeName IdentifierNode
}

func (n *AttributeAccessNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &AttributeAccessNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Receiver:      n.Receiver.splice(loc, args, unquote).(ExpressionNode),
		AttributeName: n.AttributeName.splice(loc, args, unquote).(IdentifierNode),
	}
}

func (n *AttributeAccessNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::AttributeAccessNode", env)
}

func (n *AttributeAccessNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Receiver.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	if n.AttributeName.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (*AttributeAccessNode) IsStatic() bool {
	return false
}

func (*AttributeAccessNode) Class() *value.Class {
	return value.AttributeAccessNodeClass
}

func (*AttributeAccessNode) DirectClass() *value.Class {
	return value.AttributeAccessNodeClass
}
func (n *AttributeAccessNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::AttributeAccessNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  receiver: ")
	indent.IndentStringFromSecondLine(&buff, n.Receiver.Inspect(), 1)

	buff.WriteString(",\n  attribute_name: ")
	indent.IndentStringFromSecondLine(&buff, n.AttributeName.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *AttributeAccessNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*AttributeAccessNode)
	if !ok {
		return false
	}

	return n.Span().Equal(o.Span()) &&
		n.Receiver.Equal(value.Ref(o.Receiver)) &&
		n.AttributeName.Equal(value.Ref(o.AttributeName))
}

func (n *AttributeAccessNode) String() string {
	var buff strings.Builder

	parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Receiver)
	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Receiver.String())
	if parens {
		buff.WriteRune(')')
	}

	buff.WriteString(".")
	buff.WriteString(n.AttributeName.String())

	return buff.String()
}

func (n *AttributeAccessNode) Error() string {
	return n.Inspect()
}

// Create an attribute access node eg. `foo.bar`
func NewAttributeAccessNode(loc *position.Location, recv ExpressionNode, attrName IdentifierNode) *AttributeAccessNode {
	return &AttributeAccessNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Receiver:      recv,
		AttributeName: attrName,
	}
}

// Represents subscript access eg. `arr[5]`
type SubscriptExpressionNode struct {
	TypedNodeBase
	Receiver ExpressionNode
	Key      ExpressionNode
	static   static
}

func (n *SubscriptExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	receiver := n.Receiver.splice(loc, args, unquote).(ExpressionNode)
	key := n.Key.splice(loc, args, unquote).(ExpressionNode)

	return &SubscriptExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Receiver:      receiver,
		Key:           key,
	}
}

func (n *SubscriptExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::SubscriptExpressionNode", env)
}

func (n *SubscriptExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Receiver.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	if n.Key.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *SubscriptExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SubscriptExpressionNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Receiver.Equal(value.Ref(o.Receiver)) &&
		n.Key.Equal(value.Ref(o.Key))
}

func (n *SubscriptExpressionNode) String() string {
	var buff strings.Builder

	parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Receiver)
	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Receiver.String())
	if parens {
		buff.WriteRune(')')
	}

	buff.WriteRune('[')
	buff.WriteString(n.Key.String())
	buff.WriteRune(']')

	return buff.String()
}

func (s *SubscriptExpressionNode) IsStatic() bool {
	if s.static == staticUnset {
		if areExpressionsStatic(s.Key, s.Receiver) {
			s.static = staticTrue
		} else {
			s.static = staticFalse
		}
	}
	return s.static == staticTrue
}

func (*SubscriptExpressionNode) Class() *value.Class {
	return value.SubscriptExpressionNodeClass
}

func (*SubscriptExpressionNode) DirectClass() *value.Class {
	return value.SubscriptExpressionNodeClass
}

func (n *SubscriptExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SubscriptExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  receiver: ")
	indent.IndentStringFromSecondLine(&buff, n.Receiver.Inspect(), 1)

	buff.WriteString(",\n  key: ")
	indent.IndentStringFromSecondLine(&buff, n.Key.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *SubscriptExpressionNode) Error() string {
	return n.Inspect()
}

// Create a subscript expression node eg. `arr[5]`
func NewSubscriptExpressionNode(loc *position.Location, recv, key ExpressionNode) *SubscriptExpressionNode {
	return &SubscriptExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Receiver:      recv,
		Key:           key,
	}
}

// Represents nil-safe subscript access eg. `arr?[5]`
type NilSafeSubscriptExpressionNode struct {
	TypedNodeBase
	Receiver ExpressionNode
	Key      ExpressionNode
	static   static
}

func (n *NilSafeSubscriptExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	receiver := n.Receiver.splice(loc, args, unquote).(ExpressionNode)
	key := n.Key.splice(loc, args, unquote).(ExpressionNode)

	return &NilSafeSubscriptExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Receiver:      receiver,
		Key:           key,
	}
}

func (n *NilSafeSubscriptExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::NilSafeSubscriptExpressionNode", env)
}

func (n *NilSafeSubscriptExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Receiver.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	if n.Key.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *NilSafeSubscriptExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*NilSafeSubscriptExpressionNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Receiver.Equal(value.Ref(o.Receiver)) &&
		n.Key.Equal(value.Ref(o.Key))
}

func (n *NilSafeSubscriptExpressionNode) String() string {
	var buff strings.Builder

	parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Receiver)
	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Receiver.String())
	if parens {
		buff.WriteRune(')')
	}

	buff.WriteString("?[")
	buff.WriteString(n.Key.String())
	buff.WriteString("]")

	return buff.String()
}

func (s *NilSafeSubscriptExpressionNode) IsStatic() bool {
	if s.static == staticUnset {
		if areExpressionsStatic(s.Key, s.Receiver) {
			s.static = staticTrue
		} else {
			s.static = staticFalse
		}
	}
	return s.static == staticTrue
}

func (*NilSafeSubscriptExpressionNode) Class() *value.Class {
	return value.SubscriptExpressionNodeClass
}

func (*NilSafeSubscriptExpressionNode) DirectClass() *value.Class {
	return value.SubscriptExpressionNodeClass
}

func (n *NilSafeSubscriptExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::NilSafeSubscriptExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  receiver: ")
	indent.IndentStringFromSecondLine(&buff, n.Receiver.Inspect(), 1)

	buff.WriteString(",\n  key: ")
	indent.IndentStringFromSecondLine(&buff, n.Key.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *NilSafeSubscriptExpressionNode) Error() string {
	return n.Inspect()
}

// Create a nil-safe subscript expression node eg. `arr?[5]`
func NewNilSafeSubscriptExpressionNode(loc *position.Location, recv, key ExpressionNode) *NilSafeSubscriptExpressionNode {
	return &NilSafeSubscriptExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Receiver:      recv,
		Key:           key,
	}
}

// Represents a method call eg. `'123'.()`
type CallNode struct {
	TypedNodeBase
	Receiver            ExpressionNode
	NilSafe             bool
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (n *CallNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &CallNode{
		TypedNodeBase:       TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Receiver:            n.Receiver.splice(loc, args, unquote).(ExpressionNode),
		NilSafe:             n.NilSafe,
		PositionalArguments: SpliceSlice(n.PositionalArguments, loc, args, unquote),
		NamedArguments:      SpliceSlice(n.NamedArguments, loc, args, unquote),
	}
}

func (n *CallNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::CallNode", env)
}

func (n *CallNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
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

func (n *CallNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*CallNode)
	if !ok {
		return false
	}

	if n.NilSafe != o.NilSafe {
		return false
	}

	if !n.Receiver.Equal(value.Ref(o.Receiver)) {
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

	return n.loc.Equal(o.loc)
}

func (n *CallNode) String() string {
	var buff strings.Builder

	receiverParen := ExpressionPrecedence(n) > ExpressionPrecedence(n.Receiver)

	if receiverParen {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Receiver.String())
	if receiverParen {
		buff.WriteRune(')')
	}

	if n.NilSafe {
		buff.WriteRune('?')
	}
	buff.WriteString(".(")

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

func (*CallNode) IsStatic() bool {
	return false
}

func (*CallNode) Class() *value.Class {
	return value.CallNodeClass
}

func (*CallNode) DirectClass() *value.Class {
	return value.CallNodeClass
}

func (n *CallNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::CallNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  receiver: ")
	indent.IndentStringFromSecondLine(&buff, n.Receiver.Inspect(), 1)

	fmt.Fprintf(&buff, ",\n  nil_safe: %t", n.NilSafe)

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

func (n *CallNode) Error() string {
	return n.Inspect()
}

// Create a method call node eg. `'123'.to_int()`
func NewCallNode(loc *position.Location, recv ExpressionNode, nilSafe bool, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *CallNode {
	return &CallNode{
		TypedNodeBase:       TypedNodeBase{loc: loc},
		Receiver:            recv,
		NilSafe:             nilSafe,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a method call eg. `foo.bar::[String](a)`
type GenericMethodCallNode struct {
	TypedNodeBase
	Receiver            ExpressionNode
	Op                  *token.Token
	MethodName          IdentifierNode
	TypeArguments       []TypeNode
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
	TailCall            bool
}

func (n *GenericMethodCallNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &GenericMethodCallNode{
		TypedNodeBase:       TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Receiver:            n.Receiver.splice(loc, args, unquote).(ExpressionNode),
		Op:                  n.Op.Splice(loc, unquote),
		MethodName:          n.MethodName.splice(loc, args, unquote).(IdentifierNode),
		TypeArguments:       SpliceSlice(n.TypeArguments, loc, args, unquote),
		PositionalArguments: SpliceSlice(n.PositionalArguments, loc, args, unquote),
		NamedArguments:      SpliceSlice(n.NamedArguments, loc, args, unquote),
		TailCall:            n.TailCall,
	}
}

func (n *GenericMethodCallNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::GenericMethodCallNode", env)
}

func (n *GenericMethodCallNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.MethodName.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	if n.Receiver.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	for _, arg := range n.TypeArguments {
		if arg.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
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

// Equal checks if the given GenericMethodCallNode is equal to another value.
func (n *GenericMethodCallNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*GenericMethodCallNode)
	if !ok {
		return false
	}

	if !n.MethodName.Equal(value.Ref(o.MethodName)) ||
		!n.Receiver.Equal(value.Ref(o.Receiver)) ||
		!n.Op.Equal(o.Op) ||
		n.loc.Equal(o.loc) {
		return false
	}

	if len(n.TypeArguments) != len(o.TypeArguments) ||
		len(n.NamedArguments) != len(o.NamedArguments) ||
		len(n.PositionalArguments) != len(o.PositionalArguments) {
		return false
	}

	for i, arg := range n.TypeArguments {
		if !arg.Equal(value.Ref(o.TypeArguments[i])) {
			return false
		}
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

// String returns a string representation of the GenericMethodCallNode.
func (n *GenericMethodCallNode) String() string {
	var buff strings.Builder

	receiverParen := ExpressionPrecedence(n) > ExpressionPrecedence(n.Receiver)

	if receiverParen {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Receiver.String())
	if receiverParen {
		buff.WriteRune(')')
	}

	buff.WriteString(n.Op.String())
	buff.WriteString(n.MethodName.String())
	buff.WriteString("::[")

	for i, arg := range n.TypeArguments {
		if i > 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(arg.String())
	}

	buff.WriteString("](")

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

func (*GenericMethodCallNode) IsStatic() bool {
	return false
}

func (*GenericMethodCallNode) Class() *value.Class {
	return value.GenericMethodCallNodeClass
}

func (*GenericMethodCallNode) DirectClass() *value.Class {
	return value.GenericMethodCallNodeClass
}

func (n *GenericMethodCallNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::GenericMethodCallNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  receiver: ")
	indent.IndentStringFromSecondLine(&buff, n.Receiver.Inspect(), 1)

	buff.WriteString(",\n  op: ")
	indent.IndentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString(",\n  method_name: ")
	indent.IndentStringFromSecondLine(&buff, n.MethodName.Inspect(), 1)

	buff.WriteString(",\n  type_arguments: %[\n")
	for i, element := range n.TypeArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

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

func (n *GenericMethodCallNode) Error() string {
	return n.Inspect()
}

// Create a method call node eg. `foo.bar::[String](a)`
func NewGenericMethodCallNode(loc *position.Location, recv ExpressionNode, op *token.Token, methodName IdentifierNode, typeArgs []TypeNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *GenericMethodCallNode {
	return &GenericMethodCallNode{
		TypedNodeBase:       TypedNodeBase{loc: loc},
		Receiver:            recv,
		Op:                  op,
		MethodName:          methodName,
		TypeArguments:       typeArgs,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a method call eg. `'123'.to_int()`
type MethodCallNode struct {
	TypedNodeBase
	Receiver            ExpressionNode
	Op                  *token.Token
	MethodName          IdentifierNode
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
	TailCall            bool
}

func (n *MethodCallNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &MethodCallNode{
		TypedNodeBase:       TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Receiver:            n.Receiver.splice(loc, args, unquote).(ExpressionNode),
		Op:                  n.Op.Splice(loc, unquote),
		MethodName:          n.MethodName.splice(loc, args, unquote).(IdentifierNode),
		PositionalArguments: SpliceSlice(n.PositionalArguments, loc, args, unquote),
		NamedArguments:      SpliceSlice(n.NamedArguments, loc, args, unquote),
		TailCall:            n.TailCall,
	}
}

func (n *MethodCallNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::MethodCallNode", env)
}

func (n *MethodCallNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.MethodName.traverse(n, enter, leave) == TraverseBreak {
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

func (n *MethodCallNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*MethodCallNode)
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
		n.MethodName.Equal(value.Ref(o.MethodName))
}

func (n *MethodCallNode) String() string {
	var buff strings.Builder

	receiverParen := ExpressionPrecedence(n) > ExpressionPrecedence(n.Receiver)

	if receiverParen {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Receiver.String())
	if receiverParen {
		buff.WriteRune(')')
	}

	buff.WriteString(n.Op.String())
	buff.WriteString(n.MethodName.String())
	buff.WriteString("(")

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

func (*MethodCallNode) IsStatic() bool {
	return false
}

func (*MethodCallNode) Class() *value.Class {
	return value.MethodCallNodeClass
}

func (*MethodCallNode) DirectClass() *value.Class {
	return value.MethodCallNodeClass
}

func (n *MethodCallNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::MethodCallNode{\n  loc: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  receiver: ")
	indent.IndentStringFromSecondLine(&buff, n.Receiver.Inspect(), 1)

	buff.WriteString(",\n  op: ")
	indent.IndentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString(",\n  method_name: ")
	indent.IndentStringFromSecondLine(&buff, n.MethodName.Inspect(), 1)

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

func (n *MethodCallNode) Error() string {
	return n.Inspect()
}

// Create a method call node eg. `'123'.to_int()`
func NewMethodCallNode(loc *position.Location, recv ExpressionNode, op *token.Token, methodName IdentifierNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *MethodCallNode {
	return &MethodCallNode{
		TypedNodeBase:       TypedNodeBase{loc: loc},
		Receiver:            recv,
		Op:                  op,
		MethodName:          methodName,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a function-like call eg. `to_string(123)`
type ReceiverlessMethodCallNode struct {
	TypedNodeBase
	MethodName          IdentifierNode
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
	TailCall            bool
}

func (n *ReceiverlessMethodCallNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &ReceiverlessMethodCallNode{
		TypedNodeBase:       TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		MethodName:          n.MethodName.splice(loc, args, unquote).(IdentifierNode),
		PositionalArguments: SpliceSlice(n.PositionalArguments, loc, args, unquote),
		NamedArguments:      SpliceSlice(n.NamedArguments, loc, args, unquote),
		TailCall:            n.TailCall,
	}
}

func (n *ReceiverlessMethodCallNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::ReceiverlessMethodCallNode", env)
}

func (n *ReceiverlessMethodCallNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.MethodName.traverse(n, enter, leave) == TraverseBreak {
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

func (n *ReceiverlessMethodCallNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ReceiverlessMethodCallNode)
	if !ok {
		return false
	}

	if !n.MethodName.Equal(value.Ref(o.MethodName)) ||
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

func (n *ReceiverlessMethodCallNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.MethodName.String())
	buff.WriteString("(")

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

func (*ReceiverlessMethodCallNode) IsStatic() bool {
	return false
}

func (*ReceiverlessMethodCallNode) Class() *value.Class {
	return value.ReceiverlessMethodCallNodeClass
}

func (*ReceiverlessMethodCallNode) DirectClass() *value.Class {
	return value.ReceiverlessMethodCallNodeClass
}

func (n *ReceiverlessMethodCallNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ReceiverlessMethodCallNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  method_name: ")
	indent.IndentStringFromSecondLine(&buff, n.MethodName.Inspect(), 1)

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

func (n *ReceiverlessMethodCallNode) Error() string {
	return n.Inspect()
}

// Create a function call node eg. `to_string(123)`
func NewReceiverlessMethodCallNode(loc *position.Location, methodName IdentifierNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *ReceiverlessMethodCallNode {
	return &ReceiverlessMethodCallNode{
		TypedNodeBase:       TypedNodeBase{loc: loc},
		MethodName:          methodName,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a generic function-like call eg. `foo::[Int](123)`
type GenericReceiverlessMethodCallNode struct {
	TypedNodeBase
	MethodName          IdentifierNode
	TypeArguments       []TypeNode
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
	TailCall            bool
}

func (n *GenericReceiverlessMethodCallNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &GenericReceiverlessMethodCallNode{
		TypedNodeBase:       TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		MethodName:          n.MethodName.splice(loc, args, unquote).(IdentifierNode),
		TypeArguments:       SpliceSlice(n.TypeArguments, loc, args, unquote),
		PositionalArguments: SpliceSlice(n.PositionalArguments, loc, args, unquote),
		NamedArguments:      SpliceSlice(n.NamedArguments, loc, args, unquote),
		TailCall:            n.TailCall,
	}
}

func (n *GenericReceiverlessMethodCallNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::GenericReceiverlessMethodCallNode", env)
}

func (n *GenericReceiverlessMethodCallNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.MethodName.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	for _, arg := range n.TypeArguments {
		if arg.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
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

// Equal checks if the given GenericReceiverlessMethodCallNode is equal to another value.
func (n *GenericReceiverlessMethodCallNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*GenericReceiverlessMethodCallNode)
	if !ok {
		return false
	}

	if !n.MethodName.Equal(value.Ref(o.MethodName)) || !n.loc.Equal(o.loc) {
		return false
	}

	if len(n.TypeArguments) != len(o.TypeArguments) ||
		len(n.NamedArguments) != len(o.NamedArguments) ||
		len(n.PositionalArguments) != len(o.PositionalArguments) {
		return false
	}

	for i, arg := range n.TypeArguments {
		if !arg.Equal(value.Ref(o.TypeArguments[i])) {
			return false
		}
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

// String returns a string representation of the GenericReceiverlessMethodCallNode.
func (n *GenericReceiverlessMethodCallNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.MethodName.String())
	buff.WriteString("::[")

	for i, arg := range n.TypeArguments {
		if i > 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(arg.String())
	}

	buff.WriteString("](")

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

func (*GenericReceiverlessMethodCallNode) IsStatic() bool {
	return false
}

func (*GenericReceiverlessMethodCallNode) Class() *value.Class {
	return value.GenericReceiverlessMethodCallNodeClass
}

func (*GenericReceiverlessMethodCallNode) DirectClass() *value.Class {
	return value.GenericReceiverlessMethodCallNodeClass
}

func (n *GenericReceiverlessMethodCallNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::GenericReceiverlessMethodCallNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  method_name: ")
	indent.IndentStringFromSecondLine(&buff, n.MethodName.Inspect(), 1)

	buff.WriteString(",\n  type_arguments: %[\n")
	for i, element := range n.TypeArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

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

func (n *GenericReceiverlessMethodCallNode) Error() string {
	return n.Inspect()
}

// Create a generic function call node eg. `foo::[Int](123)`
func NewGenericReceiverlessMethodCallNode(loc *position.Location, methodName IdentifierNode, typeArgs []TypeNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *GenericReceiverlessMethodCallNode {
	return &GenericReceiverlessMethodCallNode{
		TypedNodeBase:       TypedNodeBase{loc: loc},
		MethodName:          methodName,
		TypeArguments:       typeArgs,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}
