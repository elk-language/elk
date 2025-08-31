package ast

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Check whether the node is a valid right operand of the pipe operator `|>`.
func IsValidPipeExpressionTarget(node Node) bool {
	switch node.(type) {
	case *MethodCallNode, *GenericMethodCallNode, *ReceiverlessMethodCallNode,
		*GenericReceiverlessMethodCallNode, *AttributeAccessNode, *ConstructorCallNode, *CallNode:
		return true
	default:
		return false
	}
}

// Represents a method definition eg. `def foo: String then 'hello world'`
type MethodDefinitionNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Name           IdentifierNode
	TypeParameters []TypeParameterNode
	Parameters     []ParameterNode // formal parameters
	ReturnType     TypeNode
	ThrowType      TypeNode
	Body           []StatementNode // body of the method
	Flags          bitfield.BitField8
}

func (n *MethodDefinitionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	name := n.Name.splice(loc, args, unquote).(IdentifierNode)
	typeParams := SpliceSlice(n.TypeParameters, loc, args, unquote)
	params := SpliceSlice(n.Parameters, loc, args, unquote)

	var returnType TypeNode
	if n.ReturnType != nil {
		returnType = n.ReturnType.splice(loc, args, unquote).(TypeNode)
	}

	var throwType TypeNode
	if n.ThrowType != nil {
		throwType = n.ThrowType.splice(loc, args, unquote).(TypeNode)
	}

	body := SpliceSlice(n.Body, loc, args, unquote)

	return &MethodDefinitionNode{
		TypedNodeBase:          TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		DocCommentableNodeBase: n.DocCommentableNodeBase,
		Name:                   name,
		TypeParameters:         typeParams,
		Parameters:             params,
		ReturnType:             returnType,
		ThrowType:              throwType,
		Body:                   body,
		Flags:                  n.Flags,
	}
}

func (n *MethodDefinitionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::MethodDefinitionNode", env)
}

