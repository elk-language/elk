//go:build amd64 || amd64p32 || arm64 || arm64be || ppc64 || ppc64le || mips64 || mips64le || mips64p32 || mips64p32le || s390x || sparc64

package value

import (
	"unsafe"
)

func (f Float64) ToValue() Value {
	return Value{
		flag: FLOAT64_FLAG,
		data: *(*uintptr)(unsafe.Pointer(&f)),
	}
}
