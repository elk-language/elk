package main

import (
	"testing"
	"unique"

	"github.com/elk-language/elk/value"
)

type NewSymbol unique.Handle[string]

func BenchmarkGetSymbolValue(b *testing.B) {
	symVal := value.ToSymbol("string")
	b.Run("symbol", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			symVal.String()
		}
	})

	symUnique := unique.Make("string")
	b.Run("unique", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			symUnique.Value()
		}
	})
}

func BenchmarkCreateSymbol(b *testing.B) {
	b.Run("symbol", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value.ToSymbol("string")
		}
	})

	b.Run("unique", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			unique.Make("string")
		}
	})
}

func noop(bool) {}

func BenchmarkCompareSymbol(b *testing.B) {
	symVal1 := value.ToSymbol("dupa")
	symVal2 := value.ToSymbol("trupa")
	b.Run("symbol", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			noop(symVal1 == symVal2)
		}
	})

	symUnique1 := unique.Make("dupa")
	symUnique2 := unique.Make("trupa")
	b.Run("unique", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			noop(symUnique1 == symUnique2)
		}
	})
}
