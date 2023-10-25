package value

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
	nameTable map[string]Symbol
	idTable   []string
	mutex     sync.RWMutex
}

type symbolTableOption func(*symbolTableStruct)

func symbolTableWithNameTable(nameTable map[string]Symbol) symbolTableOption {
	return func(s *symbolTableStruct) {
		s.nameTable = nameTable
	}
}

func symbolTableWithIdTable(idTable []string) symbolTableOption {
	return func(s *symbolTableStruct) {
		s.idTable = idTable
	}
}

func newSymbolTable(opts ...symbolTableOption) *symbolTableStruct {
	s := &symbolTableStruct{
		nameTable: make(map[string]Symbol, SYMBOL_TABLE_INITIAL_SIZE),
		idTable:   make([]string, 0, SYMBOL_TABLE_INITIAL_SIZE),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Check if a given symbol exists.
func (s *symbolTableStruct) ExistsId(symbol Symbol) bool {
	return symbol < Symbol(len(s.idTable)) && symbol > 0
}

// Check if a symbol with the given name exists.
// This function is thread-safe.
func (s *symbolTableStruct) Exists(name string) bool {
	_, ok := s.Get(name)
	return ok
}

// Get the Symbol with the specified name.
// This function is thread-safe.
func (s *symbolTableStruct) Get(name string) (Symbol, bool) {
	s.mutex.RLock()
	val, ok := s.nameTable[name]
	s.mutex.RUnlock()
	if !ok {
		return -1, false
	}

	return val, true
}

// Get the name of the given symbol.
// This function is thread-safe.
func (s *symbolTableStruct) GetName(symbol Symbol) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if symbol >= Symbol(len(s.idTable)) || symbol < 0 {
		return "", false
	}
	val := s.idTable[symbol]

	return val, true
}

// Add a new symbol with the specified name.
// This function is idempotent, if the Symbol already exists
// nothing happens and a pointer to it gets returned.
// This function is thread-safe.
func (s *symbolTableStruct) Add(name string) Symbol {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	val, ok := s.nameTable[name]
	if ok {
		return val
	}

	symbol := Symbol(len(s.idTable))
	s.nameTable[name] = symbol
	s.idTable = append(s.idTable, name)
	return symbol
}

// The global Symbol Table of the Elk interpreter.
var SymbolTable *symbolTableStruct = newSymbolTable()
