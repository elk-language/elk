package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

// A simple setter method.
type SetterMethod struct {
	AttributeName value.Symbol
	name          value.Symbol
	frozen        bool
}

func (s *SetterMethod) Name() value.Symbol {
	return s.name
}

func (*SetterMethod) ParameterCount() int {
	return 1
}

var setterParams = []value.Symbol{value.ToSymbol("val")}

func (*SetterMethod) Parameters() []value.Symbol {
	return setterParams
}

func (*SetterMethod) OptionalParameterCount() int {
	return 0
}

func (*SetterMethod) PostRestParameterCount() int {
	return 0
}

func (*SetterMethod) NamedRestParameter() bool {
	return false
}

func (*SetterMethod) Class() *value.Class {
	return nil
}

func (*SetterMethod) DirectClass() *value.Class {
	return nil
}

func (*SetterMethod) SingletonClass() *value.Class {
	return nil
}

func (s *SetterMethod) IsFrozen() bool {
	return s.frozen
}

func (s *SetterMethod) SetFrozen() {
	s.frozen = true
}

func (s *SetterMethod) Inspect() string {
	return fmt.Sprintf("Method{name: %s, type: :setter}", s.AttributeName.Inspect())
}

func (*SetterMethod) InstanceVariables() value.SymbolMap {
	return nil
}

// Create a new getter method.
func NewSetterMethod(attrName value.Symbol, frozen bool) *SetterMethod {
	return &SetterMethod{
		AttributeName: attrName,
		name:          value.ToSymbol(attrName.ToString() + "="),
		frozen:        frozen,
	}
}

type SetterMethodOption func(*SetterMethod)

func SetterMethodWithFrozen() SetterMethodOption {
	return func(sm *SetterMethod) {
		sm.frozen = true
	}
}

func SetterMethodWithAttributeName(attrName value.Symbol) SetterMethodOption {
	return func(sm *SetterMethod) {
		sm.AttributeName = attrName
		sm.name = value.ToSymbol(attrName.ToString() + "=")
	}
}

func SetterMethodWithAttributeNameString(attrName string) SetterMethodOption {
	return func(sm *SetterMethod) {
		sm.AttributeName = value.ToSymbol(attrName)
		sm.name = value.ToSymbol(attrName + "=")
	}
}

// Create a new getter method.
func NewSetterMethodWithOptions(opts ...SetterMethodOption) *SetterMethod {
	sm := &SetterMethod{}

	for _, opt := range opts {
		opt(sm)
	}

	return sm
}

// Utility method that creates a new setter method and
// attaches it as a method to the given method map.
func DefineSetter(
	methodMap value.MethodMap,
	attrName string,
	frozen bool,
) {
	symbolName := value.ToSymbol(attrName)
	setterMethod := NewSetterMethod(
		symbolName,
		frozen,
	)
	methodMap[setterMethod.name] = setterMethod
}

// Utility method that creates a new setter method and
// attaches it as a method to the given method map.
func DefineSetterWithOptions(
	methodMap value.MethodMap,
	attrName string,
	opts ...SetterMethodOption,
) {
	setterMethod := &SetterMethod{}
	symbolName := value.ToSymbol(attrName)
	SetterMethodWithAttributeName(symbolName)(setterMethod)

	for _, opt := range opts {
		opt(setterMethod)
	}

	methodMap[setterMethod.name] = setterMethod
}

// Utility method that creates a new setter and getter method and
// attaches them as methods to the given method map.
func DefineAccessor(
	methodMap value.MethodMap,
	attrName string,
	frozen bool,
) {
	DefineGetter(methodMap, attrName, frozen)
	DefineSetter(methodMap, attrName, frozen)
}
