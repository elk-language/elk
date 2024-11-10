package types

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

type ParameterKind uint8

func (p ParameterKind) String() string {
	return parameterKindNames[p]
}

var parameterKindNames = []string{
	NormalParameterKind:         "NormalParameterKind",
	DefaultValueParameterKind:   "DefaultValueParameterKind",
	PositionalRestParameterKind: "PositionalRestParameterKind",
	NamedRestParameterKind:      "NamedRestParameterKind",
}

const (
	NormalParameterKind ParameterKind = iota
	DefaultValueParameterKind
	PositionalRestParameterKind
	NamedRestParameterKind
)

type Parameter struct {
	Name             value.Symbol
	Type             Type
	Kind             ParameterKind
	InstanceVariable bool
}

func (p *Parameter) Copy() *Parameter {
	return &Parameter{
		Name:             p.Name,
		Type:             p.Type,
		Kind:             p.Kind,
		InstanceVariable: p.InstanceVariable,
	}
}

func NewParameter(name value.Symbol, typ Type, kind ParameterKind, instanceVariable bool) *Parameter {
	return &Parameter{
		Name:             name,
		Type:             typ,
		Kind:             kind,
		InstanceVariable: instanceVariable,
	}
}

func (p *Parameter) NameWithKind() string {
	name := p.Name.String()
	if p.InstanceVariable {
		name = "@" + name
	}
	switch p.Kind {
	case NormalParameterKind:
		return name
	case DefaultValueParameterKind:
		return fmt.Sprintf("%s = x", name)
	case PositionalRestParameterKind:
		return fmt.Sprintf("*%s", name)
	case NamedRestParameterKind:
		return fmt.Sprintf("**%s", name)
	default:
		panic("invalid parameter kind")
	}
}

func (p *Parameter) inspect() string {
	buffer := new(strings.Builder)
	switch p.Kind {
	case PositionalRestParameterKind:
		buffer.WriteRune('*')
	case NamedRestParameterKind:
		buffer.WriteString("**")
	}
	buffer.WriteString(p.Name.String())

	switch p.Kind {
	case DefaultValueParameterKind:
		buffer.WriteRune('?')
	}

	buffer.WriteString(": ")
	buffer.WriteString(Inspect(p.Type))
	return buffer.String()
}

func (p *Parameter) IsPositionalRest() bool {
	return p.Kind == PositionalRestParameterKind
}

func (p *Parameter) IsNamedRest() bool {
	return p.Kind == NamedRestParameterKind
}

func (p *Parameter) HasDefaultValue() bool {
	return p.Kind == DefaultValueParameterKind
}

func (p *Parameter) ToNonLiteral(env *GlobalEnvironment) Type {
	return p
}

func (*Parameter) IsLiteral() bool {
	return false
}

func (p *Parameter) IsOptional() bool {
	switch p.Kind {
	case DefaultValueParameterKind:
		return true
	default:
		return false
	}
}

type AstNode interface {
	position.SpanInterface
}

type Method struct {
	DocComment               string
	FullName                 string
	Name                     value.Symbol
	OptionalParamCount       int
	PostParamCount           int
	abstract                 bool
	sealed                   bool
	native                   bool
	compiled                 bool
	HasNamedRestParam        bool
	InstanceVariablesChecked bool
	// used in using expression placeholders
	IsPlaceholder bool
	Checked       bool
	Replaced      bool

	Params         []*Parameter
	TypeParameters []*TypeParameter
	ReturnType     Type
	ThrowType      Type
	DefinedUnder   Namespace
	Bytecode       *vm.BytecodeFunction
	location       *position.Location
	// used to detect methods that circularly reference constants
	UsedInConstants              ds.Set[value.Symbol] // set of constants in which this method is called
	UsedConstants                ds.Set[value.Symbol] // set of constants references in this method's body
	InitialisedInstanceVariables ds.Set[value.Symbol] // a set of names of instance variables that have been initialised, used when checking constructors
	CalledMethods                []*Method            // list of methods called in this method's body
	Node                         AstNode
}

func NewMethodPlaceholder(fullName string, name value.Symbol, definedUnder Namespace, location *position.Location) *Method {
	return &Method{
		FullName:      fullName,
		Name:          name,
		DefinedUnder:  definedUnder,
		location:      location,
		IsPlaceholder: true,
	}
}