func (n *MethodDefinitionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Name.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	for _, param := range n.TypeParameters {
		if param.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
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

	for _, stmt := range n.Body {
		if stmt.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

// Check if this method definition is equal to another value.
func (n *MethodDefinitionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*MethodDefinitionNode)
	if !ok {
		return false
	}

	if len(n.TypeParameters) != len(o.TypeParameters) ||
		len(n.Parameters) != len(o.Parameters) ||
		len(n.Body) != len(o.Body) {
		return false
	}

	for i, tp := range n.TypeParameters {
		if !tp.Equal(value.Ref(o.TypeParameters[i])) {
			return false
		}
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

	return n.Flags == o.Flags &&
		n.loc.Equal(o.loc) &&
		n.Name.Equal(value.Ref(o.Name))
}

// Return a string representation of this method definition.
func (n *MethodDefinitionNode) String() string {
	var buff strings.Builder

	if n.IsAbstract() {
		buff.WriteString("abstract ")
	}
	if n.IsSealed() {
		buff.WriteString("sealed ")
	}
	if n.IsGenerator() {
		buff.WriteString("generator ")
	}
	if n.IsAsync() {
		buff.WriteString("async ")
	}
	if n.IsOverload() {
		buff.WriteString("overload ")
	}

	buff.WriteString("def ")
	buff.WriteString(n.Name.String())

	if len(n.TypeParameters) > 0 {
		buff.WriteString("[")
		for i, tp := range n.TypeParameters {
			if i > 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(tp.String())
		}
		buff.WriteString("]")
	}

	buff.WriteString("(")
	for i, param := range n.Parameters {
		if i > 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(param.String())
	}
	buff.WriteString(")")

	if n.ReturnType != nil {
		buff.WriteString(": ")
		buff.WriteString(n.ReturnType.String())
	}

	if n.ThrowType != nil {
		buff.WriteString(" ! ")
		buff.WriteString(n.ThrowType.String())
	}

	buff.WriteRune('\n')
	for _, stmt := range n.Body {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteRune('\n')
	}
	buff.WriteString("end")

	return buff.String()
}

func (*MethodDefinitionNode) IsStatic() bool {
	return false
}

// Whether the method is a setter.
func (m *MethodDefinitionNode) IsSetter() bool {
	return MethodNameIsSetter(m.Name)
}

// Whether the method is a setter.
func MethodNameIsSetter(methodNameNode IdentifierNode) bool {
	var name string

	switch n := methodNameNode.(type) {
	case *PublicIdentifierNode:
		name = n.Value
	case *PrivateIdentifierNode:
		name = n.Value
	default:
		return false
	}

	if len(name) > 0 {
		firstChar, _ := utf8.DecodeRuneInString(name)
		lastChar := name[len(name)-1]
		if (unicode.IsLetter(firstChar) || firstChar == '_') && lastChar == '=' {
			return true
		}
	}

	return false
}

// Whether the method is a setter.
func MethodNameIsSubscriptSetter(methodNameNode IdentifierNode) bool {
	var name string

	switch n := methodNameNode.(type) {
	case *PublicIdentifierNode:
		name = n.Value
	case *PrivateIdentifierNode:
		name = n.Value
	default:
		return false
	}

	return name == "[]="
}

func (m *MethodDefinitionNode) IsAbstract() bool {
	return m.Flags.HasFlag(METHOD_ABSTRACT_FLAG)
}

func (m *MethodDefinitionNode) SetAbstract() {
	m.Flags.SetFlag(METHOD_ABSTRACT_FLAG)
}

func (m *MethodDefinitionNode) IsSealed() bool {
	return m.Flags.HasFlag(METHOD_SEALED_FLAG)
}

func (m *MethodDefinitionNode) SetSealed() {
	m.Flags.SetFlag(METHOD_SEALED_FLAG)
}

func (m *MethodDefinitionNode) IsGenerator() bool {
	return m.Flags.HasFlag(METHOD_GENERATOR_FLAG)
}

func (m *MethodDefinitionNode) SetGenerator() {
	m.Flags.SetFlag(METHOD_GENERATOR_FLAG)
}

func (m *MethodDefinitionNode) IsAsync() bool {
	return m.Flags.HasFlag(METHOD_ASYNC_FLAG)
}

func (m *MethodDefinitionNode) SetAsync() {
	m.Flags.SetFlag(METHOD_ASYNC_FLAG)
}

func (m *MethodDefinitionNode) IsOverload() bool {
	return m.Flags.HasFlag(METHOD_OVERLOAD_FLAG)
}

func (m *MethodDefinitionNode) SetOverload() {
	m.Flags.SetFlag(METHOD_OVERLOAD_FLAG)
}

const (
	METHOD_ABSTRACT_FLAG bitfield.BitFlag8 = 1 << iota
	METHOD_SEALED_FLAG
	METHOD_GENERATOR_FLAG
	METHOD_ASYNC_FLAG
	METHOD_OVERLOAD_FLAG
)

func (*MethodDefinitionNode) Class() *value.Class {
	return value.MethodDefinitionNodeClass
}

func (*MethodDefinitionNode) DirectClass() *value.Class {
	return value.MethodDefinitionNodeClass
}

func (n *MethodDefinitionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::MethodDefinitionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  doc_comment: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	fmt.Fprintf(&buff, ",\n  abstract: %t", n.IsAbstract())
	fmt.Fprintf(&buff, ",\n  sealed: %t", n.IsSealed())
	fmt.Fprintf(&buff, ",\n  generator: %t", n.IsGenerator())
	fmt.Fprintf(&buff, ",\n  async: %t", n.IsAsync())
	fmt.Fprintf(&buff, ",\n  overload: %t", n.IsOverload())

	buff.WriteString(",\n  name: ")
	indent.IndentStringFromSecondLine(&buff, n.Name.Inspect(), 1)

	buff.WriteString(",\n  return_type: ")
	if n.ReturnType == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.ReturnType.Inspect(), 1)
	}

	buff.WriteString(",\n  throw_type: ")
	if n.ThrowType == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.ThrowType.Inspect(), 1)
	}

	buff.WriteString(",\n  type_parameters: %[\n")
	for i, element := range n.TypeParameters {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

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

func (p *MethodDefinitionNode) Error() string {
	return p.Inspect()
}

// Create a method definition node eg. `def foo: String then 'hello world'`
func NewMethodDefinitionNode(
	loc *position.Location,
	docComment string,
	flags bitfield.BitFlag8,
	name IdentifierNode,
	typeParams []TypeParameterNode,
	params []ParameterNode,
	returnType,
	throwType TypeNode,
	body []StatementNode,
) *MethodDefinitionNode {
	return &MethodDefinitionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Flags:          bitfield.BitField8FromBitFlag(flags),
		Name:           name,
		TypeParameters: typeParams,
		Parameters:     params,
		ReturnType:     returnType,
		ThrowType:      throwType,
		Body:           body,
	}
}

