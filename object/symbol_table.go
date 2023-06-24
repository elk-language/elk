package object

import (
	"sync"

	"github.com/elk-language/elk/config"
)

// The number of preallocated slots for symbols
// in the symbol table at startup.
var SYMBOL_TABLE_INITIAL_SIZE int

func init() {
	val, ok := config.IntFromEnvVar("ELK_SYMBOL_TABLE_INITIAL_SIZE")
	if !ok {
		SYMBOL_TABLE_INITIAL_SIZE = 128
		return
	}

	SYMBOL_TABLE_INITIAL_SIZE = val
}

// Data structure that holds Elk Symbols.
type symbolTableStruct struct {
	nameTable map[string]*Symbol
	idTable   []*Symbol
	mutex     sync.RWMutex
}

type symbolTableOption func(*symbolTableStruct)

func symbolTableWithNameTable(nameTable map[string]*Symbol) symbolTableOption {
	return func(s *symbolTableStruct) {
		s.nameTable = nameTable
	}
}

func symbolTableWithIdTable(idTable []*Symbol) symbolTableOption {
	return func(s *symbolTableStruct) {
		s.idTable = idTable
	}
}

func newSymbolTable(opts ...symbolTableOption) *symbolTableStruct {
	s := &symbolTableStruct{
		nameTable: make(map[string]*Symbol, SYMBOL_TABLE_INITIAL_SIZE),
		idTable:   make([]*Symbol, 0, SYMBOL_TABLE_INITIAL_SIZE),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Check if a symbol with the given SymbolId exists.
func (s *symbolTableStruct) ExistsId(id SymbolId) bool {
	return id < SymbolId(len(s.idTable)) && id >= 0
}

// Check if a symbol with the given name exists.
func (s *symbolTableStruct) Exists(name string) bool {
	return s.Get(name) != nil
}

// Get the Symbol with the specified name.
// This function is thread-safe.
func (s *symbolTableStruct) Get(name string) *Symbol {
	s.mutex.RLock()
	val := s.nameTable[name]
	s.mutex.RUnlock()

	return val
}

// Get the Symbol with the SymbolId.
// This function is thread-safe.
func (s *symbolTableStruct) GetId(id SymbolId) *Symbol {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if id >= SymbolId(len(s.idTable)) || id < 0 {
		return nil
	}
	val := s.idTable[id]

	return val
}

// Add a new symbol with the specified name.
// This function is idempotent, if the Symbol already exists
// nothing happens and a pointer to it gets returned.
// This function is thread-safe.
func (s *symbolTableStruct) Add(name string) *Symbol {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	val := s.nameTable[name]
	if val != nil {
		return val
	}

	id := SymbolId(len(s.idTable))
	val = newSymbol(name, id)
	s.nameTable[name] = val
	s.idTable = append(s.idTable, val)
	return val
}

// The global Symbol Table of the Elk interpreter.
var SymbolTable *symbolTableStruct = newSymbolTable()
