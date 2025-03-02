package ast

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/position"
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
	TypedNodeBaseWithLoc
	DocCommentableNodeBase
	Name           string
	TypeParameters []TypeParameterNode
	Parameters     []ParameterNode // formal parameters
	ReturnType     TypeNode
	ThrowType      TypeNode
	Body           []StatementNode // body of the method
	Flags          bitfield.BitField8
}

func (*MethodDefinitionNode) IsStatic() bool {
	return false
}

// Whether the method is a setter.
func (m *MethodDefinitionNode) IsSetter() bool {
	if len(m.Name) > 0 {
		firstChar, _ := utf8.DecodeRuneInString(m.Name)
		lastChar := m.Name[len(m.Name)-1]
		if (unicode.IsLetter(firstChar) || firstChar == '_') && lastChar == '=' {
			return true
		}
	}

	return false
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

const (
	METHOD_ABSTRACT_FLAG bitfield.BitFlag8 = 1 << iota
	METHOD_SEALED_FLAG
	METHOD_GENERATOR_FLAG
	METHOD_ASYNC_FLAG
)

func (*MethodDefinitionNode) Class() *value.Class {
	return value.MethodDefinitionNodeClass
}

func (*MethodDefinitionNode) DirectClass() *value.Class {
	return value.MethodDefinitionNodeClass
}

func (n *MethodDefinitionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::MethodDefinitionNode{\n  &: %p", n)

	buff.WriteString(",\n  doc_comment: ")
	indentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	fmt.Fprintf(&buff, ",\n  abstract: %t", n.IsAbstract())
	fmt.Fprintf(&buff, ",\n  sealed: %t", n.IsSealed())
	fmt.Fprintf(&buff, ",\n  generator: %t", n.IsGenerator())
	fmt.Fprintf(&buff, ",\n  async: %t", n.IsAsync())

	buff.WriteString(",\n  name: ")
	buff.WriteString(n.Name)

	buff.WriteString(",\n  return_type: ")
	indentStringFromSecondLine(&buff, n.ReturnType.Inspect(), 1)

	buff.WriteString(",\n  throw_type: ")
	indentStringFromSecondLine(&buff, n.ThrowType.Inspect(), 1)

	buff.WriteString(",\n  type_parameters: %%[\n")
	for i, element := range n.TypeParameters {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  parameters: %%[\n")
	for i, element := range n.Parameters {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  body: %%[\n")
	for i, element := range n.Body {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
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
	name string,
	typeParams []TypeParameterNode,
	params []ParameterNode,
	returnType,
	throwType TypeNode,
	body []StatementNode,
) *MethodDefinitionNode {
	return &MethodDefinitionNode{
		TypedNodeBaseWithLoc: TypedNodeBaseWithLoc{loc: loc},
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
	TypedNodeBaseWithLoc
	DocCommentableNodeBase
	Parameters []ParameterNode // formal parameters
	ThrowType  TypeNode
	Body       []StatementNode // body of the method
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

	fmt.Fprintf(&buff, "Std::AST::InitDefinitionNode{\n  &: %p", n)

	buff.WriteString(",\n  doc_comment: ")
	indentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	buff.WriteString(",\n  throw_type: ")
	indentStringFromSecondLine(&buff, n.ThrowType.Inspect(), 1)

	buff.WriteString(",\n  parameters: %%[\n")
	for i, element := range n.Parameters {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  body: %%[\n")
	for i, element := range n.Body {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
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
		TypedNodeBaseWithLoc: TypedNodeBaseWithLoc{loc: loc},
		Parameters:           params,
		ThrowType:            throwType,
		Body:                 body,
	}
}

// Represents a method signature definition eg. `sig to_string(val: Int): String`
type MethodSignatureDefinitionNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Name           string
	TypeParameters []TypeParameterNode
	Parameters     []ParameterNode // formal parameters
	ReturnType     TypeNode
	ThrowType      TypeNode
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

	fmt.Fprintf(&buff, "Std::AST::MethodSignatureDefinitionNode{\n  &: %p", n)

	buff.WriteString(",\n  doc_comment: ")
	indentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	buff.WriteString(",\n  return_type: ")
	indentStringFromSecondLine(&buff, n.ReturnType.Inspect(), 1)

	buff.WriteString(",\n  throw_type: ")
	indentStringFromSecondLine(&buff, n.ThrowType.Inspect(), 1)

	buff.WriteString(",\n  type_parameters: %%[\n")
	for i, element := range n.TypeParameters {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  parameters: %%[\n")
	for i, element := range n.Parameters {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *MethodSignatureDefinitionNode) Error() string {
	return n.Inspect()
}

// Create a method signature node eg. `sig to_string(val: Int): String`
func NewMethodSignatureDefinitionNode(span *position.Span, docComment, name string, typeParams []TypeParameterNode, params []ParameterNode, returnType, throwType TypeNode) *MethodSignatureDefinitionNode {
	return &MethodSignatureDefinitionNode{
		TypedNodeBase: TypedNodeBase{span: span},
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
	NewName string
	OldName string
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
	return fmt.Sprintf("Std::AST::AliasDeclarationEntry{&: %p, new_name: %s, old_name: %s}", n, n.NewName, n.OldName)
}

func (n *AliasDeclarationEntry) Error() string {
	return n.Inspect()
}

// Create an alias alias entry eg. `new_name old_name`
func NewAliasDeclarationEntry(span *position.Span, newName, oldName string) *AliasDeclarationEntry {
	return &AliasDeclarationEntry{
		NodeBase: NodeBase{span: span},
		NewName:  newName,
		OldName:  oldName,
	}
}

// Represents a new alias declaration eg. `alias push append, add plus`
type AliasDeclarationNode struct {
	TypedNodeBase
	Entries []*AliasDeclarationEntry
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

	fmt.Fprintf(&buff, "Std::AST::AliasDeclarationNode{\n  &: %p", n)

	buff.WriteString(",\n  entries: %%[\n")
	for i, element := range n.Entries {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *AliasDeclarationNode) Error() string {
	return n.Inspect()
}

// Create an alias declaration node eg. `alias push append, add plus`
func NewAliasDeclarationNode(span *position.Span, entries []*AliasDeclarationEntry) *AliasDeclarationNode {
	return &AliasDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Entries:       entries,
	}
}

// Represents a new getter declaration eg. `getter foo: String`
type GetterDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Entries []ParameterNode
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

	fmt.Fprintf(&buff, "Std::AST::GetterDeclarationNode{\n  &: %p", n)

	buff.WriteString(",\n  doc_comment: ")
	indentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	buff.WriteString(",\n  entries: %%[\n")
	for i, element := range n.Entries {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *GetterDeclarationNode) Error() string {
	return n.Inspect()
}

// Create a getter declaration node eg. `getter foo: String`
func NewGetterDeclarationNode(span *position.Span, docComment string, entries []ParameterNode) *GetterDeclarationNode {
	return &GetterDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
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

	fmt.Fprintf(&buff, "Std::AST::SetterDeclarationNode{\n  &: %p", n)

	buff.WriteString(",\n  doc_comment: ")
	indentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	buff.WriteString(",\n  entries: %%[\n")
	for i, element := range n.Entries {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *SetterDeclarationNode) Error() string {
	return n.Inspect()
}

// Create a setter declaration node eg. `setter foo: String`
func NewSetterDeclarationNode(span *position.Span, docComment string, entries []ParameterNode) *SetterDeclarationNode {
	return &SetterDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
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

	fmt.Fprintf(&buff, "Std::AST::AttrDeclarationNode{\n  &: %p", n)

	buff.WriteString(",\n  doc_comment: ")
	indentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	buff.WriteString(",\n  entries: %%[\n")
	for i, element := range n.Entries {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *AttrDeclarationNode) Error() string {
	return n.Inspect()
}

// Create an attribute declaration node eg. `attr foo: String`
func NewAttrDeclarationNode(span *position.Span, docComment string, entries []ParameterNode) *AttrDeclarationNode {
	return &AttrDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Entries: entries,
	}
}
