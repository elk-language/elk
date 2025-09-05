package types

import (
	"fmt"
	"iter"
	"strings"

	"github.com/elk-language/elk/bitfield"
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

func (p *Parameter) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *Parameter {
	newParam := p.Copy()
	newParam.Type = DeepCopyEnv(newParam.Type, oldEnv, newEnv)
	return newParam
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

const (
	METHOD_ABSTRACT_FLAG bitfield.BitFlag16 = 1 << iota
	METHOD_SEALED_FLAG
	METHOD_NATIVE_FLAG
	METHOD_COMPILED_FLAG
	METHOD_NAMED_REST_PARAM_FLAG
	METHOD_INSTANCE_VARIABLES_CHECKED_FLAG
	METHOD_ATTRIBUTE_FLAG
	METHOD_GENERATOR_FLAG
	METHOD_ASYNC_FLAG
	METHOD_OVERLOAD_FLAG
	METHOD_MACRO_FLAG
	// used in using expression placeholders
	METHOD_PLACEHOLDER_FLAG
	METHOD_CHECKED_FLAG
	METHOD_REPLACED_FLAG
)

type Method struct {
	DocComment         string
	FullName           string
	Name               value.Symbol
	OptionalParamCount int
	PostParamCount     int
	OverloadId         int
	Flags              bitfield.BitField16

	Params         []*Parameter
	TypeParameters []*TypeParameter
	Overloads      []*Method
	Base           *Method
	ReturnType     Type
	ThrowType      Type
	DefinedUnder   Namespace
	Body           value.Method
	location       *position.Location
	// used to detect methods that circularly reference constants
	UsedInConstants              ds.Set[value.Symbol] // set of constants in which this method is called
	UsedConstants                ds.Set[value.Symbol] // set of constants references in this method's body
	InitialisedInstanceVariables ds.Set[value.Symbol] // a set of names of instance variables that have been initialised, used when checking constructors
	CalledMethods                []*Method            // list of methods called in this method's body
	Node                         AstNode
}

func NewMethodPlaceholder(fullName string, name value.Symbol, definedUnder Namespace, location *position.Location) *Method {
	m := &Method{
		FullName:     fullName,
		Name:         name,
		DefinedUnder: definedUnder,
		location:     location,
	}
	m.SetPlaceholder(true)
	return m
}

func (m *Method) WithoutOverloads() *Method {
	if len(m.Overloads) == 0 {
		return m
	}

	copy := m.Copy()
	copy.Overloads = nil
	return copy
}

func (m *Method) Copy() *Method {
	return &Method{
		Base:                         m.Base,
		OverloadId:                   m.OverloadId,
		FullName:                     m.FullName,
		Name:                         m.Name,
		DocComment:                   m.DocComment,
		Params:                       m.Params,
		OptionalParamCount:           m.OptionalParamCount,
		PostParamCount:               m.PostParamCount,
		Flags:                        m.Flags,
		TypeParameters:               m.TypeParameters,
		ReturnType:                   m.ReturnType,
		ThrowType:                    m.ThrowType,
		DefinedUnder:                 m.DefinedUnder,
		Body:                         m.Body,
		location:                     m.location,
		UsedInConstants:              m.UsedInConstants,
		UsedConstants:                m.UsedConstants,
		CalledMethods:                m.CalledMethods,
		InitialisedInstanceVariables: m.InitialisedInstanceVariables,
		Node:                         m.Node,
		Overloads:                    m.Overloads,
	}
}

func (m *Method) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *Method {
	var newDefinedUnder Namespace

	if m.DefinedUnder != nil && !IsCallable(m.DefinedUnder) {
		newDefinedUnder = DeepCopyEnv(m.DefinedUnder, oldEnv, newEnv).(Namespace)
		if newMethod := newDefinedUnder.Method(m.Name); newMethod != nil {
			return newMethod
		}
	}

	newMethod := &Method{
		OverloadId:                   m.OverloadId,
		FullName:                     m.FullName,
		Name:                         m.Name,
		DocComment:                   m.DocComment,
		OptionalParamCount:           m.OptionalParamCount,
		PostParamCount:               m.PostParamCount,
		Flags:                        m.Flags,
		Body:                         m.Body,
		location:                     m.location,
		UsedInConstants:              m.UsedInConstants,
		UsedConstants:                m.UsedConstants,
		InitialisedInstanceVariables: m.InitialisedInstanceVariables,
		Node:                         m.Node,
	}
	if newDefinedUnder != nil {
		newMethod.DefinedUnder = newDefinedUnder
		newMethod.DefinedUnder.SetMethod(newMethod.Name, newMethod)
	}

	newOverloads := make([]*Method, len(m.Overloads))
	for i, overload := range m.Overloads {
		newOverloads[i] = overload.DeepCopyEnv(oldEnv, newEnv)
	}
	newMethod.Overloads = newOverloads

	if m.Base != nil {
		newMethod.Base = DeepCopyEnv(m.Base, oldEnv, newEnv).(*Method)
	}

	newMethod.ThrowType = DeepCopyEnv(m.ThrowType, oldEnv, newEnv)
	newMethod.ReturnType = DeepCopyEnv(m.ReturnType, oldEnv, newEnv)

	newParameters := make([]*Parameter, len(m.Params))
	for i, param := range m.Params {
		newParameters[i] = param.DeepCopyEnv(oldEnv, newEnv)
	}
	newMethod.Params = newParameters

	newTypeParameters := make([]*TypeParameter, len(m.TypeParameters))
	for i, typeParam := range m.TypeParameters {
		newTypeParameters[i] = typeParam.DeepCopyEnv(oldEnv, newEnv)
	}
	newMethod.TypeParameters = newTypeParameters

	newCalledMethods := make([]*Method, len(m.CalledMethods))
	for i, calledMethod := range m.CalledMethods {
		newCalledMethods[i] = calledMethod.DeepCopyEnv(oldEnv, newEnv)
	}
	newMethod.CalledMethods = newCalledMethods

	return newMethod
}

func (m *Method) CreateAlias(newName value.Symbol) *Method {
	alias := m.Copy()
	alias.Name = newName
	alias.Base = m
	return alias
}

func (m *Method) AllOverloads() iter.Seq[*Method] {
	return func(yield func(*Method) bool) {
		if !yield(m) {
			return
		}

		for _, overload := range m.Overloads {
			if !yield(overload) {
				return
			}
		}
	}
}

func (m *Method) ReverseOverloads() iter.Seq[*Method] {
	return func(yield func(*Method) bool) {
		for i := len(m.Overloads) - 1; i >= 0; i-- {
			overload := m.Overloads[i]
			if !yield(overload) {
				return
			}
		}

		if !yield(m) {
			return
		}
	}
}

func (m *Method) RegisterOverload(overload *Method) {
	m.Overloads = append(m.Overloads, overload)
	overload.OverloadId = len(m.Overloads)
	overload.Name = value.ToSymbol(fmt.Sprintf("%s@%d", overload.Name.String(), len(m.Overloads)))
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

func (m *Method) AreInstanceVariablesChecked() bool {
	return m.Flags.HasFlag(METHOD_INSTANCE_VARIABLES_CHECKED_FLAG)
}

func (m *Method) SetInstanceVariablesChecked(val bool) *Method {
	m.SetFlag(METHOD_INSTANCE_VARIABLES_CHECKED_FLAG, val)
	return m
}

func (m *Method) HasNamedRestParam() bool {
	return m.Flags.HasFlag(METHOD_NAMED_REST_PARAM_FLAG)
}

func (m *Method) SetNamedRestParam(val bool) *Method {
	m.SetFlag(METHOD_NAMED_REST_PARAM_FLAG, val)
	return m
}

func (m *Method) IsPlaceholder() bool {
	return m.Flags.HasFlag(METHOD_PLACEHOLDER_FLAG)
}

func (m *Method) SetPlaceholder(val bool) *Method {
	m.SetFlag(METHOD_PLACEHOLDER_FLAG, val)
	return m
}

func (m *Method) IsGenerator() bool {
	return m.Flags.HasFlag(METHOD_GENERATOR_FLAG)
}

func (m *Method) SetGenerator(val bool) *Method {
	m.SetFlag(METHOD_GENERATOR_FLAG, val)
	return m
}

func (m *Method) IsAsync() bool {
	return m.Flags.HasFlag(METHOD_ASYNC_FLAG)
}

func (m *Method) SetAsync(val bool) *Method {
	m.SetFlag(METHOD_ASYNC_FLAG, val)
	return m
}

func (m *Method) IsOverload() bool {
	return m.Flags.HasFlag(METHOD_OVERLOAD_FLAG)
}

func (m *Method) SetOverload(val bool) *Method {
	m.SetFlag(METHOD_OVERLOAD_FLAG, val)
	return m
}

func (m *Method) IsRegisteredOverload() bool {
	return m.OverloadId != 0
}

func (m *Method) IsMacro() bool {
	return m.Flags.HasFlag(METHOD_MACRO_FLAG)
}

func (m *Method) SetMacro(val bool) *Method {
	m.SetFlag(METHOD_MACRO_FLAG, val)
	return m
}

func (m *Method) IsReplaced() bool {
	return m.Flags.HasFlag(METHOD_REPLACED_FLAG)
}

func (m *Method) SetReplaced(val bool) *Method {
	m.SetFlag(METHOD_REPLACED_FLAG, val)
	return m
}

func (m *Method) IsChecked() bool {
	return m.Flags.HasFlag(METHOD_CHECKED_FLAG)
}

func (m *Method) SetChecked(val bool) *Method {
	m.SetFlag(METHOD_CHECKED_FLAG, val)
	return m
}

func (m *Method) IsAbstract() bool {
	return m.Flags.HasFlag(METHOD_ABSTRACT_FLAG)
}

func (m *Method) SetAbstract(abstract bool) *Method {
	m.SetFlag(METHOD_ABSTRACT_FLAG, abstract)
	return m
}

func (m *Method) IsDefinable() bool {
	if m.IsCompiled() || m.IsMacro() {
		return false
	}

	if m.Base != nil {
		_, hasBytecode := m.Base.Body.(*vm.BytecodeFunction)
		return hasBytecode || m.Base.IsAttribute()
	}

	_, hasBytecode := m.Body.(*vm.BytecodeFunction)
	return hasBytecode || m.IsAttribute()
}

func (m *Method) IsSetter() bool {
	nameString := m.Name.String()
	if len(nameString) < 1 {
		return false
	}

	return nameString[len(nameString)-1] == '='
}

func (m *Method) IsCompilable() bool {
	return !m.IsAbstract() && !m.IsCompiled() && !m.IsAttribute()
}

func (m *Method) IsCompiled() bool {
	return m.Flags.HasFlag(METHOD_COMPILED_FLAG)
}

func (m *Method) SetCompiled(compiled bool) *Method {
	m.SetFlag(METHOD_COMPILED_FLAG, compiled)
	return m
}

func (m *Method) IsSealed() bool {
	return m.Flags.HasFlag(METHOD_SEALED_FLAG)
}

func (m *Method) SetSealed(sealed bool) *Method {
	m.SetFlag(METHOD_SEALED_FLAG, sealed)
	return m
}

func (m *Method) IsNative() bool {
	return m.Flags.HasFlag(METHOD_NATIVE_FLAG)
}

func (m *Method) SetNative(native bool) *Method {
	m.SetFlag(METHOD_NATIVE_FLAG, native)
	return m
}

func (m *Method) IsAttribute() bool {
	return m.Flags.HasFlag(METHOD_ATTRIBUTE_FLAG)
}

func (m *Method) SetAttribute(val bool) *Method {
	m.SetFlag(METHOD_ATTRIBUTE_FLAG, val)
	return m
}

func (m *Method) SetFlag(flag bitfield.BitFlag16, val bool) {
	if val {
		m.Flags.SetFlag(flag)
	} else {
		m.Flags.UnsetFlag(flag)
	}
}

func NewMethod(docComment string, flags bitfield.BitFlag16, name value.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType Type, throwType Type, definedUnder Namespace) *Method {
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
		Name:               name,
		TypeParameters:     typeParams,
		DocComment:         docComment,
		Params:             params,
		ReturnType:         returnType,
		ThrowType:          throwType,
		DefinedUnder:       definedUnder,
		OptionalParamCount: optParamCount,
		PostParamCount:     postParamCount,
		UsedInConstants:    make(ds.Set[value.Symbol]),
		UsedConstants:      make(ds.Set[value.Symbol]),
		Flags:              bitfield.BitField16FromBitFlag(flags),
	}
	if hasNamedRestParam {
		m.SetNamedRestParam(true)
	}
	if name == symbol.S_init {
		m.InitialisedInstanceVariables = make(ds.Set[value.Symbol])
	}

	return m
}

func (m *Method) RequiredParamCount() int {
	requiredParamCount := len(m.Params) - m.OptionalParamCount
	if m.HasNamedRestParam() {
		requiredParamCount--
	}
	if m.HasPositionalRestParam() {
		requiredParamCount--
	}
	return requiredParamCount
}

func (m *Method) ExpectedParamCountString() string {
	requiredParamCount := m.RequiredParamCount()
	if m.HasNamedRestParam() || m.HasPositionalRestParam() {
		return fmt.Sprintf("%d...", requiredParamCount)
	}

	if requiredParamCount == len(m.Params) {
		return fmt.Sprintf("%d", requiredParamCount)
	}

	return fmt.Sprintf("%d...%d", requiredParamCount, len(m.Params))
}

func (m *Method) NamedRestParam() *Parameter {
	if m.HasNamedRestParam() {
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
	if m.HasNamedRestParam() {
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
		return fmt.Sprintf("%s.:%s", scope.Name(), methodName.String())
	case *Interface:
		return fmt.Sprintf("%s.:%s", scope.Name(), methodName.String())
	case *Module:
		return fmt.Sprintf("%s::%s", scope.Name(), methodName.String())
	case *SingletonClass:
		return fmt.Sprintf("%s::%s", scope.AttachedObject.Name(), methodName.String())
	case *Callable:
		return "call"
	case *MixinWithWhere:
		return inspectMethod(scope.Namespace, methodName)
	default:
		panic(fmt.Sprintf("method with invalid DefinedUnder: %#v, name: %s", namespace, methodName.String()))
	}
}

func (m *Method) InspectSignature(showModifiers bool) string {
	buffer := new(strings.Builder)
	if showModifiers {
		if m.IsAbstract() {
			buffer.WriteString("abstract ")
		}
		if m.IsSealed() {
			buffer.WriteString("sealed ")
		}
		if m.IsNative() {
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
