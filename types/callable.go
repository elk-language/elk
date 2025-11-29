package types

import (
	"strings"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type Callable struct {
	Body      *Method
	IsClosure bool
}

func NewCallable(method *Method, isClosure bool) *Callable {
	return &Callable{
		Body:      method,
		IsClosure: isClosure,
	}
}

func NewCallableWithMethod(docComment string, flags bitfield.BitFlag16, name value.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType Type, throwType Type, isClosure bool) *Callable {
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
	callable.Body = method
	return callable
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
	}
}

func (c *Callable) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *Callable {
	newCallable := c.Copy()
	newCallable.Body = c.Body.DeepCopyEnv(oldEnv, newEnv)
	c.Body.DefinedUnder = newCallable

	return newCallable
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

func (c *Callable) Constants() ConstantMap {
	return nil
}

func (c *Callable) Constant(name value.Symbol) (Constant, bool) {
	return Constant{}, false
}

func (c *Callable) ConstantString(name string) (Constant, bool) {
	return Constant{}, false
}

func (c *Callable) DefineConstant(name value.Symbol, val Type) {
	panic("cannot define constants on callables")
}

func (c *Callable) DefineConstantWithFullName(name value.Symbol, fullName string, val Type) {
	panic("cannot define constants on callables")
}

func (c *Callable) Subtypes() ConstantMap {
	return nil
}

func (c *Callable) Subtype(name value.Symbol) (Constant, bool) {
	return Constant{}, false
}

func (c *Callable) SubtypeString(name string) (Constant, bool) {
	return Constant{}, false
}

func (c *Callable) MustSubtypeString(name string) Type {
	return nil
}

func (c *Callable) MustSubtype(name value.Symbol) Type {
	return nil
}

func (c *Callable) DefineSubtype(name value.Symbol, val Type) {
	panic("cannot define subtypes on callables")
}

func (c *Callable) DefineSubtypeWithFullName(name value.Symbol, fullName string, val Type) {
	panic("cannot define subtypes on callables")
}

func (c *Callable) Methods() MethodMap {
	if c.Body == nil {
		return make(MethodMap)
	}
	m := make(MethodMap)
	m[symbol.L_call] = c.Body
	return m
}

func (c *Callable) Method(name value.Symbol) *Method {
	if name == symbol.L_call {
		return c.Body
	}
	return nil
}

func (c *Callable) MethodString(name string) *Method {
	if name == "call" {
		return c.Body
	}
	return nil
}

func (c *Callable) DefineMethod(docComment string, flags bitfield.BitFlag16, name value.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType, throwType Type) *Method {
	panic("cannot define methods on callables")
}

func (c *Callable) SetMethod(name value.Symbol, method *Method) {
}

func (c *Callable) InstanceVariables() TypeMap {
	return nil
}

func (c *Callable) InstanceVariable(name value.Symbol) Type {
	return nil
}

func (c *Callable) InstanceVariableString(name string) Type {
	return nil
}

func (c *Callable) DefineInstanceVariable(name value.Symbol, val Type) {
	panic("cannot define instance variables on callables")
}

func (c *Callable) DefineClass(docComment string, primitive, abstract, sealed, noinit bool, name value.Symbol, parent Namespace, env *GlobalEnvironment) *Class {
	panic("cannot define classes on callables")
}

func (c *Callable) DefineModule(docComment string, name value.Symbol, env *GlobalEnvironment) *Module {
	panic("cannot define module on callables")
}

func (c *Callable) DefineMixin(docComment string, abstract bool, name value.Symbol, env *GlobalEnvironment) *Mixin {
	panic("cannot define mixins on callables")
}

func (c *Callable) DefineInterface(docComment string, name value.Symbol, env *GlobalEnvironment) *Interface {
	panic("cannot define interfaces on callables")
}

func (c *Callable) inspect() string {
	buffer := new(strings.Builder)
	if c.IsClosure {
		buffer.WriteString("%|")
	} else {
		buffer.WriteRune('|')
	}
	firstIteration := true
	for _, param := range c.Body.Params {
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
	buffer.WriteRune('|')
	returnType := c.Body.ReturnType
	if returnType == nil {
		returnType = Void{}
	}
	buffer.WriteString(": ")
	buffer.WriteString(Inspect(returnType))

	throwType := c.Body.ThrowType
	if throwType != nil && !IsNever(throwType) {
		buffer.WriteString(" ! ")
		buffer.WriteString(Inspect(throwType))
	}

	return buffer.String()
}

func (c *Callable) ToNonLiteral(env *GlobalEnvironment) Type {
	return c
}

func (*Callable) IsLiteral() bool {
	return false
}

func (c *Callable) IsGeneric() bool {
	return false
}

func (c *Callable) TypeParameters() []*TypeParameter {
	return nil
}

func (c *Callable) SetTypeParameters(t []*TypeParameter) {
	panic("cannot set type parameters on a callable")
}
