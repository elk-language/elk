//go:build amd64 || amd64p32 || arm64 || arm64be || ppc64 || ppc64le || mips64 || mips64le || mips64p32 || mips64p32le || s390x || sparc64

package value

func (d DateSpan) ToValue() Value {
	return Value{
		flag: DATE_SPAN_FLAG,
		data: uintptr(uint32(d.months))<<32 | uintptr(uint32(d.days)),
	}
}