// Represents a constructor definition eg. `init then 'hello world'`
type InitDefinitionNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Parameters []ParameterNode // formal parameters
	ThrowType  TypeNode
	Body       []StatementNode // body of the method
}

func (n *InitDefinitionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	params := SpliceSlice(n.Parameters, loc, args, unquote)

	var throwType TypeNode
	if n.ThrowType != nil {
		throwType = n.ThrowType.splice(loc, args, unquote).(TypeNode)
	}

	body := SpliceSlice(n.Body, loc, args, unquote)

	return &InitDefinitionNode{
		TypedNodeBase:          TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		DocCommentableNodeBase: n.DocCommentableNodeBase,
		Parameters:             params,
		ThrowType:              throwType,
		Body:                   body,
	}
}

func (n *InitDefinitionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::InitDefinitionNode", env)
}

func (n *InitDefinitionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
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

	if n.ThrowType != nil {
		if n.ThrowType.traverse(n, enter, leave) == TraverseBreak {
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

// Check if this node equals another node.
func (n *InitDefinitionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*InitDefinitionNode)
	if !ok {
		return false
	}

	if !n.loc.Equal(o.loc) ||
		n.comment != o.comment {
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

	if n.ThrowType == o.ThrowType {
	} else if n.ThrowType == nil || o.ThrowType == nil {
		return false
	} else if !n.ThrowType.Equal(value.Ref(o.ThrowType)) {
		return false
	}

	return true
}

// Return a string representation of the node.
func (n *InitDefinitionNode) String() string {
	var buff strings.Builder

	doc := n.DocComment()
	if len(doc) > 0 {
		buff.WriteString("##[\n")
		indent.IndentString(&buff, doc, 1)
		buff.WriteString("\n]##\n")
	}

	buff.WriteString("init")

	if len(n.Parameters) > 0 {
		buff.WriteString("(")
		for i, param := range n.Parameters {
			if i != 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(param.String())
		}
		buff.WriteString(")")
	}

	if n.ThrowType != nil {
		buff.WriteString(" ! ")
		buff.WriteString(n.ThrowType.String())
	}
	buff.WriteRune('\n')

	for _, stmt := range n.Body {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteRune('\n')
	}
	buff.WriteString("end")

	return buff.String()
}

func (*InitDefinitionNode) IsStatic() bool {
	return false
}

func (*InitDefinitionNode) Class() *value.Class {
	return value.InitDefinitionNodeClass
}

func (*InitDefinitionNode) DirectClass() *value.Class {
	return value.InitDefinitionNodeClass
}

func (n *InitDefinitionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::InitDefinitionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  doc_comment: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	buff.WriteString(",\n  throw_type: ")
	if n.ThrowType == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.ThrowType.Inspect(), 1)
	}

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

func (n *InitDefinitionNode) Error() string {
	return n.Inspect()
}

// Create a constructor definition node eg. `init then 'hello world'`
func NewInitDefinitionNode(loc *position.Location, params []ParameterNode, throwType TypeNode, body []StatementNode) *InitDefinitionNode {
	return &InitDefinitionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Parameters:    params,
		ThrowType:     throwType,
		Body:          body,
	}
}

// Represents a method signature definition eg. `sig to_string(val: Int): String`
type MethodSignatureDefinitionNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Name           IdentifierNode
	TypeParameters []TypeParameterNode
	Parameters     []ParameterNode // formal parameters
	ReturnType     TypeNode
	ThrowType      TypeNode
}

