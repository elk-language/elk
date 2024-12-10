//go:build 386 || arm || armbe || mips || mipsle || ppc || s390 || sparc

package value

const FloatPrecision = 24

// Elk's Float value
type Float float32