func (m *Method) Copy() *Method {
	return &Method{
		FullName:                     m.FullName,
		Name:                         m.Name,
		DocComment:                   m.DocComment,
		Params:                       m.Params,
		OptionalParamCount:           m.OptionalParamCount,
		PostParamCount:               m.PostParamCount,
		abstract:                     m.abstract,
		sealed:                       m.sealed,
		native:                       m.native,
		TypeParameters:               m.TypeParameters,
		IsPlaceholder:                m.IsPlaceholder,
		HasNamedRestParam:            m.HasNamedRestParam,
		ReturnType:                   m.ReturnType,
		ThrowType:                    m.ThrowType,
		DefinedUnder:                 m.DefinedUnder,
		Bytecode:                     m.Bytecode,
		location:                     m.location,
		UsedInConstants:              m.UsedInConstants,
		UsedConstants:                m.UsedConstants,
		CalledMethods:                m.CalledMethods,
		InitialisedInstanceVariables: m.InitialisedInstanceVariables,
		Node:                         m.Node,
	}
}

func (m *Method) Location() *position.Location {
	return m.location
}

func (m *Method) SetLocation(location *position.Location) {
	m.location = location
}

func (m *Method) IsInit() bool {
	return m.Name == symbol.S_init
}

func (m *Method) IsGeneric() bool {
	return len(m.TypeParameters) > 0
}

func (m *Method) IsAbstract() bool {
	return m.abstract
}

func (m *Method) SetAbstract(abstract bool) *Method {
	m.abstract = abstract
	return m
}

func (m *Method) IsCompilable() bool {
	return m.Bytecode != nil && !m.compiled
}

func (m *Method) IsCompiled() bool {
	return m.compiled
}

func (m *Method) SetCompiled(compiled bool) *Method {
	m.compiled = compiled
	return m
}

func (m *Method) IsSealed() bool {
	return m.sealed
}

func (m *Method) SetSealed(sealed bool) *Method {
	m.sealed = sealed
	return m
}

func (m *Method) IsNative() bool {
	return m.native
}

func (m *Method) SetNative(native bool) *Method {
	m.native = native
	return m
}

func NewMethod(docComment string, abstract, sealed, native bool, name value.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType Type, throwType Type, definedUnder Namespace) *Method {
	var optParamCount int
	var hasNamedRestParam bool
	postParamCount := -1
	for _, param := range params {
		switch param.Kind {
		case NormalParameterKind:
			if postParamCount != -1 {
				postParamCount++
			}
		case DefaultValueParameterKind:
			optParamCount++
			if postParamCount != -1 {
				postParamCount++
			}
		case PositionalRestParameterKind:
			postParamCount++
		case NamedRestParameterKind:
			hasNamedRestParam = true
		}
	}
	m := &Method{
		abstract:           abstract,
		sealed:             sealed,
		native:             native,
		Name:               name,
		TypeParameters:     typeParams,
		DocComment:         docComment,
		Params:             params,
		ReturnType:         returnType,
		ThrowType:          throwType,
		DefinedUnder:       definedUnder,
		HasNamedRestParam:  hasNamedRestParam,
		OptionalParamCount: optParamCount,
		PostParamCount:     postParamCount,
		UsedInConstants:    make(ds.Set[value.Symbol]),
		UsedConstants:      make(ds.Set[value.Symbol]),
	}
	if name == symbol.S_init {
		m.InitialisedInstanceVariables = make(ds.Set[value.Symbol])
	}

	return m
}

func (m *Method) RequiredParamCount() int {
	requiredParamCount := len(m.Params) - m.OptionalParamCount
	if m.HasNamedRestParam {
		requiredParamCount--
	}
	if m.HasPositionalRestParam() {
		requiredParamCount--
	}
	return requiredParamCount
}

func (m *Method) ExpectedParamCountString() string {
	requiredParamCount := m.RequiredParamCount()
	if m.HasNamedRestParam || m.HasPositionalRestParam() {
		return fmt.Sprintf("%d...", requiredParamCount)
	}

	if requiredParamCount == len(m.Params) {
		return fmt.Sprintf("%d", requiredParamCount)
	}

	return fmt.Sprintf("%d...%d", requiredParamCount, len(m.Params))
}