func (n *MethodSignatureDefinitionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	name := n.Name.splice(loc, args, unquote).(IdentifierNode)
	typeParams := SpliceSlice(n.TypeParameters, loc, args, unquote)
	params := SpliceSlice(n.Parameters, loc, args, unquote)

	var returnType TypeNode
	if n.ReturnType != nil {
		returnType = n.ReturnType.splice(loc, args, unquote).(TypeNode)
	}

	var throwType TypeNode
	if n.ThrowType != nil {
		throwType = n.ThrowType.splice(loc, args, unquote).(TypeNode)
	}

	return &MethodSignatureDefinitionNode{
		TypedNodeBase:          TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		DocCommentableNodeBase: n.DocCommentableNodeBase,
		Name:                   name,
		TypeParameters:         typeParams,
		Parameters:             params,
		ReturnType:             returnType,
		ThrowType:              throwType,
	}
}

func (n *MethodSignatureDefinitionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::MethodSignatureDefinitionNode", env)
}

func (n *MethodSignatureDefinitionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Name.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	for _, param := range n.TypeParameters {
		if param.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
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

func (n *MethodSignatureDefinitionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*MethodSignatureDefinitionNode)
	if !ok {
		return false
	}

	if !n.Name.Equal(value.Ref(o.Name)) ||
		!n.loc.Equal(o.loc) ||
		n.comment != o.comment {
		return false
	}

	if len(n.TypeParameters) != len(o.TypeParameters) ||
		len(n.Parameters) != len(o.Parameters) {
		return false
	}

	for i, param := range n.TypeParameters {
		if !param.Equal(value.Ref(o.TypeParameters[i])) {
			return false
		}
	}

	for i, param := range n.Parameters {
		if !param.Equal(value.Ref(o.Parameters[i])) {
			return false
		}
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

	return true
}

func (n *MethodSignatureDefinitionNode) String() string {
	var buff strings.Builder

	doc := n.DocComment()
	if len(doc) > 0 {
		buff.WriteString("##[\n")
		indent.IndentString(&buff, doc, 1)
		buff.WriteString("\n]##\n")
	}

	buff.WriteString("sig ")
	buff.WriteString(n.Name.String())

	if len(n.TypeParameters) > 0 {
		buff.WriteRune('[')
		for i, param := range n.TypeParameters {
			if i > 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(param.String())
		}
		buff.WriteRune(']')
	}

	if len(n.Parameters) > 0 {
		buff.WriteRune('(')
		for i, param := range n.Parameters {
			if i > 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(param.String())
		}
		buff.WriteRune(')')
	}

	if n.ReturnType != nil {
		buff.WriteString(": ")
		buff.WriteString(n.ReturnType.String())
	}

	if n.ThrowType != nil {
		buff.WriteString(" ! ")
		buff.WriteString(n.ThrowType.String())
	}

	return buff.String()
}

func (*MethodSignatureDefinitionNode) IsStatic() bool {
	return false
}

func (*MethodSignatureDefinitionNode) Class() *value.Class {
	return value.MethodSignatureDefinitionNodeClass
}

func (*MethodSignatureDefinitionNode) DirectClass() *value.Class {
	return value.MethodSignatureDefinitionNodeClass
}

func (n *MethodSignatureDefinitionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::MethodSignatureDefinitionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  name: ")
	indent.IndentStringFromSecondLine(&buff, n.Name.Inspect(), 1)

	buff.WriteString(",\n  doc_comment: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	buff.WriteString(",\n  return_type: ")
	if n.ReturnType == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.ReturnType.Inspect(), 1)
	}

	buff.WriteString(",\n  throw_type: ")
	if n.ThrowType == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.ThrowType.Inspect(), 1)
	}

	buff.WriteString(",\n  type_parameters: %[\n")
	for i, element := range n.TypeParameters {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  parameters: %[\n")
	for i, element := range n.Parameters {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *MethodSignatureDefinitionNode) Error() string {
	return n.Inspect()
}

// Create a method signature node eg. `sig to_string(val: Int): String`
func NewMethodSignatureDefinitionNode(loc *position.Location, docComment string, name IdentifierNode, typeParams []TypeParameterNode, params []ParameterNode, returnType, throwType TypeNode) *MethodSignatureDefinitionNode {
	return &MethodSignatureDefinitionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Name:           name,
		TypeParameters: typeParams,
		Parameters:     params,
		ReturnType:     returnType,
		ThrowType:      throwType,
	}
}

// A single alias entry eg. `new_name old_name`
type AliasDeclarationEntry struct {
	NodeBase
	NewName IdentifierNode
	OldName IdentifierNode
}

func (n *AliasDeclarationEntry) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &AliasDeclarationEntry{
		NodeBase: n.NodeBase,
		NewName:  n.NewName.splice(loc, args, unquote).(IdentifierNode),
		OldName:  n.OldName.splice(loc, args, unquote).(IdentifierNode),
	}
}

