package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

// A simple getter method.
type GetterMethod struct {
	AttributeName value.Symbol
	frozen        bool
}

func (g *GetterMethod) Name() value.Symbol {
	return g.AttributeName
}

func (*GetterMethod) ParameterCount() int {
	return 0
}

func (*GetterMethod) Parameters() []value.Symbol {
	return nil
}

func (*GetterMethod) OptionalParameterCount() int {
	return 0
}

func (*GetterMethod) PostRestParameterCount() int {
	return 0
}

func (*GetterMethod) NamedRestParameter() bool {
	return false
}

func (*GetterMethod) Class() *value.Class {
	return nil
}

func (*GetterMethod) DirectClass() *value.Class {
	return nil
}

func (*GetterMethod) SingletonClass() *value.Class {
	return nil
}

func (g *GetterMethod) IsFrozen() bool {
	return g.frozen
}

func (g *GetterMethod) SetFrozen() {
	g.frozen = true
}

func (g *GetterMethod) Inspect() string {
	return fmt.Sprintf("Method{name: %s, type: :getter}", g.AttributeName.Inspect())
}

func (*GetterMethod) InstanceVariables() value.SymbolMap {
	return nil
}

func (g *GetterMethod) Call(self value.Value) (value.Value, value.Value) {
	iv := self.InstanceVariables()
	if iv == nil {
		return nil, value.NewCantAccessInstanceVariablesOnPrimitiveError(self.Inspect())
	}
	result, ok := iv.Get(g.AttributeName)
	if !ok {
		return value.Nil, nil
	}
	return result, nil
}

// Create a new getter method.
func NewGetterMethod(attrName value.Symbol, frozen bool) *GetterMethod {
	return &GetterMethod{
		AttributeName: attrName,
		frozen:        frozen,
	}
}

type GetterMethodOption func(*GetterMethod)

func GetterMethodWithFrozen() GetterMethodOption {
	return func(gm *GetterMethod) {
		gm.frozen = true
	}
}

func GetterMethodWithAttributeName(attrName value.Symbol) GetterMethodOption {
	return func(gm *GetterMethod) {
		gm.AttributeName = attrName
	}
}

func GetterMethodWithAttributeNameString(attrName string) GetterMethodOption {
	return func(gm *GetterMethod) {
		gm.AttributeName = value.ToSymbol(attrName)
	}
}

// Create a new getter method.
func NewGetterMethodWithOptions(opts ...GetterMethodOption) *GetterMethod {
	gm := &GetterMethod{}

	for _, opt := range opts {
		opt(gm)
	}

	return gm
}

// Utility method that creates a new getter method and
// attaches it as a method to the given method map.
func DefineGetter(
	methodMap value.MethodMap,
	name string,
	frozen bool,
) {
	symbolName := value.ToSymbol(name)
	getterMethod := NewGetterMethod(
		symbolName,
		frozen,
	)
	methodMap[symbolName] = getterMethod
}

// Utility method that creates a new getter method and
// attaches it as a method to the given method map.
func DefineGetterWithOptions(
	methodMap value.MethodMap,
	name string,
	opts ...GetterMethodOption,
) {
	getterMethod := &GetterMethod{}
	symbolName := value.ToSymbol(name)
	GetterMethodWithAttributeName(symbolName)(getterMethod)

	for _, opt := range opts {
		opt(getterMethod)
	}

	methodMap[symbolName] = getterMethod
}
