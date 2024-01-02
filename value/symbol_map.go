package value

import (
	"fmt"
	"maps"
	"strings"
)

var SymbolMapClass *Class // ::Std::SymbolMap

// Simple map with symbol keys.
type SymbolMap map[Symbol]Value

func (SymbolMap) Class() *Class {
	return SymbolMapClass
}

func (SymbolMap) DirectClass() *Class {
	return SymbolMapClass
}

func (SymbolMap) SingletonClass() *Class {
	return nil
}

func (s SymbolMap) InstanceVariables() SymbolMap {
	return nil
}

func (s SymbolMap) Copy() Value {
	newMap := make(SymbolMap, len(s))
	maps.Copy(newMap, s)
	return newMap
}

func (s SymbolMap) Inspect() string {
	if len(s) == 0 {
		return "{}"
	}

	var buff strings.Builder
	buff.WriteString("{")
	firstIteration := true

	for symbol, val := range s {
		if !firstIteration {
			buff.WriteString(", ")
		}
		buff.WriteString(fmt.Sprintf("%s: %s", symbol.InspectContent(), val.Inspect()))
		firstIteration = false
	}

	buff.WriteString("}")
	return buff.String()
}

// Get a value stored using the given key.
func (s SymbolMap) Get(key Symbol) Value {
	return s[key]
}

func (s SymbolMap) Delete(key Symbol) {
	delete(s, key)
}

// Check if the given key exists.
func (s SymbolMap) Has(key Symbol) bool {
	_, ok := s[key]
	return ok
}

// Set the passed value under the given key.
func (s SymbolMap) Set(key Symbol, val Value) {
	s[key] = val
}

// Get a value stored using the given key.
func (s SymbolMap) GetString(key string) Value {
	symbol, ok := SymbolTable.Get(key)
	if !ok {
		return nil
	}

	return s[symbol]
}

// Check if the given key exists.
func (s SymbolMap) HasString(key string) bool {
	if s.GetString(key) == nil {
		return false
	}

	return true
}

func (s SymbolMap) DeleteString(key string) {
	symbol, ok := SymbolTable.Get(key)
	if !ok {
		return
	}

	delete(s, symbol)
}

// Set the passed value under the given key.
func (s SymbolMap) SetString(key string, val Value) {
	s[SymbolTable.Add(key)] = val
}

func initSymbolMap() {
	SymbolMapClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("SymbolMap", SymbolMapClass)
}