func (n *AliasDeclarationEntry) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::AliasDeclarationEntry", env)
}

func (n *AliasDeclarationEntry) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.OldName.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	if n.NewName.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (*AliasDeclarationEntry) IsStatic() bool {
	return false
}

func (*AliasDeclarationEntry) Class() *value.Class {
	return value.MethodSignatureDefinitionNodeClass
}

func (*AliasDeclarationEntry) DirectClass() *value.Class {
	return value.MethodSignatureDefinitionNodeClass
}

func (n *AliasDeclarationEntry) Inspect() string {
	var buff strings.Builder

	buff.WriteString("Std::Elk::AST::AliasDeclarationEntry{\n")

	fmt.Fprintf(&buff, "  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  new_name: %s")
	indent.IndentStringFromSecondLine(&buff, n.NewName.Inspect(), 1)

	buff.WriteString(",\n  old_name: %s")
	indent.IndentStringFromSecondLine(&buff, n.OldName.Inspect(), 1)

	buff.WriteRune('}')

	return buff.String()
}

func (n *AliasDeclarationEntry) Equal(other value.Value) bool {
	if !other.IsReference() {
		return false
	}

	o, ok := other.AsReference().(*AliasDeclarationEntry)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.NewName.Equal(value.Ref(o.NewName)) &&
		n.OldName.Equal(value.Ref(o.OldName))
}

func (n *AliasDeclarationEntry) String() string {
	var buff strings.Builder

	buff.WriteString(n.NewName.String())
	buff.WriteRune(' ')
	buff.WriteString(n.OldName.String())

	return buff.String()
}

func (n *AliasDeclarationEntry) Error() string {
	return n.Inspect()
}

// Create an alias alias entry eg. `new_name old_name`
func NewAliasDeclarationEntry(loc *position.Location, newName, oldName IdentifierNode) *AliasDeclarationEntry {
	return &AliasDeclarationEntry{
		NodeBase: NodeBase{loc: loc},
		NewName:  newName,
		OldName:  oldName,
	}
}

// Represents a new alias declaration eg. `alias push append, add plus`
type AliasDeclarationNode struct {
	TypedNodeBase
	Entries []*AliasDeclarationEntry
}

func (n *AliasDeclarationNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &AliasDeclarationNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Entries:       SpliceSlice(n.Entries, loc, args, unquote),
	}
}

func (n *AliasDeclarationNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::AliasDeclarationNode", env)
}

