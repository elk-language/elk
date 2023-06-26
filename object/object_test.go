package object

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// BENCHMARK: Try implementing this idea and measure
// the difference once the VM is fully operational.
//
// func Noop(...any) {}

// // Mimics the structure of the Go's interface.
// type ElkInterface struct {
// 	// itabs are allocated in non garbage collected memory
// 	// so using uintptr should be safe
// 	itab uintptr
// 	data unsafe.Pointer
// }

// const (
// 	UINT64_BITMASK  uintptr = 0b00
// 	INT64_BITMASK   uintptr = 0b01
// 	FLOAT64_BITMASK uintptr = 0b10
// )

// func Int64ToInterface(i Int64) ElkInterface {
// 	return ElkInterface{
// 		itab: uintptr(i),
// 		data: unsafe.Pointer(INT64_BITMASK),
// 	}
// }

// func (e ElkInterface) IsInt64() bool {
// 	return uintptr(e.data) == INT64_BITMASK
// }

// func (e ElkInterface) IsUInt64() bool {
// 	return uintptr(e.data) == UINT64_BITMASK
// }

// func (e ElkInterface) IsFloat64() bool {
// 	return uintptr(e.data) == FLOAT64_BITMASK
// }

// func (e ElkInterface) IsRefValue() bool {
// 	return uintptr(e.data) > FLOAT64_BITMASK
// }

// func (e ElkInterface) ToInt64() Int64 {
// 	return Int64(e.itab)
// }

// func (e ElkInterface) ToValue() Value {
// 	if !e.IsRefValue() {
// 		return nil
// 	}
// 	return *(*Value)(unsafe.Pointer(&e))
// }

// func BenchmarkEfaceScalar(b *testing.B) {
// 	var Uint UInt32
// 	b.Run("uint32", func(b *testing.B) {
// 		for i := 0; i < b.N; i++ {
// 			Uint = UInt32(i)
// 		}
// 	})
// 	var Eface Value
// 	b.Run("eface32", func(b *testing.B) {
// 		for i := 0; i < b.N; i++ {
// 			Eface = UInt32(i)
// 		}
// 	})
// 	var ElkIface ElkInterface
// 	b.Run("elkInterface", func(b *testing.B) {
// 		for i := 0; i < b.N; i++ {
// 			ElkIface = Int64ToInterface(Int64(i))
// 		}
// 	})
// 	Noop(Uint, Eface, ElkIface)
// }

func TestObjectInspect(t *testing.T) {
	tests := map[string]struct {
		obj  *Object
		want string
	}{
		"anonymous class and empty ivars": {
			obj:  NewObject(ObjectWithClass(NewClass())),
			want: `<anonymous>{}`,
		},
		"named class and empty ivars": {
			obj:  NewObject(ObjectWithClass(ExceptionClass)),
			want: `Std::Exception{}`,
		},
		"named class and ivars": {
			obj: NewObject(
				ObjectWithClass(ExceptionClass),
				ObjectWithInstanceVariables(
					SimpleSymbolMap{
						SymbolTable.Add("message").Id: String("foo bar!"),
					},
				),
			),
			want: `Std::Exception{ message: "foo bar!" }`,
		},
		"anonymous class and ivars": {
			obj: NewObject(
				ObjectWithClass(NewClass()),
				ObjectWithInstanceVariables(
					SimpleSymbolMap{
						SymbolTable.Add("message").Id: String("foo bar!"),
					},
				),
			),
			want: `<anonymous>{ message: "foo bar!" }`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.obj.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
