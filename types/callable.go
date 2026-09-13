package types

import (
	"encoding/binary"
	"strings"

	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value/symbol"
)

type Callable struct {
	Body      Ref[*Method]
	IsClosure bool
	id        ID
}

func NewCallable(method *Method, isClosure bool) *Callable {
	t := &Callable{
		Body:      method.ToRef(),
		IsClosure: isClosure,
	}
	return t
}

func NewCallableWithMethod(docComment string, flags bitfield.BitFlag16, name symbol.Symbol, typeParams []Ref[*TypeParameter], params []Ref[*Parameter], returnType Type, throwType Type, isClosure bool) *Callable {
	callable := NewCallable(nil, isClosure)
	method := NewMethod(
		docComment,
		flags,
		name,
		typeParams,
		params,
		returnType,
		throwType,
		callable,
	)
	callable.Body = method.ToRef()
	return callable
}

func (c *Callable) HashUint64() uint64 {
	d := xxhash.New()
	d.WriteString("callable:")

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(c.Body))
	d.Write(b)

	var isClosureByte byte
	if c.IsClosure {
		isClosureByte = 1
	}
	d.Write([]byte{byte(isClosureByte)})

	return d.Sum64()
}

func (c *Callable) EqualAny(other any) bool {
	o, ok := other.(*Callable)
	if !ok {
		return false
	}

	if c.id > 0 {
		return c.id == o.ID()
	}

	return c.IsClosure == o.IsClosure &&
		c.Body == o.Body
}

func (c *Callable) ID() ID {
	return c.id
}

func (c *Callable) SetID(id ID) {
	c.id = id
}

func (c *Callable) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(c, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(c, parent)
	}
}

func (c *Callable) Copy() *Callable {
	return &Callable{
		Body: c.Body,
		id:   c.id,
	}
}

func IsCallable(namespace Namespace) bool {
	_, ok := namespace.(*Callable)
	return ok
}

func IsClosure(namespace Namespace) bool {
	c, ok := namespace.(*Callable)
	if !ok {
		return false
	}
	return c.IsClosure
}

func (c *Callable) Name() string {
	return "callable"
}

func (c *Callable) DocComment() string {
	return ""
}

func (c *Callable) SetDocComment(string) {
	panic("cannot set doc comment on callables")
}

func (c *Callable) AppendDocComment(string) {
	panic("cannot append doc comment on callables")
}

func (c *Callable) Parent() Namespace {
	return nil
}

func (c *Callable) SetParent(Namespace) {
	panic("cannot set parent of callable")
}

func (c *Callable) Singleton() *SingletonClass {
	return nil
}

func (c *Callable) SetSingleton(*SingletonClass) {
	panic("cannot set singleton class of callable")
}

func (c *Callable) IsAbstract() bool {
	return true
}

func (c *Callable) IsDefined() bool {
	return false
}

func (c *Callable) IsNative() bool {
	return false
}

func (c *Callable) SetDefined(bool) {
	panic("cannot set `compiled` in callable")
}

func (c *Callable) IsSealed() bool {
	return true
}

func (c *Callable) IsPrimitive() bool {
	return true
}

func (c *Callable) IsImmutable() bool {
	return true
}

func (c *Callable) Constants() ConstantMap {
	return nil
}

func (c *Callable) Constant(name symbol.Symbol) (Constant, bool) {
	return Constant{}, false
}

func (c *Callable) ConstantString(name string) (Constant, bool) {
	return Constant{}, false
}

func (c *Callable) DefineConstant(name symbol.Symbol, val Type) {
	panic("cannot define constants on callables")
}

func (c *Callable) DefineConstantWithFullName(name symbol.Symbol, fullName string, val Type) {
	panic("cannot define constants on callables")
}

func (c *Callable) Subtypes() ConstantMap {
	return nil
}

func (c *Callable) Subtype(name symbol.Symbol) (Constant, bool) {
	return Constant{}, false
}

func (c *Callable) SubtypeString(name string) (Constant, bool) {
	return Constant{}, false
}

func (c *Callable) MustSubtypeString(name string) Type {
	return nil
}

func (c *Callable) MustSubtype(name symbol.Symbol) Type {
	return nil
}

func (c *Callable) DefineSubtype(name symbol.Symbol, val Type) {
	panic("cannot define subtypes on callables")
}

func (c *Callable) DefineSubtypeWithFullName(name symbol.Symbol, fullName string, val Type) {
	panic("cannot define subtypes on callables")
}

func (c *Callable) Methods() MethodMap {
	if c.Body.IsZero() {
		return make(MethodMap)
	}

	m := make(MethodMap)
	m[symbol.L_call] = c.Body
	return m
}

func (c *Callable) Method(name symbol.Symbol) *Method {
	if name == symbol.L_call {
		return c.Body.Get()
	}
	return nil
}

func (c *Callable) MethodString(name string) *Method {
	if name == "call" {
		return c.Body.Get()
	}
	return nil
}

func (c *Callable) DefineMethod(docComment string, flags bitfield.BitFlag16, name symbol.Symbol, typeParams []Ref[*TypeParameter], params []Ref[*Parameter], returnType, throwType Type) *Method {
	panic("cannot define methods on callables")
}

func (c *Callable) SetMethod(name symbol.Symbol, method *Method) {
}

func (c *Callable) InstanceVariables() InstanceVariableMap {
	return nil
}

func (c *Callable) InstanceVariable(name symbol.Symbol) *InstanceVariable {
	return nil
}

func (c *Callable) InstanceVariableString(name string) *InstanceVariable {
	return nil
}

func (c *Callable) DefineInstanceVariable(name symbol.Symbol, ivar *InstanceVariable) {
	panic("cannot define instance variables on callables")
}

func (c *Callable) DefineClass(docComment string, primitive, abstract, sealed, noinit, immutable bool, name symbol.Symbol, parent Namespace) *Class {
	panic("cannot define classes on callables")
}

func (c *Callable) DefineModule(docComment string, name symbol.Symbol) *Module {
	panic("cannot define module on callables")
}

func (c *Callable) DefineMixin(docComment string, abstract bool, name symbol.Symbol) *Mixin {
	panic("cannot define mixins on callables")
}

func (c *Callable) DefineInterface(docComment string, name symbol.Symbol) *Interface {
	panic("cannot define interfaces on callables")
}

func (c *Callable) inspect() string {
	buffer := new(strings.Builder)
	body := c.Body.Get()
	if body.IsPure() {
		buffer.WriteString("pure ")
	}
	if c.IsClosure {
		buffer.WriteString("%|")
	} else {
		buffer.WriteRune('|')
	}
	firstIteration := true
	for _, paramRef := range body.Params {
		param := paramRef.Get()
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
		buffer.WriteString(Inspect(param.Type.Get()))
	}
	buffer.WriteRune('|')
	returnType := body.ReturnType.Get()
	if returnType == nil {
		returnType = Void{}
	}
	buffer.WriteString(": ")
	buffer.WriteString(Inspect(returnType))

	throwType := body.ThrowType.Get()
	if throwType != nil && !IsNever(throwType) {
		buffer.WriteString(" ! ")
		buffer.WriteString(Inspect(throwType))
	}

	return buffer.String()
}

func (c *Callable) ToNonLiteral() Type {
	return c
}

func (*Callable) IsLiteral() bool {
	return false
}

func (c *Callable) IsGeneric() bool {
	return false
}

func (c *Callable) TypeParameters() []Ref[*TypeParameter] {
	return nil
}

func (c *Callable) SetTypeParameters(t []Ref[*TypeParameter]) {
	panic("cannot set type parameters on a callable")
}