func (n *AliasDeclarationNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, entry := range n.Entries {
		if entry.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *AliasDeclarationNode) String() string {
	var buff strings.Builder

	buff.WriteString("alias ")
	for i, entry := range n.Entries {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(entry.String())
	}

	return buff.String()
}

func (n *AliasDeclarationNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*AliasDeclarationNode)
	if !ok {
		return false
	}

	if len(n.Entries) != len(o.Entries) {
		return false
	}

	if !n.loc.Equal(o.loc) {
		return false
	}

	for i, entry := range n.Entries {
		if !entry.Equal(value.Ref(o.Entries[i])) {
			return false
		}
	}

	return true
}

func (*AliasDeclarationNode) IsStatic() bool {
	return false
}

func (*AliasDeclarationNode) Class() *value.Class {
	return value.AliasDeclarationNodeClass
}

func (*AliasDeclarationNode) DirectClass() *value.Class {
	return value.AliasDeclarationNodeClass
}

func (n *AliasDeclarationNode) Inspect() string {
	var buff strings.Builder

	buff.WriteString("Std::Elk::AST::AliasDeclarationNode{\n")

	fmt.Fprintf(&buff, "  location: %s,\n", (*value.Location)(n.loc).Inspect())
	buff.WriteString("  entries: %[\n")
	for i, element := range n.Entries {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *AliasDeclarationNode) Error() string {
	return n.Inspect()
}

// Create an alias declaration node eg. `alias push append, add plus`
func NewAliasDeclarationNode(loc *position.Location, entries []*AliasDeclarationEntry) *AliasDeclarationNode {
	return &AliasDeclarationNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Entries:       entries,
	}
}

// Represents a new getter declaration eg. `getter foo: String`
type GetterDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Entries []ParameterNode
}

func (n *GetterDeclarationNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &GetterDeclarationNode{
		TypedNodeBase:          TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		DocCommentableNodeBase: n.DocCommentableNodeBase,
		Entries:                SpliceSlice(n.Entries, loc, args, unquote),
	}
}

func (n *GetterDeclarationNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::GetterDeclarationNode", env)
}

func (n *GetterDeclarationNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, entry := range n.Entries {
		if entry.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

// Equal checks if this node equals the other node.
func (n *GetterDeclarationNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*GetterDeclarationNode)
	if !ok {
		return false
	}

	if !n.loc.Equal(o.loc) ||
		n.comment != o.comment {
		return false
	}

	if len(n.Entries) != len(o.Entries) {
		return false
	}

	for i, entry := range n.Entries {
		if !entry.Equal(value.Ref(o.Entries[i])) {
			return false
		}
	}

	return true
}

// String returns the string representation of this node.
func (n *GetterDeclarationNode) String() string {
	var buff strings.Builder

	doc := n.DocComment()
	if len(doc) > 0 {
		buff.WriteString("##[\n")
		indent.IndentString(&buff, doc, 1)
		buff.WriteString("\n]##\n")
	}

	buff.WriteString("getter ")

	for i, entry := range n.Entries {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(entry.String())
	}

	return buff.String()
}

func (*GetterDeclarationNode) IsStatic() bool {
	return false
}

func (*GetterDeclarationNode) Class() *value.Class {
	return value.GetterDeclarationNodeClass
}

func (*GetterDeclarationNode) DirectClass() *value.Class {
	return value.GetterDeclarationNodeClass
}

func (n *GetterDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::GetterDeclarationNode{\n  loc: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  doc_comment: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	buff.WriteString(",\n  entries: %[\n")
	for i, element := range n.Entries {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *GetterDeclarationNode) Error() string {
	return n.Inspect()
}

// Create a getter declaration node eg. `getter foo: String`
func NewGetterDeclarationNode(loc *position.Location, docComment string, entries []ParameterNode) *GetterDeclarationNode {
	return &GetterDeclarationNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Entries: entries,
	}
}

// Represents a new setter declaration eg. `setter foo: String`
type SetterDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Entries []ParameterNode
}

func (n *SetterDeclarationNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &SetterDeclarationNode{
		TypedNodeBase:          TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		DocCommentableNodeBase: n.DocCommentableNodeBase,
		Entries:                SpliceSlice(n.Entries, loc, args, unquote),
	}
}

func (n *SetterDeclarationNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::SetterDeclarationNode", env)
}

func (n *SetterDeclarationNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, entry := range n.Entries {
		if entry.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *SetterDeclarationNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SetterDeclarationNode)
	if !ok {
		return false
	}

	if len(n.Entries) != len(o.Entries) {
		return false
	}

	for i, entry := range n.Entries {
		if !entry.Equal(value.Ref(o.Entries[i])) {
			return false
		}
	}

	return n.loc.Equal(o.loc) &&
		n.comment == o.comment
}

func (n *SetterDeclarationNode) String() string {
	var buff strings.Builder

	doc := n.DocComment()
	if len(doc) > 0 {
		buff.WriteString("##[\n")
		indent.IndentString(&buff, doc, 1)
		buff.WriteString("\n]##\n")
	}

	buff.WriteString("setter ")
	for i, entry := range n.Entries {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(entry.String())
	}

	return buff.String()
}

func (*SetterDeclarationNode) IsStatic() bool {
	return false
}

func (*SetterDeclarationNode) Class() *value.Class {
	return value.SetterDeclarationNodeClass
}

func (*SetterDeclarationNode) DirectClass() *value.Class {
	return value.SetterDeclarationNodeClass
}

func (n *SetterDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SetterDeclarationNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  doc_comment: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	buff.WriteString(",\n  entries: %[\n")
	for i, element := range n.Entries {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *SetterDeclarationNode) Error() string {
	return n.Inspect()
}

// Create a setter declaration node eg. `setter foo: String`
func NewSetterDeclarationNode(loc *position.Location, docComment string, entries []ParameterNode) *SetterDeclarationNode {
	return &SetterDeclarationNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Entries: entries,
	}
}

// Represents a new setter declaration eg. `attr foo: String`
type AttrDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Entries []ParameterNode
}