func (m *Method) NamedRestParam() *Parameter {
	if m.HasNamedRestParam {
		return m.Params[len(m.Params)-1]
	}
	return nil
}

func (m *Method) HasPositionalRestParam() bool {
	return m.PostParamCount != -1
}

func (m *Method) PositionalRestParamIndex() int {
	if m.PostParamCount == -1 {
		return -1
	}

	index := len(m.Params) - 1 - m.PostParamCount
	if m.HasNamedRestParam {
		index--
	}
	return index
}

func (m *Method) PositionalRestParam() *Parameter {
	index := m.PositionalRestParamIndex()
	if index == -1 {
		return nil
	}
	return m.Params[index]
}

func (m *Method) inspect() string {
	return inspectMethod(m.DefinedUnder, m.Name)
}

func inspectMethod(namespace Namespace, methodName value.Symbol) string {
	switch scope := namespace.(type) {
	case *Class, *Mixin:
		return fmt.Sprintf("%s.:%s", scope.Name(), methodName)
	case *Interface:
		return fmt.Sprintf("%s.:%s", scope.Name(), methodName)
	case *Module:
		return fmt.Sprintf("%s::%s", scope.Name(), methodName)
	case *Closure:
		return "call"
	case *MixinWithWhere:
		return inspectMethod(scope.Namespace, methodName)
	default:
		panic(fmt.Sprintf("method with invalid DefinedUnder: %#v", namespace))
	}
}

func (m *Method) InspectSignature(showModifiers bool) string {
	buffer := new(strings.Builder)
	if showModifiers {
		if m.abstract {
			buffer.WriteString("abstract ")
		}
		if m.sealed {
			buffer.WriteString("sealed ")
		}
		if m.native {
			buffer.WriteString("native ")
		}
	}
	buffer.WriteString("def ")
	buffer.WriteString(m.Name.String())

	if len(m.TypeParameters) > 0 {
		buffer.WriteRune('[')
		firstIteration := true
		for _, param := range m.TypeParameters {
			if !firstIteration {
				buffer.WriteString(", ")
			} else {
				firstIteration = false
			}
			switch param.Variance {
			case COVARIANT:
				buffer.WriteRune('+')
			case CONTRAVARIANT:
				buffer.WriteRune('-')
			case BIVARIANT:
				buffer.WriteString("+-")
			}
			buffer.WriteString(param.Name.String())
			if !IsNever(param.LowerBound) {
				buffer.WriteString(" > ")
				buffer.WriteString(Inspect(param.LowerBound))
			}
			if !IsAny(param.UpperBound) {
				buffer.WriteString(" < ")
				buffer.WriteString(Inspect(param.UpperBound))
			}
		}
		buffer.WriteRune(']')
	}

	buffer.WriteRune('(')
	firstIteration := true
	for _, param := range m.Params {
		if !firstIteration {
			buffer.WriteString(", ")
		} else {
			firstIteration = false
		}
		if param.IsPositionalRest() {
			buffer.WriteRune('*')
		} else if param.IsNamedRest() {
			buffer.WriteString("**")
		}
		buffer.WriteString(param.Name.String())
		if param.HasDefaultValue() {
			buffer.WriteRune('?')
		}
		buffer.WriteString(": ")
		buffer.WriteString(Inspect(param.Type))
	}
	buffer.WriteRune(')')
	returnType := m.ReturnType
	if returnType == nil {
		returnType = Void{}
	}
	buffer.WriteString(": ")
	buffer.WriteString(Inspect(returnType))

	throwType := m.ThrowType
	if throwType != nil && !IsNever(throwType) {
		buffer.WriteString(" ! ")
		buffer.WriteString(Inspect(throwType))
	}

	return buffer.String()
}

func (m *Method) InspectSignatureWithColor(showModifiers bool) string {
	return lexer.Colorize(m.InspectSignature(showModifiers))
}

func (m *Method) ToNonLiteral(env *GlobalEnvironment) Type {
	return m
}

func (*Method) IsLiteral() bool {
	return false
}
