package types

import (
	"strings"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type Closure struct {
	Body *Method
}

func NewClosure(method *Method) *Closure {
	return &Closure{
		Body: method,
	}
}

func NewClosureWithMethod(docComment string, abstract, sealed, native bool, name value.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType Type, throwType Type) *Closure {
	closure := NewClosure(nil)
	method := NewMethod(
		docComment,
		abstract,
		sealed,
		native,
		name,
		typeParams,
		params,
		returnType,
		throwType,
		closure,
	)
	closure.Body = method
	return closure
}

func IsClosure(namespace Namespace) bool {
	_, ok := namespace.(*Closure)
	return ok
}

func (c *Closure) Name() string {
	return c.inspect()
}

func (c *Closure) DocComment() string {
	return ""
}

func (c *Closure) SetDocComment(string) {
	panic("cannot set doc comment on closures")
}

func (c *Closure) AppendDocComment(string) {
	panic("cannot append doc comment on closures")
}

func (c *Closure) Parent() Namespace {
	return nil
}

func (c *Closure) SetParent(Namespace) {
	panic("cannot set parent of closure")
}

func (c *Closure) Singleton() *SingletonClass {
	return nil
}

func (c *Closure) IsAbstract() bool {
	return true
}

func (c *Closure) IsSealed() bool {
	return true
}

func (c *Closure) IsPrimitive() bool {
	return true
}

func (c *Closure) Constants() *TypeMap {
	return nil
}

func (c *Closure) Constant(name value.Symbol) Type {
	return nil
}

func (c *Closure) ConstantString(name string) Type {
	return nil
}

func (c *Closure) DefineConstant(name value.Symbol, val Type) {
	panic("cannot define constants on closures")
}

func (c *Closure) Subtypes() *TypeMap {
	return nil
}

func (c *Closure) Subtype(name value.Symbol) Type {
	return nil
}

func (c *Closure) SubtypeString(name string) Type {
	return nil
}

func (c *Closure) DefineSubtype(name value.Symbol, val Type) {
	panic("cannot define subtypes on closures")
}

func (c *Closure) Methods() *MethodMap {
	if c.Body == nil {
		return NewMethodMap()
	}
	m := NewMethodMap()
	m.Set(symbol.M_call, c.Body)
	return m
}

func (c *Closure) Method(name value.Symbol) *Method {
	if name == symbol.M_call {
		return c.Body
	}
	return nil
}

func (c *Closure) MethodString(name string) *Method {
	if name == "call" {
		return c.Body
	}
	return nil
}

func (c *Closure) DefineMethod(docComment string, abstract, sealed, native bool, name value.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType, throwType Type) *Method {
	panic("cannot define methods on closures")
}

func (c *Closure) SetMethod(name value.Symbol, method *Method) {
}

func (c *Closure) InstanceVariables() *TypeMap {
	return nil
}

func (c *Closure) InstanceVariable(name value.Symbol) Type {
	return nil
}

func (c *Closure) InstanceVariableString(name string) Type {
	return nil
}

func (c *Closure) DefineInstanceVariable(name value.Symbol, val Type) {
	panic("cannot define instance variables on closures")
}

func (c *Closure) DefineClass(docComment string, primitive, abstract, sealed bool, name value.Symbol, parent Namespace, env *GlobalEnvironment) *Class {
	panic("cannot define classes on closures")
}

func (c *Closure) DefineModule(docComment string, name value.Symbol, env *GlobalEnvironment) *Module {
	panic("cannot define module on closures")
}

func (c *Closure) DefineMixin(docComment string, abstract bool, name value.Symbol, env *GlobalEnvironment) *Mixin {
	panic("cannot define mixins on closures")
}

func (c *Closure) DefineInterface(docComment string, name value.Symbol, env *GlobalEnvironment) *Interface {
	panic("cannot define interfaces on closures")
}

func (c *Closure) inspect() string {
	buffer := new(strings.Builder)
	buffer.WriteRune('|')
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

func (c *Closure) ToNonLiteral(env *GlobalEnvironment) Type {
	return c
}

func (*Closure) IsLiteral() bool {
	return false
}

func (c *Closure) IsGeneric() bool {
	return false
}

func (c *Closure) TypeParameters() []*TypeParameter {
	return nil
}

func (c *Closure) SetTypeParameters(t []*TypeParameter) {
	panic("cannot set type parameters on a closure")
}
