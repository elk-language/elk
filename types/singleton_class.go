package types

import (
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value/symbol"
)

// Type that represents the singleton class of a mixin, class etc.
type SingletonClass struct {
	AttachedObject Ref[Namespace]
	Class
	id ID
}

func (c *SingletonClass) EqualAny(other any) bool {
	o, ok := other.(*SingletonClass)
	if !ok {
		return false
	}

	if c.id > 0 {
		return c.id == o.ID()
	}

	return c.name == o.name
}

func (c *SingletonClass) ToRef() Ref[*SingletonClass] {
	return Ref[*SingletonClass](c.id)
}

func (c *SingletonClass) SetParent(parent Namespace) {
	c.parent = ToRef(parent)
}

func (c *SingletonClass) RemoveTemporaryParents() {
	if _, ok := c.parent.Get().(*TemporaryParent); !ok {
		return
	}

	c.parent = ToRef(c.Superclass())
}

func NewSingletonClass(attached Namespace, parent Namespace) *SingletonClass {
	singleton := &SingletonClass{
		AttachedObject: ToRef(attached),
		Class: Class{
			parent:        ToRef(parent),
			NamespaceBase: MakeNamespaceBase("", "&"+attached.Name()),
		},
	}
	return singleton
}

func (s *SingletonClass) ToNonLiteral() Type {
	return s
}

func (s *SingletonClass) Copy() *SingletonClass {
	return &SingletonClass{
		AttachedObject: s.AttachedObject,
		Class:          s.Class,
	}
}

func (c *SingletonClass) DefineMethod(docComment string, flags bitfield.BitFlag16, name symbol.Symbol, typeParams []Ref[*TypeParameter], params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(docComment, flags, name, typeParams, params, returnType, throwType, c)
	c.SetMethod(name, method)
	return method
}
