//go:build amd64 || amd64p32 || arm64 || arm64be || ppc64 || ppc64le || mips64 || mips64le || mips64p32 || mips64p32le || s390x || sparc64

package value

import (
	"unsafe"
)

func (i UInt64) ToValue() Value {
	inline := inlineValue{
		data: unsafe.Pointer(uintptr(UINT64_FLAG)),
		tab:  *(*uintptr)(unsafe.Pointer(&i)),
	}

	return *(*Value)(unsafe.Pointer(&inline))
}
