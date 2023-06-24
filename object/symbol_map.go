package object

import (
	"fmt"
	"strings"
)

var SymbolMapClass *Class // ::Std::SymbolMap

// Map with symbol keys.
type SymbolMap struct {
	SimpleSymbolMap
	frozen bool
}

func (s *SymbolMap) IsFrozen() bool {
	return s.frozen
}

func (s *SymbolMap) SetFrozen() {
	s.frozen = true
}

func (s *SymbolMap) Class() *Class {
	return SymbolMapClass
}

func (s *SymbolMap) InstanceVariables() SimpleSymbolMap {
	return nil
}

// Simple map with symbol keys.
type SimpleSymbolMap map[SymbolId]Value

func (s SimpleSymbolMap) Inspect() string {
	var buff strings.Builder
	buff.WriteString("{ ")
	firstIteration := true

	for key, val := range s {
		if !firstIteration {
			buff.WriteString(", ")
		}
		symbol := SymbolTable.GetId(key)
		buff.WriteString(fmt.Sprintf("%s: %s", symbol.Name, val.Inspect()))
		firstIteration = false
	}

	buff.WriteString(" }")
	return buff.String()
}

// Get a value stored using the given key.
func (s SimpleSymbolMap) Get(key *Symbol) Value {
	return s[key.Id]
}

// Set the passed value under the given key.
func (s SimpleSymbolMap) Set(key *Symbol, val Value) {
	s[key.Id] = val
}

// Get a value stored using the given key.
func (s SimpleSymbolMap) GetString(key string) Value {
	symbol := SymbolTable.Get(key)
	if symbol == nil {
		return nil
	}

	return s[symbol.Id]
}

// Set the passed value under the given key.
func (s SimpleSymbolMap) SetString(key string, val Value) {
	s[SymbolTable.Add(key).Id] = val
}

// Get a value stored using the given ID.
func (s SimpleSymbolMap) GetId(id SymbolId) Value {
	return s[id]
}

// Set the passed value under the given ID.
func (s SimpleSymbolMap) SetId(id SymbolId, val Value) {
	s[id] = val
}

func initSymbolMap() {
	SymbolMapClass = NewClass(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstant("SymbolMap", SymbolMapClass)
}
