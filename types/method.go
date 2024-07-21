package types

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
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

func (p *Parameter) IsPositionalRest() bool {
	return p.Kind == PositionalRestParameterKind
}

func (p *Parameter) IsNamedRest() bool {
	return p.Kind == NamedRestParameterKind
}

func (p *Parameter) HasDefaultValue() bool {
	return p.Kind == DefaultValueParameterKind
}

func (p *Parameter) IsOptional() bool {
	switch p.Kind {
	case DefaultValueParameterKind:
		return true
	default:
		return false
	}
}

type Method struct {
	Name               value.Symbol
	DocComment         string
	Params             []*Parameter
	OptionalParamCount int
	PostParamCount     int
	abstract           bool
	sealed             bool
	native             bool
	HasNamedRestParam  bool
	ReturnType         Type
	ThrowType          Type
	DefinedUnder       Namespace
	Bytecode           *vm.BytecodeFunction
	span               *position.Span
}

func (m *Method) Span() *position.Span {
	return m.span
}

func (m *Method) SetSpan(span *position.Span) {
	m.span = span
}

func (m *Method) IsAbstract() bool {
	return m.abstract
}

func (m *Method) SetAbstract(abstract bool) *Method {
	m.abstract = abstract
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

func NewMethod(docComment string, abstract, sealed, native bool, name value.Symbol, params []*Parameter, returnType Type, throwType Type, definedUnder Namespace) *Method {
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

	return &Method{
		abstract:           abstract,
		sealed:             sealed,
		native:             native,
		Name:               name,
		DocComment:         docComment,
		Params:             params,
		ReturnType:         returnType,
		ThrowType:          throwType,
		DefinedUnder:       definedUnder,
		HasNamedRestParam:  hasNamedRestParam,
		OptionalParamCount: optParamCount,
		PostParamCount:     postParamCount,
	}
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
	switch scope := m.DefinedUnder.(type) {
	case *Class, *Mixin, *Interface:
		return fmt.Sprintf("%s.:%s", scope.Name(), m.Name)
	case *Module:
		return fmt.Sprintf("%s::%s", scope.Name(), m.Name)
	default:
		panic(fmt.Sprintf("method without DefinedUnder: %#v", m.DefinedUnder))
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
	buffer.WriteRune('(')
	firstIteration := true
	for _, param := range m.Params {
		if !firstIteration {
			buffer.WriteString(", ")
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