func (n *AttrDeclarationNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &AttrDeclarationNode{
		TypedNodeBase:          TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		DocCommentableNodeBase: n.DocCommentableNodeBase,
		Entries:                SpliceSlice(n.Entries, loc, args, unquote),
	}
}

func (n *AttrDeclarationNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::AttrDeclarationNode", env)
}

func (n *AttrDeclarationNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, entry := range n.Entries {
		if entry.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (*AttrDeclarationNode) IsStatic() bool {
	return false
}

func (*AttrDeclarationNode) Class() *value.Class {
	return value.AttrDeclarationNodeClass
}

func (*AttrDeclarationNode) DirectClass() *value.Class {
	return value.AttrDeclarationNodeClass
}

func (n *AttrDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::AttrDeclarationNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  doc_comment: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	buff.WriteString(",\n  entries: %[\n")
	for i, element := range n.Entries {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *AttrDeclarationNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*AttrDeclarationNode)
	if !ok {
		return false
	}

	if !n.loc.Equal(o.loc) {
		return false
	}

	if n.DocComment() != o.DocComment() {
		return false
	}

	if len(n.Entries) != len(o.Entries) {
		return false
	}

	for i, entry := range n.Entries {
		if !entry.Equal(value.Ref(o.Entries[i])) {
			return false
		}
	}

	return true
}

func (n *AttrDeclarationNode) String() string {
	var buff strings.Builder

	doc := n.DocComment()
	if len(doc) > 0 {
		buff.WriteString("##[\n")
		indent.IndentString(&buff, doc, 1)
		buff.WriteString("\n]##\n")
	}
	buff.WriteString("attr ")

	for i, entry := range n.Entries {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(entry.String())
	}

	return buff.String()
}

func (n *AttrDeclarationNode) Error() string {
	return n.Inspect()
}

// Create an attribute declaration node eg. `attr foo: String`
func NewAttrDeclarationNode(loc *position.Location, docComment string, entries []ParameterNode) *AttrDeclarationNode {
	return &AttrDeclarationNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Entries: entries,
	}
}
